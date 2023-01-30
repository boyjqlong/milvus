package tidbparser

import (
	"sync"

	tparser "github.com/pingcap/tidb/parser"
)

var parserPool = &sync.Pool{New: func() interface{} { return tparser.New() }}
