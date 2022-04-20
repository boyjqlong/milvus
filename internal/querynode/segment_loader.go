// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package querynode

import (
	"context"
	"errors"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/milvus-io/milvus/internal/proto/schemapb"

	"go.uber.org/zap"

	"github.com/milvus-io/milvus/internal/common"
	etcdkv "github.com/milvus-io/milvus/internal/kv/etcd"
	"github.com/milvus-io/milvus/internal/log"
	"github.com/milvus-io/milvus/internal/metrics"
	"github.com/milvus-io/milvus/internal/mq/msgstream"
	"github.com/milvus-io/milvus/internal/proto/commonpb"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
	"github.com/milvus-io/milvus/internal/proto/querypb"
	"github.com/milvus-io/milvus/internal/storage"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/internal/util/funcutil"
	"github.com/milvus-io/milvus/internal/util/metricsinfo"
	"github.com/milvus-io/milvus/internal/util/timerecord"
)

const timeoutForEachRead = 10 * time.Second

// segmentLoader is only responsible for loading the field data from binlog
type segmentLoader struct {
	historicalReplica ReplicaInterface
	streamingReplica  ReplicaInterface

	dataCoord types.DataCoord

	cm     storage.ChunkManager // minio cm
	etcdKV *etcdkv.EtcdKV

	factory msgstream.Factory
}

func (loader *segmentLoader) getFieldType(segment *Segment, fieldID FieldID) (schemapb.DataType, error) {
	var coll *Collection
	var err error

	switch segment.getType() {
	case segmentTypeGrowing:
		coll, err = loader.streamingReplica.getCollectionByID(segment.collectionID)
		if err != nil {
			return schemapb.DataType_None, err
		}
	case segmentTypeSealed:
		coll, err = loader.historicalReplica.getCollectionByID(segment.collectionID)
		if err != nil {
			return schemapb.DataType_None, err
		}
	default:
		return schemapb.DataType_None, fmt.Errorf("invalid segment type: %s", segment.getType().String())
	}

	return coll.getFieldType(fieldID)
}

func (loader *segmentLoader) loadSegment(req *querypb.LoadSegmentsRequest, segmentType segmentType) error {
	if req.Base == nil {
		return fmt.Errorf("nil base message when load segment, collectionID = %d", req.CollectionID)
	}

	// no segment needs to load, return
	if len(req.Infos) == 0 {
		return nil
	}

	var metaReplica ReplicaInterface
	switch segmentType {
	case segmentTypeGrowing:
		metaReplica = loader.streamingReplica
	case segmentTypeSealed:
		metaReplica = loader.historicalReplica
	default:
		err := fmt.Errorf("illegal segment type when load segment, collectionID = %d", req.CollectionID)
		log.Error("load segment failed, illegal segment type", zap.Int64("loadSegmentRequest msgID", req.Base.MsgID), zap.Error(err))
		return err
	}

	log.Debug("segmentLoader start loading...",
		zap.Any("collectionID", req.CollectionID),
		zap.Any("numOfSegments", len(req.Infos)),
		zap.Any("loadType", segmentType),
	)
	// check memory limit
	concurrencyLevel := runtime.GOMAXPROCS(0)
	for ; concurrencyLevel > 1; concurrencyLevel /= 2 {
		err := loader.checkSegmentSize(req.CollectionID, req.Infos, concurrencyLevel)
		if err == nil {
			break
		}
	}

	err := loader.checkSegmentSize(req.CollectionID, req.Infos, concurrencyLevel)
	if err != nil {
		log.Error("load failed, OOM if loaded", zap.Int64("loadSegmentRequest msgID", req.Base.MsgID), zap.Error(err))
		return err
	}

	newSegments := make(map[UniqueID]*Segment)
	segmentGC := func() {
		for _, s := range newSegments {
			deleteSegment(s)
		}
	}

	for _, info := range req.Infos {
		segmentID := info.SegmentID
		partitionID := info.PartitionID
		collectionID := info.CollectionID

		collection, err := loader.historicalReplica.getCollectionByID(collectionID)
		if err != nil {
			segmentGC()
			return err
		}
		segment, err := newSegment(collection, segmentID, partitionID, collectionID, "", segmentType, true)
		if err != nil {
			log.Error("load segment failed when create new segment",
				zap.Int64("collectionID", collectionID),
				zap.Int64("partitionID", partitionID),
				zap.Int64("segmentID", segmentID),
				zap.Int32("segment type", int32(segmentType)),
				zap.Error(err))
			segmentGC()
			return err
		}

		newSegments[segmentID] = segment
	}

	loadSegmentFunc := func(idx int) error {
		loadInfo := req.Infos[idx]
		collectionID := loadInfo.CollectionID
		partitionID := loadInfo.PartitionID
		segmentID := loadInfo.SegmentID
		segment := newSegments[segmentID]
		tr := timerecord.NewTimeRecorder("loadDurationPerSegment")
		err = loader.loadSegmentInternal(segment, loadInfo)
		if err != nil {
			log.Error("load segment failed when load data into memory",
				zap.Int64("collectionID", collectionID),
				zap.Int64("partitionID", partitionID),
				zap.Int64("segmentID", segmentID),
				zap.Int32("segment type", int32(segmentType)),
				zap.Error(err))
			return err
		}

		metrics.QueryNodeLoadSegmentLatency.WithLabelValues(fmt.Sprint(Params.QueryNodeCfg.QueryNodeID)).Observe(float64(tr.ElapseSpan().Milliseconds()))

		return nil
	}
	// start to load
	err = funcutil.ProcessFuncParallel(len(req.Infos), concurrencyLevel, loadSegmentFunc, "loadSegmentFunc")
	if err != nil {
		segmentGC()
		return err
	}

	// set segment to meta replica
	for _, s := range newSegments {
		err = metaReplica.setSegment(s)
		if err != nil {
			log.Error("load segment failed, set segment to meta failed",
				zap.Int64("collectionID", s.collectionID),
				zap.Int64("partitionID", s.partitionID),
				zap.Int64("segmentID", s.segmentID),
				zap.Int64("loadSegmentRequest msgID", req.Base.MsgID),
				zap.Error(err))
			segmentGC()
			return err
		}
	}

	return nil
}

