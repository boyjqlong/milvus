package querycoord

import (
	"github.com/milvus-io/milvus/internal/proto/commonpb"
)

type extraIndexInfo struct {
	indexID        UniqueID
	indexName      string
	indexParams    []*commonpb.KeyValuePair
	indexSize      uint64
	indexFilePaths []string
}
