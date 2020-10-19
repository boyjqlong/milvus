package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/czs007/suvlim/internal/conf"
	"github.com/czs007/suvlim/internal/reader"
	"strconv"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

    var yamlFile string
	flag.StringVar(&yamlFile, "yaml", "", "yaml file")
	flag.Parse()
	// flag.Usage()
	fmt.Println("yaml file: ", yamlFile)
	conf.LoadConfig(yamlFile)

	pulsarAddr := "pulsar://"
	pulsarAddr += conf.Config.Pulsar.Address
	pulsarAddr += ":"
	pulsarAddr += strconv.FormatInt(int64(conf.Config.Pulsar.Port), 10)
	reader.StartQueryNode(ctx, pulsarAddr)
}