func (loader *segmentLoader) loadSegmentInternal(segment *Segment,
	loadInfo *querypb.SegmentLoadInfo) error {
	collectionID := loadInfo.CollectionID
	partitionID := loadInfo.PartitionID
	segmentID := loadInfo.SegmentID
	log.Debug("start loading segment data into memory",
		zap.Int64("collectionID", collectionID),
		zap.Int64("partitionID", partitionID),
		zap.Int64("segmentID", segmentID))

	pkFieldID, err := loader.historicalReplica.getPKFieldIDByCollectionID(collectionID)
	if err != nil {
		return err
	}

	var nonIndexedFieldBinlogs []*datapb.FieldBinlog
	var pkFieldBinlogs []*datapb.FieldBinlog
	if segment.getType() == segmentTypeSealed {
		fieldID2IndexInfo := make(map[int64]*querypb.FieldIndexInfo)
		for _, indexInfo := range loadInfo.IndexInfos {
			fieldID := indexInfo.FieldID
			fieldID2IndexInfo[fieldID] = indexInfo
		}

		indexedFieldInfos := make(map[int64]*IndexedFieldInfo)

		for _, fieldBinlog := range loadInfo.BinlogPaths {
			fieldID := fieldBinlog.FieldID
			if indexInfo, ok := fieldID2IndexInfo[fieldID]; ok {
				fieldInfo := &IndexedFieldInfo{
					fieldBinlog: fieldBinlog,
					indexInfo:   indexInfo,
				}
				indexedFieldInfos[fieldID] = fieldInfo

				if fieldBinlog.FieldID == pkFieldID {
					// pk field data should always be loaded.
					pkFieldBinlogs = append(pkFieldBinlogs, fieldBinlog)
				}
			} else {
				nonIndexedFieldBinlogs = append(nonIndexedFieldBinlogs, fieldBinlog)
			}
		}

		if err := loader.loadIndexedFieldData(segment, indexedFieldInfos); err != nil {
			return err
		}
	} else {
		nonIndexedFieldBinlogs = loadInfo.BinlogPaths
	}

	if err := loader.loadFiledBinlogData(segment, pkFieldBinlogs); err != nil {
		return err
	}
	if err := loader.loadFiledBinlogData(segment, nonIndexedFieldBinlogs); err != nil {
		return err
	}

	if pkFieldID == common.InvalidFieldID {
		log.Warn("segment primary key field doesn't exist when load segment")
	} else {
		log.Debug("loading bloom filter...")
		pkStatsBinlogs := loader.filterPKStatsBinlogs(loadInfo.Statslogs, pkFieldID)
		err = loader.loadSegmentBloomFilter(segment, pkStatsBinlogs)
		if err != nil {
			return err
		}
	}

	log.Debug("loading delta...")
	err = loader.loadDeltaLogs(segment, loadInfo.Deltalogs)
	return err
}

