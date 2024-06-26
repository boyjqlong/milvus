package funcutil

import (
	"fmt"
	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus/pkg/common"
)

func IsTextIndex(indexParams []*commonpb.KeyValuePair) bool {
	indexType, err := GetAttrByKeyFromRepeatedKV(common.IndexTypeKey, indexParams)
	if err != nil {
		return false
	}
	return indexType == "internal_fts"
}

func ComposeTextIndexName(collectionID, fieldID int64, collectionName string) string {
	return fmt.Sprintf("_internal_fts_%s_%d_%d", collectionName, collectionID, fieldID)
}
