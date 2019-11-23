// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package resp provides way to send MQTT message.
package resp

import (
	"github.com/eclipse/paho.mqtt.golang"
)

// Writer allows us to send MQTT message
type Writer struct{
	client mqtt.Client
}

// NewWriter returns a writer that can send MQTT message to MQTT client c
func NewWriter(c mqtt.Client) *Writer{
	return &Writer{
		client : c,
	}
}

// Write is used to send MQTT message
func(writer *Writer) Write(topic string, qos byte, retained bool, msg []byte) error{
	token := writer.client.Publish(topic, qos, retained, msg)
	token.Wait()
	return token.Error()
}