func (loader *segmentLoader) filterPKStatsBinlogs(fieldBinlogs []*datapb.FieldBinlog, pkFieldID int64) []string {
	result := make([]string, 0)
	for _, fieldBinlog := range fieldBinlogs {
		if fieldBinlog.FieldID == pkFieldID {
			for _, binlog := range fieldBinlog.GetBinlogs() {
				result = append(result, binlog.GetLogPath())
			}
		}
	}
	return result
}

func (loader *segmentLoader) loadFiledBinlogData(segment *Segment, fieldBinlogs []*datapb.FieldBinlog) error {
	if len(fieldBinlogs) <= 0 {
		return nil
	}

	segmentType := segment.getType()
	iCodec := storage.InsertCodec{}
	blobs := make([]*storage.Blob, 0)
	for _, fieldBinlog := range fieldBinlogs {
		for _, path := range fieldBinlog.Binlogs {
			binLog, err := loader.cm.Read(path.GetLogPath())
			if err != nil {
				return err
			}
			blob := &storage.Blob{
				Key:   path.GetLogPath(),
				Value: binLog,
			}
			blobs = append(blobs, blob)
		}
	}

	_, _, insertData, err := iCodec.Deserialize(blobs)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	switch segmentType {
	case segmentTypeGrowing:
		tsData, ok := insertData.Data[common.TimeStampField]
		if !ok {
			return errors.New("cannot get timestamps from insert data")
		}
		utss := make([]uint64, tsData.RowNum())
		for i := 0; i < tsData.RowNum(); i++ {
			utss[i] = uint64(tsData.GetRow(i).(int64))
		}

		rowIDData, ok := insertData.Data[common.RowIDField]
		if !ok {
			return errors.New("cannot get row ids from insert data")
		}

		return loader.loadGrowingSegments(segment, rowIDData.(*storage.Int64FieldData).Data, utss, insertData)
	case segmentTypeSealed:
		return loader.loadSealedSegments(segment, insertData)
	default:
		err := errors.New(fmt.Sprintln("illegal segment type when load segment, collectionID = ", segment.collectionID))
		return err
	}
}

func (loader *segmentLoader) loadIndexedFieldData(segment *Segment, vecFieldInfos map[int64]*IndexedFieldInfo) error {
	for fieldID, fieldInfo := range vecFieldInfos {
		if fieldInfo.indexInfo == nil || !fieldInfo.indexInfo.EnableIndex {
			fieldBinlog := fieldInfo.fieldBinlog
			err := loader.loadFiledBinlogData(segment, []*datapb.FieldBinlog{fieldBinlog})
			if err != nil {
				return err
			}
			log.Debug("load vector field's binlog data done", zap.Int64("segmentID", segment.ID()), zap.Int64("fieldID", fieldID))
		} else {
			indexInfo := fieldInfo.indexInfo
			err := loader.loadFieldIndexData(segment, indexInfo)
			if err != nil {
				return err
			}
			log.Debug("load vector field's index data done", zap.Int64("segmentID", segment.ID()), zap.Int64("fieldID", fieldID))
		}
		segment.setIndexedFieldInfo(fieldID, fieldInfo)
	}

	return nil
}

