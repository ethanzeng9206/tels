package service

import (
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/spf13/viper"
)

type influxdbConfig struct {
	Url 		string
	Token 		string
	Org 		string
	Bucket 		string
}
var Influxdb influxdbConfig

func Init()  {
	viper.SetConfigFile("conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.UnmarshalKey("influxdb", &Influxdb)
}

func WirteDataToDb(p *write.Point) {
	client := influxdb2.NewClient(Influxdb.Url, Influxdb.Token)
	defer client.Close()

	writeAPI := client.WriteAPI(Influxdb.Org, Influxdb.Bucket)
	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

func DbClient() influxdb2.Client {
	Init()
	client := influxdb2.NewClient(Influxdb.Url, Influxdb.Token)
	client.Options().SetBatchSize(120000)
	client.Options().SetMaxRetries(240000)
	return client
}

