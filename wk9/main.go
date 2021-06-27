package main

import (
	"flag"
	"go-error/wk9/g"
	"go-error/wk9/server"
)

func main() {
	cfgTmp := flag.String("c", "cfg.json", "configuration file")
	g.ParseConfig(*cfgTmp)

	server.InitTCP(g.Config().Addr)
}