func (loader *segmentLoader) loadFieldIndexData(segment *Segment, indexInfo *querypb.FieldIndexInfo) error {
	indexBuffer := make([][]byte, 0)
	indexCodec := storage.NewIndexFileBinlogCodec()
	filteredPaths := make([]string, 0, len(indexInfo.IndexFilePaths))
	for _, p := range indexInfo.IndexFilePaths {
		log.Debug("load index file", zap.String("path", p))
		indexPiece, err := loader.cm.Read(p)
		if err != nil {
			return err
		}

		if path.Base(p) != storage.IndexParamsKey {
			data, _, _, _, err := indexCodec.Deserialize([]*storage.Blob{{Key: path.Base(p), Value: indexPiece}})
			if err != nil {
				return err
			}
			indexBuffer = append(indexBuffer, data[0].Value)
			filteredPaths = append(filteredPaths, p)
		}
	}
	// 2. use index bytes and index path to update segment
	indexInfo.IndexFilePaths = filteredPaths
	fieldType, err := loader.getFieldType(segment, indexInfo.FieldID)
	if err != nil {
		return err
	}
	return segment.segmentLoadIndexData(indexBuffer, indexInfo, fieldType)
}

func (loader *segmentLoader) loadGrowingSegments(segment *Segment,
	ids []UniqueID,
	timestamps []Timestamp,
	insertData *storage.InsertData) error {
	numRows := len(ids)
	if numRows != len(timestamps) || insertData == nil {
		return errors.New(fmt.Sprintln("illegal insert data when load segment, collectionID = ", segment.collectionID))
	}

	log.Debug("start load growing segments...",
		zap.Any("collectionID", segment.collectionID),
		zap.Any("segmentID", segment.ID()),
		zap.Any("numRows", len(ids)),
	)

	// 1. do preInsert
	var numOfRecords = len(ids)
	offset, err := segment.segmentPreInsert(numOfRecords)
	if err != nil {
		return err
	}
	log.Debug("insertNode operator", zap.Int("insert size", numOfRecords), zap.Int64("insert offset", offset), zap.Int64("segment id", segment.ID()))

	// 2. update bloom filter
	insertRecord, err := storage.TransferInsertDataToInsertRecord(insertData)
	if err != nil {
		return err
	}
	tmpInsertMsg := &msgstream.InsertMsg{
		InsertRequest: internalpb.InsertRequest{
			CollectionID: segment.collectionID,
			Timestamps:   timestamps,
			RowIDs:       ids,
			NumRows:      uint64(numRows),
			FieldsData:   insertRecord.FieldsData,
			Version:      internalpb.InsertDataVersion_ColumnBased,
		},
	}
	pks, err := getPrimaryKeys(tmpInsertMsg, loader.streamingReplica)
	if err != nil {
		return err
	}
	segment.updateBloomFilter(pks)

	// 3. do insert
	err = segment.segmentInsert(offset, ids, timestamps, insertRecord)
	if err != nil {
		return err
	}
	log.Debug("Do insert done in segment loader", zap.Int("len", numOfRecords), zap.Int64("segmentID", segment.ID()), zap.Int64("collectionID", segment.collectionID))

	return nil
}

func (loader *segmentLoader) loadSealedSegments(segment *Segment, insertData *storage.InsertData) error {
	insertRecord, err := storage.TransferInsertDataToInsertRecord(insertData)
	if err != nil {
		return err
	}
	numRows := insertRecord.NumRows
	for _, fieldData := range insertRecord.FieldsData {
		fieldID := fieldData.FieldId
		if fieldID == common.TimeStampField {
			timestampsData := insertData.Data[fieldID].(*storage.Int64FieldData)
			segment.setIDBinlogRowSizes(timestampsData.NumRows)
		}

		err := segment.segmentLoadFieldData(fieldID, numRows, fieldData)
		if err != nil {
			// TODO: return or continue?
			return err
		}
	}
	return nil
}

func (loader *segmentLoader) loadSegmentBloomFilter(segment *Segment, binlogPaths []string) error {
	if len(binlogPaths) == 0 {
		log.Info("there are no stats logs saved with segment", zap.Any("segmentID", segment.segmentID))
		return nil
	}

	values, err := loader.cm.MultiRead(binlogPaths)
	if err != nil {
		return err
	}
	blobs := make([]*storage.Blob, 0)
	for i := 0; i < len(values); i++ {
		blobs = append(blobs, &storage.Blob{Value: values[i]})
	}

	stats, err := storage.DeserializeStats(blobs)
	if err != nil {
		return err
	}
	for _, stat := range stats {
		if stat.BF == nil {
			log.Warn("stat log with nil bloom filter", zap.Int64("segmentID", segment.segmentID), zap.Any("stat", stat))
			continue
		}
		err = segment.pkFilter.Merge(stat.BF)
		if err != nil {
			return err
		}
	}
	return nil
}

