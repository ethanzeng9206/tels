package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"tels/pb/huawei"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
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

		// w := req.GetReqId()
		// fmt.Println(w)
		// jsonDato := req.GetDataJson()
		// fmt.Println(jsonData)

		data := req.GetData()
		// jsonData := req.GetDataJson()
		// fmt.Println(jsonData)
		// err = ioutil.WriteFile("tmp/test.json", []byte(jsonData), 0644)
		// if err != nil {
		// 	return fmt.Errorf("cannot write JSON data to file: %w", err)
		// }
		// defer DbClient().Close()
		// writeAPI := DbClient().WriteAPI("its", "tels")
		
		var huaweiTel = &huawei.Telemetry{}
		err = proto.Unmarshal(data, huaweiTel)
		if err != nil {
			fmt.Printf("解析头数据失败：%s\n", err)
		}

		if huaweiTel.SensorPath == "huawei-ifm:ifm/interfaces/interface/mib-statistics" {
			// defer DbClient().Close()
			// writeAPI := DbClient().WriteAPI("its", "tels")
			ifmIntArry := huaweiTel.GetDataGpb().GetRow()
			for _, ifmIntData := range ifmIntArry {
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	
				var ifmIntInfo = &huawei.Ifm{}
				err = proto.Unmarshal((ifmIntData.GetContent()), ifmIntInfo)
				if err != nil {
					fmt.Printf("解析ifm数据失败：%s\n", err)
				}
				ifmIntArry := ifmIntInfo.Interfaces.GetInterface()
				for _, ifmInt := range ifmIntArry {
					p := influxdb2.NewPointWithMeasurement("ifm").
					AddTag("Device", huaweiTel.NodeIdStr).
					AddTag("Port", ifmInt.Name).
					AddField("SendBytes", ifmInt.GetMibStatistics().GetSendByte()).
					AddField("ReceiveBytes",ifmInt.GetMibStatistics().GetReceiveByte()).
					AddField("SendPacket", ifmInt.GetMibStatistics().GetSendPacket()).
					AddField("ReceivePacket",ifmInt.GetMibStatistics().GetReceivePacket()).
					AddField("SendErrorPacket", ifmInt.GetMibStatistics().GetSendErrorPacket()).
					AddField("ReceiveErrorPacket", ifmInt.GetMibStatistics().GetReceiveErrorPacket()).
					AddField("SendDropPacket", ifmInt.GetMibStatistics().GetSendDropPacket()).
					AddField("ReceiveDropPacket", ifmInt.GetMibStatistics().GetReceiveDropPacket()).
					SetTime(time.UnixMilli(int64(ifmIntData.GetTimestamp())))
					// SetTime(time.Now())/
					// writeAPI.WritePoint(p)
					// DbClient().WritePoint(p)
					WirteDataToDb(p)
	
					// if ifmInt.Name == "GigabitEthernet2/0/2" {
					// 	fmt.Println(time.UnixMilli(int64(ifmIntData.GetTimestamp())))
					// 	fmt.Println(ifmInt.Name)
					// 	fmt.Println(ifmInt.GetMibStatistics().GetSendByte())
					// 	fmt.Println(ifmInt.GetMibStatistics().GetReceiveByte())
					// }
				}
			}
		}

		if huaweiTel.SensorPath == "huawei-debug:debug/cpu-infos/cpu-info" {
			// defer DbClient().Close()
			// writeAPI := DbClient().WriteAPI("its", "tels")
			cpuData := huaweiTel.GetDataGpb().GetRow()
			for _, cpuInfosArry := range cpuData {
				var cpuInfoArry = &huawei.Debug{}
				err := proto.Unmarshal(cpuInfosArry.GetContent(), cpuInfoArry)
				if err != nil {
					fmt.Printf("解析cpu数据失败：%s\n", err)
				}
				cpuInfos := cpuInfoArry.GetCpuInfos().GetCpuInfo()
				for _, cpuInfo := range cpuInfos {
					// fmt.Println(cpuInfo.Position)
					// fmt.Println(cpuInfo.GetSystemCpuUsage())
					// fmt.Println(ProtobuffToJson(cpuInfo))
					// fmt.Println("-----------------------------------")
					p := influxdb2.NewPointWithMeasurement("cpu").
					AddTag("Device", huaweiTel.NodeIdStr).
					AddTag("Position", cpuInfo.Position).
					AddField("CPU Usage", cpuInfo.SystemCpuUsage)
					WirteDataToDb(p)
					// writeAPI.WritePoint(p)
				} 
			}
		}

		if huaweiTel.SensorPath == "huawei-debug:debug/memory-infos/memory-info" {
			// defer DbClient().Close()
			// writeAPI := DbClient().WriteAPI("its", "tels")
			memData := huaweiTel.GetDataGpb().GetRow()
			for _, memInfosArry := range memData {
				var memInfoArry = &huawei.Debug{}
				err := proto.Unmarshal(memInfosArry.GetContent(), memInfoArry)
				if err != nil {
					fmt.Printf("解析内存数据失败：%s\n", err)
				}
				memInfos := memInfoArry.GetMemoryInfos().GetMemoryInfo()
				for _, memInfo := range memInfos {
					// fmt.Println(ProtobuffToJson(memInfo))
					p := influxdb2.NewPointWithMeasurement("mem").
					AddTag("Device", huaweiTel.NodeIdStr).
					AddTag("Position", memInfo.Position).
					AddField("Mem Total", memInfo.GetOsMemoryTotal()).
					AddField("Mem Usage", memInfo.GetDoMemoryUse())
					WirteDataToDb(p)
					// writeAPI.WritePoint(p)
				}

			}
		}
		// writeAPI.Flush()
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


func ConvertBytes(bytes []byte, charset string) ([]byte, error){
	switch charset {
	case "GB18030":
		decodeBytes, err := simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
		if err != nil {
			return nil, err
		}
		return decodeBytes, nil
	case "GBK":
		decodeBytes, err := simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
		if err != nil {
			return nil, err
		}
		return decodeBytes, nil
	default:
		return bytes, nil
	}
}

