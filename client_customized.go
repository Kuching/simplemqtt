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

// newTLSConfig_C returns a SSL/TLS configuration to be used when connecting
// to an MQTT broker.
// Modify this function if you want your own customized SSL/TLS configuration.
func newTLSConfig_C(conf config.Config) *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(conf.MqttCert + "/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// // Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(conf.MqttCert + "/client.crt", conf.MqttCert + "/client.key")
	if err != nil {
		panic(err)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var defaultHandler_C mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	glog.V(1).Infoln(msg.Topic(), ": Handler NOT FOUND")
}

// NewClient_C returns a MQTT client which uses configuration from conf to set broker, clientID in client options, 
// and mqttCertificatePath in TLSConfig.
// And this MQTT client also has the following client options.
// Modify this function if you want your own customized MQTT client.
func NewClient_C(conf config.Config) mqtt.Client {
	tlsconfig := newTLSConfig_C(conf)
	// 
	// See more options setting in github.com/eclipse/paho.mqtt.golang/options.go
	//
	opts := mqtt.NewClientOptions()
	// AddBroker adds a broker URI to the list of brokers to be used. The format should be
	// scheme://host:port
	// Where "scheme" is one of "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
	// and "port" is the port on which the broker is accepting connections.
	//
	// Default values for hostname is "127.0.0.1", for schema is "tcp://".
	//
	// An example broker URI would look like: tcp://foobar.com:1883
	opts.AddBroker(conf.MqttServer)
	// SetClientID will set the client id to be used by this client when
	// connecting to the MQTT broker. According to the MQTT v3.1 specification,
	// a client id mus be no longer than 23 characters.
	opts.SetClientID(conf.MqttClientId)
	// SetTLSConfig will set an SSL/TLS configuration to be used when connecting
	// to an MQTT broker. Please read the official Go documentation for more
	// information.
	opts.SetTLSConfig(tlsconfig)
	// SetKeepAlive will set the amount of time (in seconds) that the client
	// should wait before sending a PING request to the broker. This will
	// allow the client to know that a connection has not been lost with the
	// server.
	opts.SetKeepAlive(10 * time.Second)
	// SetDefaultPublishHandler sets the MessageHandler that will be called when a message
	// is received that does not match any known subscriptions.
	opts.SetDefaultPublishHandler(defaultHandler_C)
	// SetPingTimeout will set the amount of time (in seconds) that the client
	// will wait after sending a PING request to the broker, before deciding
	// that the connection has been lost. Default is 10 seconds.
	opts.SetPingTimeout(10 * time.Second)
	// SetCleanSession will set the "clean session" flag in the connect message
	// when this client connects to an MQTT broker. By setting this flag, you are
	// indicating that no messages saved by the broker for this client should be
	// delivered. Any messages that were going to be sent by this client before
	// diconnecting previously but didn't will not be sent upon connecting to the
	// broker.
	opts.SetCleanSession(false)
	// SetOrderMatters will set the message routing to guarantee order within
	// each QoS level. By default, this value is true. If set to false,
	// this flag indicates that messages can be delivered asynchronously
	// from the client to the application and possibly arrive out of order.
	// opts.SetOrderMatters(false)

	Client := mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return Client
}
