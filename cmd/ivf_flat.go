package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/milvus-io/milvus/internal/common"

	"github.com/milvus-io/milvus/internal/proto/commonpb"

	"github.com/milvus-io/milvus/internal/proto/schemapb"

	"github.com/milvus-io/milvus/internal/proto/etcdpb"
	"github.com/milvus-io/milvus/internal/storage"

	"github.com/milvus-io/milvus/internal/indexnode"
)

func ServePProf() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
			panic("Start web service for pprof failed")
		}
	}()
}

func generateFloatVectors(nb, dim int) []float32 {
	vectors := make([]float32, 0)
	for i := 0; i < nb; i++ {
		for j := 0; j < dim; j++ {
			vectors = append(vectors, rand.Float32())
		}
	}
	return vectors
}

func generateInt64Array(numRows int, random bool) []int64 {
	ret := make([]int64, 0, numRows)
	for i := 0; i < numRows; i++ {
		if random {
			ret = append(ret, int64(rand.Int()))
		} else {
			ret = append(ret, int64(i)+2)
		}
	}
	return ret
}

var typeParams map[string]string
var indexType = "IVF_SQ8"
var indexParams = map[string]string{
	"index_type":  indexType,
	"metric_type": "L2",
	"dim":         "128",
	"nlist":       "100",
}
var nb = 1000000
var dim = 128

// var vectors = generateFloatVectors(nb, dim)
var tss = generateInt64Array(nb, false)

func indexOnlyCGO(i int, vectors []float32) {
	fmt.Printf("round %d\n", i)
	index, err := indexnode.NewCIndex(typeParams, indexParams)
	if err != nil {
		panic(err)
	}
	if err := index.BuildFloatVecIndexWithoutIds(vectors); err != nil {
		panic(err)
	}
	if _, err := index.Serialize(); err != nil {
		panic(err)
	}
	if err := index.Delete(); err != nil {
		panic(err)
	}
}

func indexWithCodecAndCGO(i int, vectors []float32) {
	fmt.Printf("round %d\n", i)
	codec := storage.NewInsertCodec(&etcdpb.CollectionMeta{
		Schema: &schemapb.CollectionSchema{
			Name:        "codec",
			Description: "",
			AutoID:      false,
			Fields: []*schemapb.FieldSchema{
				{
					FieldID:  101,
					Name:     "fvec",
					DataType: schemapb.DataType_FloatVector,
					IndexParams: []*commonpb.KeyValuePair{
						{
							Key:   "dim",
							Value: "128",
						},
					},
				},
			},
		},
	})
	blobs, _, err := codec.Serialize(1, 2, &storage.InsertData{
		Data: map[storage.FieldID]storage.FieldData{
			common.TimeStampField: &storage.Int64FieldData{
				NumRows: []int64{int64(nb)},
				Data:    tss,
			},
			101: &storage.FloatVectorFieldData{
				NumRows: []int64{int64(nb)},
				Data:    vectors,
				Dim:     dim,
			},
		},
		Infos: nil,
	})
	if err != nil {
		panic(err)
	}
	_, _, data, err := codec.Deserialize(blobs)
	if err != nil {
		panic(err)
	}
	indexOnlyCGO(i, data.Data[101].(*storage.FloatVectorFieldData).Data)
}

func main() {
	if len(os.Args) < 1 {
		_, _ = fmt.Fprint(os.Stderr, "usage: ivf_flat --withCodec\n")
		return
	}

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var withCodec bool
	flags.BoolVar(&withCodec, "withCodec", false, "with data codec")
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(-1)
	}

	ServePProf()

	i := 0
	for {
		var vectors = generateFloatVectors(nb, dim)
		if withCodec {
			indexWithCodecAndCGO(i, vectors)
		} else {
			indexOnlyCGO(i, vectors)
		}
		i++
	}
}
