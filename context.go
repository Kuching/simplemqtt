// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package smqtt implements functions for manipulation of MQTT API.
// It provides a gin-like way to create MQTT API.
package smqtt

import (
	"math"
	"encoding/json"

	"github.com/eclipse/paho.mqtt.golang"
	
)

const abortIndex int8 = math.MaxInt8 / 2

// Context allows us to pass variables between middleware,
// manage the flow, validate the JSON of a MQTT Message and 
// render a JSON response.
type Context struct {
	handlers []Handler
	router *Router
	keys map[string]interface{}
	index int8
}

// When a response is rendered using context.Respond, it will be recorded
// in context.keys in the form of Response.
type Response struct {
	Topic string
	Qos byte
	Msg []byte
}

// Set is used to store a new key/value pair in COntext.keys.
// It will create an empty map when context.keys is currently nil.
func (c *Context) Set(key string, value interface{}) {
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// Get returns the value for the given key.
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.keys != nil {
		value, exists = c.keys[key]
	}
	return
}

// MustGet returns the value for the given key if it existes, 
// otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// Next is used to execute context.handlers[context.index+1:]
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort is used to abort the execution of context.handlers.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Start is used to call context.handlers[0]
func (c *Context) Start() {
	c.index = 0
	s := int8(len(c.handlers))
	if c.index < s{
		c.handlers[c.index](c)
	}	
}

// BindJSON is used to unmarshal the MQTT message to obj
func (c *Context) BindJSON(obj interface{}) error {
	msg := c.MustGet("mqtt-msg").(mqtt.Message)
	return json.Unmarshal(msg.Payload(), &obj)
}

// Respond is used to write and send a MQTT message to certain topic as response
func(c *Context) Respond(topic string, qos byte, retained bool, msg []byte) error {
	//err := c.router.getWriter().Write(topic, qos, retained, msg)
	token := c.router.getClient().Publish(topic, qos, retained, msg)	
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.Set("mqtt-resp", Response{topic, qos, msg})
	return nil
}

func (c *Context) setHandlers(handlers []Handler){
	c.handlers = handlers
}

// reset is used to clear context.keys
func (c *Context) reset(){
	c.keys = make(map[string]interface{})
}