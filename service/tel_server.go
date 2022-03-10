package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"tels/pb/huawei"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/internal/status"
)


type TelServer struct {
	*huawei.UnimplementedGRPCDataserviceServer
}



func NewTelServer() *TelServer {
	return &TelServer{
		&huawei.UnimplementedGRPCDataserviceServer{},
	}
}

func (s *TelServer) DataPublish(stream huawei.GRPCDataservice_DataPublishServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		// req, err := stream.Recv()
		req, err := stream.Recv()
		
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			log.Print("can not receive stream request: ", err )
			return err
		}

		w := req.GetReqId()
		fmt.Println(w)
			
		// data := req.GetData()
		// fmt.Println(string(data))

		// jsonDato := req.GetDataJson()
		// fmt.Println(jsonData)

		// d := req.String()
		// fmt.Println(d)
		
		// e := req.GetErrors()
		// fmt.Println(e)

		data := req.GetData()
		// cpuInfo := &huawei.Debug_CpuInfos_CpuInfo{}
		// memInfo := &huawei.Debug_MemoryInfos_MemoryInfo{}
		// devmPhy := &huawei.Devm_PhysicalEntitys_PhysicalEntity{}
		// devmPhyClass := &huawei.Devm_PhysicalEntitys_PhysicalEntity_Class_name
		// ifm := &huawei.Ifm_Interfaces_Interface{
		// 	Name: "",
		// 	Index: 0,
		// 	MibStatistics: &huawei.Ifm_Interfaces_Interface_MibStatistics{},
		// 	CommonStatistics: &huawei.Ifm_Interfaces_Interface_CommonStatistics{},
		// }
		
		

		// err = proto.Unmarshal(data, cpuInfo)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// err = proto.Unmarshal(data, memInfo)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// err = proto.Unmarshal(data, devmPhy)
		// err = proto.Unmarshal(data, ifm)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(ProtobuffToJson(cpuInfo))
		// fmt.Println(ProtobuffToJson(memInfo))
		// fmt.Println(ProtobuffToJson(devmPhy))
		// ifmData, err := ProtobuffToJson(ifm)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(ifmData)
		
		ifmInfo := &huawei.Ifm_Interfaces_Interface{}
		err = proto.Unmarshal(data, ifmInfo)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ifmInfo)
		// for _, ifmIntInfo := range ifmInfo.Interface {
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	fmt.Println(ProtobuffToJson(ifmData))
		// }


		fmt.Println(" -----------------------------------")
	}

	return nil
}

func contextError(ctx context.Context) error{
	switch ctx.Err() {
	case context.Canceled:
		return fmt.Errorf("request is Canceled")
	case context.DeadlineExceeded:
		return fmt.Errorf("Deadline is exceeded")
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobuffToJson(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON data to file: %w", err)
	}

	return nil
}

func ProtobuffToJson(message proto.Message) (string, error) {
	marshalerOp := protojson.MarshalOptions{
		Indent: " ",
		UseProtoNames: true,
		EmitUnpopulated: true,
	}
	marshler, err := marshalerOp.Marshal(message)
	return string(marshler), err
}