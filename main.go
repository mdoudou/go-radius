package main

import (
	"github.com/hel2o/go-radius/db"
	"github.com/hel2o/go-radius/g"
	"github.com/hel2o/go-radius/rpc"
	"github.com/hel2o/go-radius/server"
)

func main() {
	g.Init()
	db.InitDB()
	go rpc.RpcServerStart()
	server.StartRadius()

}
