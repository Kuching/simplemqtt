// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package core provides way to send MQTT message.
package core

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/eclipse/paho.mqtt.golang"
)

// Core allows us to send MQTT message
type Core struct{
	client mqtt.Client
}

// NewCore returns a Responser that can send MQTT message with MQTT client c
func NewCore(conf *ClientConfig) (*Core, error){
	c, err := newClient(conf)
	if err != nil {
		return nil, ClientError(err.Error())
	}
	return &Core{
		client : c,
	}, nil
}

// Pub is used to send MQTT message
func(c *Core) Pub(topic string, qos byte, retained bool, data interface{}) error{
	msg, _ := json.Marshal(data)
	token := c.client.Publish(topic, qos, retained, msg)
	token.Wait()
	if err := token.Error(); err != nil {
		glog.V(1).Infoln(err)
		glog.V(1).Infoln("<<< MQTT")
		return err
	}
	glog.V(1).Infof("%s %s\n%s", "<<< MQTT", topic, string(msg))
	return nil
}