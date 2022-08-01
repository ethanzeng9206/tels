package decode

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"tels/pb/huawei"
	"tels/service"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	// 这里需要加for循环，否则会丢数据，时间节点不匹配
	// for {
	// 	// 这里不要直接返回error,不然如果程序运行过程中出现错误就直接终止了
	// 	err := contextError(stream.Context())
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// }
	for i:=0; i<6; i++ {
		go func (i int)  {
			for stream != nil {
				decode(stream)
			}
		}(i)
	}
	return nil
}

func decode(stream huawei.GRPCDataservice_DataPublishServer) error {
	err := contextError(stream.Context())
	if err != nil {
		service.Logger.Error(err.Error())
	}

	req, err := stream.Recv()
	if err == io.EOF {
		service.Logger.Info("no more received stream")
	}
	if err != nil {
		service.Logger.Errorf("can not receive stream request: %s" , err)
	}
	
	rawData := req.GetData()

	var telData = &huawei.Telemetry{}
	err = proto.Unmarshal(rawData, telData)
	if err != nil {
		// log.Printf("解析头数据失败：%s\n", err)
		service.Logger.Errorf("解析头数据失败: %s", err)
	}

	client := service.DbClient()
	writeApi := client.WriteAPI("its", "tels")

	switch telData.SensorPath {
	case "huawei-ifm:ifm/interfaces/interface/mib-statistics":
		service.Logger.Info("接收来自%s的网卡流量数据包", telData.GetNodeIdStr())
		ifmArryData := telData.GetDataGpb().GetRow()
		for _, ifmRawData := range ifmArryData {
			var ifmData = &huawei.Ifm{}
			err = proto.Unmarshal((ifmRawData.GetContent()), ifmData)
			if err != nil {
				service.Logger.Errorf("解析ifm数据失败： %s", err)
			}

			ifmIntData := ifmData.Interfaces.GetInterface()
			for _, intData := range ifmIntData {
				p := influxdb2.NewPointWithMeasurement("ifmTraffic").
				AddTag("Device", telData.GetNodeIdStr()).AddTag("Port", intData.GetName()).
				AddField("SendBits", intData.GetMibStatistics().GetSendByte() * 8).AddField("ReceiveBits", intData.GetMibStatistics().GetReceiveByte() * 8).
				AddField("SendPacket", intData.GetMibStatistics().GetSendPacket()).AddField("ReceivePacket", intData.GetMibStatistics().GetReceivePacket()).
				SetTime(time.UnixMilli(int64(ifmRawData.GetTimestamp())))
				writeApi.WritePoint(p)
			}
		}
	case "huawei-debug:debug/cpu-infos/cpu-info":
		service.Logger.Info("接收来自%s的CPU使用数据包", telData.GetNodeIdStr())
		devmArryData := telData.GetDataGpb().GetRow()
		for _, devmRawData := range devmArryData{
			var devmData = &huawei.Debug{}
			err := proto.Unmarshal(devmRawData.GetContent(), devmData)
			if err != nil {
				service.Logger.Errorf("解析devm数据失败: %s", err)
			}

			cpuInfos := devmData.GetCpuInfos().GetCpuInfo()

			for _, cpuInfo := range cpuInfos {
				p := influxdb2.NewPointWithMeasurement("cpu").
				AddTag("Device", telData.NodeIdStr).
				AddTag("Position", cpuInfo.Position).
				AddField("CPU Usage", cpuInfo.SystemCpuUsage).
				SetTime(time.UnixMilli(int64(devmRawData.GetTimestamp())))
				// WirteDataToDb(p)
				writeApi.WritePoint(p)
			}
		}
	case "huawei-debug:debug/memory-infos/memory-info":
		service.Logger.Info("接收来自%s的内存使用数据包", telData.GetNodeIdStr())
		devmArryData := telData.GetDataGpb().GetRow()
		for _, devmRawData := range devmArryData{
			var devmData = &huawei.Debug{}
			err := proto.Unmarshal(devmRawData.GetContent(), devmData)
			if err != nil {
				service.Logger.Errorf("解析devm数据失败: %s", err)
			}

			memInfos := devmData.GetMemoryInfos().GetMemoryInfo()

			for _, memInfo := range memInfos {
				p := influxdb2.NewPointWithMeasurement("mem").
				AddTag("Device", telData.NodeIdStr).
				AddTag("Position", memInfo.Position).
				AddField("Mem Total", memInfo.GetOsMemoryTotal()).
				AddField("Mem Usage", memInfo.GetDoMemoryUse()).
				SetTime(time.UnixMilli(int64(devmRawData.GetTimestamp())))
				writeApi.WritePoint(p)
			}
		}
	default:	
		break
	}
	writeApi.Flush()

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
