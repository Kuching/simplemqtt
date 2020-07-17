// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package smqtt implements functions for manipulation of MQTT API.
// It provides a gin-like way to create MQTT API.
package smqtt

import (
	"sync"
	"bytes"
	"time"
	"fmt"
	"runtime"
	"io/ioutil"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	
)

// This definition indicates that a handler implementaion 
// must and only has a *Context as parameter 
type Handler func(*Context) 

// Router allows us to create Group and read/write MQTT message
type Router struct {
	group *Group
	client mqtt.Client
	// pool is a sync.Pool for context-multiplex
	pool sync.Pool
}

// NewRouter returns a *Router that uses two default middlewares logger and recovery
func NewRouter(c mqtt.Client) *Router {	
	router := &Router{
		group: &Group{},
		client: c,
	}
	router.Use(logger(), recovery())
	router.pool.New = func() interface{} {
		return router.allocateContext()
	}
	return router
}

// Group creates a new group
func (router *Router) Group(topic string) *Group{
	group := &Group{
		topic: topic,
		handlers: router.group.handlers,
		router: router,
	}
	return group
}

// Use add middlewares to router.group.handlers
func (router *Router) Use(handler ...Handler) {
	for _, h := range handler{
		router.group.handlers = append(router.group.handlers, h)
	}
}

func (router *Router) getClient() mqtt.Client{
	return router.client
}

func (router *Router) allocateContext() *Context{
	return &Context{router: router}
}

// logger combine the MQTT request and response as log content, and log after response. 
func logger() Handler {
	return func(c *Context) {
		start := time.Now()
		msg := c.MustGet("mqtt-msg").(mqtt.Message)
		sid := c.MustGet("mqtt-session").(string)
		b := bytes.NewBufferString(fmt.Sprintf("\n[%s] %s %s", sid, ">>> MQTT ", msg.Topic()))
		b.WriteString(fmt.Sprintf("\n[%s] %s", sid, string(msg.Payload())))
		c.Next()
		end := time.Now()

		resp, exists := c.Get("mqtt-resp")
		if exists {
			r := resp.(Response)
			b.WriteString(fmt.Sprintf("\n[%s] %s %s", sid, "<<< MQTT ", r.Topic))
			b.WriteString(fmt.Sprintf("\n[%s] %s", sid, string(r.Msg)))
		}else{
			b.WriteString(fmt.Sprintf("\n[%s] %s", sid, "<<< MQTT "))
		}
		latency := end.Sub(start)
		b.WriteString(fmt.Sprintf("\n[%s] %v - %v |%13v\n\n",
			sid,
			start.Format("0102 15:04:05.000000"),
			end.Format("0102 15:04:05.000000"),
			latency,
		))
		glog.V(1).Info(b.String())
	}
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// recovery enables program to recover after it panics
func recovery() Handler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				glog.Infof("[Recovery] panic recovered:\n%s\n%s", err, stack)
			}
		}()
		c.Next()
	}
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

