package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"
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

// CreatePayload creates a json payload from a map
func CreatePayload(data map[string]interface{}) ([]byte, error) {

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	return payloadBytes, nil
}

// SendMQTTCommand sends a command to the printer with json payload
func SendMQTTCommand(payload []byte) error {
	mqttclinet := *GetMQTTClient()

	// Check if the MQTT client is connected
	if !mqttclinet.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	printerModel, err := kobrautils.GetPrinterModel()
	if err != nil {
		return err
	}

	printerID, err := kobrautils.GetPrinterID()
	if err != nil {
		return err
	}

	// Convert map string interface to interface

	// anycubic/anycubicCloud/v1/server/printer/<PRINTER_MODEL_ID>/<PRINTER_ID>/response
	token := mqttclinet.Publish(fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/%s/%s/response", printerModel, printerID), 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("[SendMQTTCommand]: error sending mqtt command: %v", token.Error())
	} else {
		fmt.Println("MQTT Command Sent")
	}
	if utils.IsDev() {
		fmt.Println("MQTT URL:", fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/%s/%s/response", printerModel, printerID))
		jsonpayload, _ := json.Marshal(payload)
		fmt.Println("MQTT Payload:", string(jsonpayload))
	}
	return nil
}
