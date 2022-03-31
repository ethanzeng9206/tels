package main

import (
	"fmt"
	"log"
	"net"
	"tels/pb/huawei"
	"tels/service"

	"google.golang.org/grpc"
)

func main()  {
	log.Println("telemery server")
	grpcServer := grpc.NewServer()
	// client := huawei.NewGRPCDataserviceClient()
	var telServer = service.NewTelServer()

	huawei.RegisterGRPCDataserviceServer(grpcServer, telServer)


	address := fmt.Sprint("0.0.0.0:10061")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("can not start server: err", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start server: err", err)
	}


}