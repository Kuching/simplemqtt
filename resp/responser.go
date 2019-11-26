// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package resp provides way to send MQTT message.
package resp

import (
	"github.com/eclipse/paho.mqtt.golang"
)

// Responser allows us to send MQTT message
type Responser struct{
	client mqtt.Client
}

// NewResponser returns a Responser that can send MQTT message with MQTT client c
func NewResponser(c mqtt.Client) *Responser{
	return &Responser{
		client : c,
	}
}

// Publish is used to send MQTT message
func(resp *Responser) Publish(topic string, qos byte, retained bool, msg []byte) error{
	token := resp.client.Publish(topic, qos, retained, msg)
	token.Wait()
	return token.Error()
}