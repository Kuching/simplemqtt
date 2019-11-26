// Copyright 2019 Lya.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package smqtt implements functions for manipulation of MQTT API.
// It provides a gin-like way to create MQTT API.
package smqtt

import (
	"time"
	"github.com/eclipse/paho.mqtt.golang"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"tct/common/config"
	"github.com/golang/glog"
)

// newTLSConfig returns a SSL/TLS configuration to be used when connecting
// to an MQTT broker.
func newTLSConfig(conf config.Config) *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(conf.MqttCertificatePath + "/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	cert, err := tls.LoadX509KeyPair(conf.MqttCertificatePath + "/client.crt", conf.MqttCertificatePath + "/client.key")
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		RootCAs: certpool,
		ClientAuth: tls.NoClientCert,
		ClientCAs: nil,
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
	}
}

var defaultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	glog.V(1).Infoln(msg.Topic(), ": Handler NOT FOUND")
}

// NewClient returns a MQTT client which uses configuration from conf to set broker, clientID in client options, 
// and mqttCertificatePath in TLSConfig.
// And this MQTT client also has the following client options.
func NewClient(conf config.Config)  mqtt.Client {
	tlsconfig := newTLSConfig(conf)
	opts := mqtt.NewClientOptions().AddBroker(conf.MqttServer).SetClientID(conf.MqttClientId).SetTLSConfig(tlsconfig)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetDefaultPublishHandler(defaultHandler)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetCleanSession(false)
	opts.SetOrderMatters(false)
	Client := mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return Client
}
