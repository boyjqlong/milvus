# Generate Parser with Antlr4

## Install Antlr4

Please follow [install antlr4](https://github.com/antlr/antlr4/blob/master/doc/go-target.md) to install the antlr tool.

The version of antlr tool: `4.11.1`.

## Code Generate

```shell
go generate ./...
```

After code generation, the code tree will be like below:

```
.
├── antlr-4.11.1-complete.jar
├── generate.go
├── generate.sh
├── MySqlLexer.g4
├── mysql_lexer.go
├── MySqlLexer.interp
├── MySqlLexer.tokens
├── mysqlparser_base_visitor.go
├── MySqlParser.g4
├── mysql_parser.go
├── MySqlParser.interp
├── MySqlParser.tokens
├── mysqlparser_visitor.go
└── README.md
```
