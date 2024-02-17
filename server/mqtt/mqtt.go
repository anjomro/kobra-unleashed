package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/anjomro/kobra-unleashed/server/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile(utils.GetEnv("MQTT_CAFILE", "certs/ca.crt"))
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(utils.GetEnv("MQTT_CLIENT_CERT", "certs/client.crt"), utils.GetEnv("MQTT_CLIENT_KEY", "certs/client.key"))
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf)

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

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var mqtttClient MQTT.Client

func GetMQTTClient() *MQTT.Client {
	if mqtttClient == nil {
		opts := MQTT.NewClientOptions()
		opts.AddBroker(utils.GetEnv("MQTT_BROKER", "mqtts://localhost:8883"))
		opts.SetClientID(utils.GetEnv("MQTT_CLIENT_ID", "kobra-client"))
		opts.SetTLSConfig(NewTLSConfig())
		opts.SetDefaultPublishHandler(f)

		mqtttClient = MQTT.NewClient(opts)
		if token := mqtttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		} else {
			fmt.Println("MQTT Connected")
		}

	}
	return &mqtttClient
}
