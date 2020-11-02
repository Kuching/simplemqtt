// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package smqtt implements functions for manipulation of MQTT API.
// It provides a gin-like way to create MQTT API.
package core

import (
	"time"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
)

type ClientConfig struct {
	Broker string
	ClientId string 
	CaCertPath string
	ClientCertPath string
	ClientKeyPath string
	KeepAlive int
	PingTimeout int
	CleanSession bool
	OrderMatters bool
}

var defaultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	glog.V(1).Infoln(msg.Topic(), ": Handler NOT FOUND")
}

// newTLSConfig returns a SSL/TLS configuration to be used when connecting
// to an MQTT broker.
func newTLSConfig(caCertPath, clientCertPath, clientKeyPath string) (*tls.Config, error) {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(caCertPath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs: certpool,
		ClientAuth: tls.NoClientCert,
		ClientCAs: nil,
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
	}, nil
}

// NewClient returns a MQTT client
func newClient(conf *ClientConfig) (mqtt.Client, error) {
	if conf.Broker == "" {
		return nil, ParamRequiredError("Broker")
	}
	if conf.ClientId == "" {
		return nil, ParamRequiredError("ClientId")
	}
	if conf.CaCertPath == "" {
		return nil, ParamRequiredError("CaCertPath")
	}
	if conf.ClientCertPath == "" {
		return nil, ParamRequiredError("ClientCertPath")
	}
	if conf.ClientKeyPath == "" {
		return nil, ParamRequiredError("ClientKeyPath")
	}
	tlsconfig, err := newTLSConfig(conf.CaCertPath, conf.ClientCertPath, conf.ClientKeyPath)
	if err != nil {
		return nil, err
	}
	opts := mqtt.NewClientOptions().AddBroker(conf.Broker).SetClientID(conf.ClientId).SetTLSConfig(tlsconfig)
	if conf.KeepAlive == 0 {
		opts.SetKeepAlive(10 * time.Second)
	}else {
		opts.SetKeepAlive(time.Duration(conf.KeepAlive) * time.Second)
	}
	if conf.PingTimeout == 0 {
		opts.SetPingTimeout(10 * time.Second)
	}else {
		opts.SetPingTimeout(time.Duration(conf.PingTimeout) * time.Second)
	} 
	opts.SetCleanSession(conf.CleanSession)
	opts.SetOrderMatters(conf.OrderMatters)
	opts.SetDefaultPublishHandler(defaultHandler)
	
	Client := mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return Client, nil
}