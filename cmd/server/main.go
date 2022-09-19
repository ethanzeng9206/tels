package main

import (
	"fmt"
	"log"
	"net"
	"tels/decode"
	"tels/pb/huawei"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var address string

func Init()  {
	viper.SetConfigFile("conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	address = viper.GetString("server.address")
}


func main()  {
	Init()

	log.Println("telemery server")
	grpcServer := grpc.NewServer()
	// client := huawei.NewGRPCDataserviceClient()
	var telServer = decode.NewTelServer()

	huawei.RegisterGRPCDataserviceServer(grpcServer, telServer)


	// address := fmt.Sprint("0.0.0.0:10061")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("can not start server: err", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start server: err", err)
	}


}