// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package smqtt implements functions for manipulation of MQTT API.
// It provides a gin-like way to create MQTT API.
package smqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/Kuching/simple_mqtt/session"
)

// Group allows us to create MQTT API and manipulate them as a group.
type Group struct {
	router *Router
	topic string
	handlers []Handler
}

// Use adds middleware to group.handlers.
func (group *Group) Use(handler ...Handler) {
	group.handlers = combine(group.handlers, handler)
}

// Listen is used to make real subscription to MQTT broker.
func (group *Group) Listen(sub_topic string, qos byte, handler ...Handler) error{
	handlers := combine(group.handlers, handler)
	var msg_handler = func(client mqtt.Client, msg mqtt.Message){
		c := group.router.pool.Get().(*Context)
		c.reset()
		c.setHandlers(handlers)
		c.Set("mqtt-msg", msg)
		c.Set("mqtt-session", SessionID())
		c.Start()
		group.router.pool.Put(c)
	}
	token := group.router.getClient().Subscribe(group.topic+sub_topic, qos, msg_handler)
	token.Wait()
	return token.Error()
}

func combine(hds1 []Handler, hds2 []Handler) []Handler {
	for _, h := range(hds2){
		hds1 = append(hds1, h)
	}
	return hds1
}