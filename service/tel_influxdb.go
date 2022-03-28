package service

import (
	"github.com/influxdata/influxdb-client-go/v2"
	// "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)



func WirteDataToDb(p *write.Point) {
	client := influxdb2.NewClient("http://10.158.1.123:8086", "6xbs-nUI6oQjDvhUSUc8VBtzdvl324yY5JdKYv9enGrDrwf4jhjYTaiaAxHWRJzQchodRMVdaBJRnLsNzruUUA==")
	defer client.Close()
	// client.Options().SetFlushInterval(1)
	
	writeAPI := client.WriteAPI("its", "tels")
	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

func DbClient() influxdb2.Client {
	client := influxdb2.NewClient("http://10.158.1.123:8086", "6xbs-nUI6oQjDvhUSUc8VBtzdvl324yY5JdKYv9enGrDrwf4jhjYTaiaAxHWRJzQchodRMVdaBJRnLsNzruUUA==")
	// writeAPI := client.WriteAPI("its", "tels")
	return client
}

