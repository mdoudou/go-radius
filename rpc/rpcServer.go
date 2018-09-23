package rpc

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/hel2o/go-radius/g"
)

type AgentRequest struct {
	IpAddress     string
	ServerVersion string
}
type AgentResponse struct {
	Status    int    `json:"Status"`
	Version   string `json:"Version"`
	StartTime int64  `json:"StartTime"`
}

func (t *Agent) ReportAgentStatus(args *AgentRequest, resp *AgentResponse) error {
	if g.Config().Rpc.Management == args.IpAddress {
		resp.Status = 1
		resp.Version = g.VERSION
		resp.StartTime = g.StartTime
		return nil
	}
	return errors.New("auth fail")

}

type Agent int

func RpcServerStart() {
	addr := g.Config().Rpc.Listen

	server := rpc.NewServer()
	server.Register(new(Agent))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("rpc listen error:", e)
	} else {
		log.Println("rpc listening", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("rpc listener accept fail:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
