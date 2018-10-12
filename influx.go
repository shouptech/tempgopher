package main

import (
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// WriteStateToInflux writes a State object to an Influx database
func WriteStateToInflux(s State, config Influx) error {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		UserAgent:          config.UserAgent,
		Timeout:            time.Duration(config.Timeout * 1000000000),
		InsecureSkipVerify: config.InsecureSkipVerify,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.Database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{"alias": s.Alias}
	fields := map[string]interface{}{"value": s.Temp}
	pt, err := client.NewPoint("temperature", tags, fields, s.When)
	if err != nil {
		return err
	}

	bp.AddPoint(pt)
	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}