func (loader *segmentLoader) loadDeltaLogs(segment *Segment, deltaLogs []*datapb.FieldBinlog) error {
	dCodec := storage.DeleteCodec{}
	var blobs []*storage.Blob
	for _, deltaLog := range deltaLogs {
		for _, log := range deltaLog.GetBinlogs() {
			value, err := loader.cm.Read(log.GetLogPath())
			if err != nil {
				return err
			}
			blob := &storage.Blob{
				Key:   log.GetLogPath(),
				Value: value,
			}
			blobs = append(blobs, blob)
		}
	}
	if len(blobs) == 0 {
		log.Info("there are no delta logs saved with segment", zap.Any("segmentID", segment.segmentID))
		return nil
	}
	_, _, deltaData, err := dCodec.Deserialize(blobs)
	if err != nil {
		return err
	}

	err = segment.segmentLoadDeletedRecord(deltaData.Pks, deltaData.Tss, deltaData.RowCount)
	if err != nil {
		return err
	}
	return nil
}

func (loader *segmentLoader) FromDmlCPLoadDelete(ctx context.Context, collectionID int64, position *internalpb.MsgPosition) error {
	log.Debug("from dml check point load delete", zap.Any("position", position), zap.Any("msg id", position.MsgID))
	stream, err := loader.factory.NewMsgStream(ctx)
	if err != nil {
		return err
	}

	defer func() {
		metrics.QueryNodeNumConsumers.WithLabelValues(fmt.Sprint(Params.QueryNodeCfg.QueryNodeID)).Dec()
		stream.Close()
	}()

	pChannelName := funcutil.ToPhysicalChannel(position.ChannelName)
	position.ChannelName = pChannelName

	stream.AsConsumer([]string{pChannelName}, fmt.Sprintf("querynode-%d-%d", Params.QueryNodeCfg.QueryNodeID, collectionID))
	lastMsgID, err := stream.GetLatestMsgID(pChannelName)
	if err != nil {
		return err
	}

	if lastMsgID.AtEarliestPosition() {
		log.Debug("there is no more delta msg", zap.Int64("Collection ID", collectionID), zap.String("channel", pChannelName))
		return nil
	}

	metrics.QueryNodeNumConsumers.WithLabelValues(fmt.Sprint(Params.QueryNodeCfg.QueryNodeID)).Inc()
	err = stream.Seek([]*internalpb.MsgPosition{position})
	if err != nil {
		return err
	}
	stream.Start()

	delData := &deleteData{
		deleteIDs:        make(map[UniqueID][]primaryKey),
		deleteTimestamps: make(map[UniqueID][]Timestamp),
		deleteOffset:     make(map[UniqueID]int64),
	}

	log.Debug("start read delta msg from seek position to last position",
		zap.Int64("Collection ID", collectionID), zap.String("channel", pChannelName))
	hasMore := true
	for hasMore {
		select {
		case <-ctx.Done():
			break
		case msgPack, ok := <-stream.Chan():
			if !ok {
				log.Warn("fail to read delta msg", zap.String("pChannelName", pChannelName), zap.Any("msg id", position.GetMsgID()), zap.Error(err))
				return err
			}

			if msgPack == nil {
				continue
			}

			for _, tsMsg := range msgPack.Msgs {
				if tsMsg.Type() == commonpb.MsgType_Delete {
					dmsg := tsMsg.(*msgstream.DeleteMsg)
					if dmsg.CollectionID != collectionID {
						continue
					}
					log.Debug("delete pk",
						zap.Any("pk", dmsg.PrimaryKeys),
						zap.String("vChannelName", position.GetChannelName()),
						zap.Any("msg id", position.GetMsgID()),
					)
					processDeleteMessages(loader.historicalReplica, dmsg, delData)
				}

				ret, err := lastMsgID.LessOrEqualThan(tsMsg.Position().MsgID)
				if err != nil {
					log.Warn("check whether current MsgID less than last MsgID failed",
						zap.Int64("Collection ID", collectionID), zap.String("channel", pChannelName), zap.Error(err))
					return err
				}

				if ret {
					hasMore = false
					break
				}
			}
		}
	}

	log.Debug("All data has been read, there is no more data", zap.Int64("Collection ID", collectionID),
		zap.String("channel", pChannelName), zap.Any("msg id", position.GetMsgID()))
	for segmentID, pks := range delData.deleteIDs {
		segment, err := loader.historicalReplica.getSegmentByID(segmentID)
		if err != nil {
			log.Debug(err.Error())
			continue
		}
		offset := segment.segmentPreDelete(len(pks))
		delData.deleteOffset[segmentID] = offset
	}

	wg := sync.WaitGroup{}
	for segmentID := range delData.deleteOffset {
		wg.Add(1)
		go deletePk(loader.historicalReplica, delData, segmentID, &wg)
	}
	wg.Wait()
	log.Debug("from dml check point load done", zap.Any("msg id", position.GetMsgID()))
	return nil
}

