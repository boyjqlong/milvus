package parser

import (
	tparser "github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/mysql"
)

const (
	defaultCharset   = "utf8mb4"
	defaultCollation = "utf8mb4_bin"
)

type Config struct {
	Charset   string
	Collation string
	SQLMode   mysql.SQLMode

	EnableWindowFunction        bool
	EnableStrictDoubleTypeCheck bool
	SkipPositionRecording       bool
	CharsetClient               string // CharsetClient indicates how to decode the original SQL.
}

type Option = func(c *Config)

func (c *Config) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Config) ToTidbParserConfig() tparser.ParserConfig {
	return tparser.ParserConfig{
		EnableWindowFunction:        c.EnableWindowFunction,
		EnableStrictDoubleTypeCheck: c.EnableStrictDoubleTypeCheck,
		SkipPositionRecording:       c.SkipPositionRecording,
	}
}

func DefaultParserConfig() Config {
	return Config{
		Charset:   defaultCharset,
		Collation: defaultCollation,
		SQLMode:   mysql.ModeNone,
	}
}