func deletePk(replica ReplicaInterface, deleteData *deleteData, segmentID UniqueID, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("QueryNode::iNode::delete", zap.Any("SegmentID", segmentID))
	targetSegment, err := replica.getSegmentByID(segmentID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if targetSegment.segmentType != segmentTypeSealed {
		return
	}

	ids := deleteData.deleteIDs[segmentID]
	timestamps := deleteData.deleteTimestamps[segmentID]
	offset := deleteData.deleteOffset[segmentID]

	err = targetSegment.segmentDelete(offset, ids, timestamps)
	if err != nil {
		log.Warn("QueryNode: targetSegmentDelete failed", zap.Error(err))
		return
	}
	log.Debug("Do delete done", zap.Int("len", len(deleteData.deleteIDs[segmentID])), zap.Int64("segmentID", segmentID), zap.Any("segmentType", targetSegment.segmentType))
}

// JoinIDPath joins ids to path format.
func JoinIDPath(ids ...UniqueID) string {
	idStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idStr = append(idStr, strconv.FormatInt(id, 10))
	}
	return path.Join(idStr...)
}

func (loader *segmentLoader) checkSegmentSize(collectionID UniqueID, segmentLoadInfos []*querypb.SegmentLoadInfo, concurrency int) error {
	usedMem := metricsinfo.GetUsedMemoryCount()
	totalMem := metricsinfo.GetMemoryCount()
	if len(segmentLoadInfos) < concurrency {
		concurrency = len(segmentLoadInfos)
	}

	if usedMem == 0 || totalMem == 0 {
		return fmt.Errorf("get memory failed when checkSegmentSize, collectionID = %d", collectionID)
	}

	usedMemAfterLoad := usedMem
	maxSegmentSize := uint64(0)
	for _, loadInfo := range segmentLoadInfos {
		segmentSize := uint64(loadInfo.SegmentSize)
		usedMemAfterLoad += segmentSize
		if segmentSize > maxSegmentSize {
			maxSegmentSize = segmentSize
		}
	}

	toMB := func(mem uint64) float64 {
		return float64(mem) / 1024 / 1024
	}

	// when load segment, data will be copied from go memory to c++ memory
	if usedMemAfterLoad+maxSegmentSize*uint64(concurrency) > uint64(float64(totalMem)*Params.QueryNodeCfg.OverloadedMemoryThresholdPercentage) {
		return fmt.Errorf("load segment failed, OOM if load, collectionID = %d, maxSegmentSize = %.2f MB, concurrency = %d, usedMemAfterLoad = %.2f MB, totalMem = %.2f MB, thresholdFactor = %f",
			collectionID, toMB(maxSegmentSize), concurrency, toMB(usedMemAfterLoad), toMB(totalMem), Params.QueryNodeCfg.OverloadedMemoryThresholdPercentage)
	}

	return nil
}

func newSegmentLoader(
	historicalReplica ReplicaInterface,
	streamingReplica ReplicaInterface,
	etcdKV *etcdkv.EtcdKV,
	cm storage.ChunkManager,
	factory msgstream.Factory) *segmentLoader {

	return &segmentLoader{
		historicalReplica: historicalReplica,
		streamingReplica:  streamingReplica,

		cm:     cm,
		etcdKV: etcdKV,

		factory: factory,
	}
}
