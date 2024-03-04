package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"

	"github.com/anjomro/kobra-unleashed/server/structs"
	"github.com/anjomro/kobra-unleashed/server/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile(utils.GetEnv("MQTT_CAFILE", "/user/ca.crt"))
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(utils.GetEnv("MQTT_CLIENT_CERT", "/user/client.crt"), utils.GetEnv("MQTT_CLIENT_KEY", "/user/client.key"))
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

var MQTTClient *MQTT.Client

func GetMQTTClient() *MQTT.Client {
	if MQTTClient == nil {
		opts := MQTT.NewClientOptions()
		opts.AddBroker(utils.GetEnv("MQTT_BROKER", "ssl://localhost:8883"))
		opts.SetClientID(utils.GetEnv("MQTT_CLIENT_ID", "go-simple"))
		opts.SetTLSConfig(NewTLSConfig())
		opts.SetDefaultPublishHandler(f)
		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		MQTTClient = &c
	}
	return MQTTClient
}

func getCommandTopic() (string, error) {
	// Returns the topic where a command should be published
	printerModel, err := kobrautils.GetPrinterModel()
	if err != nil {
		return "", err
	}

	printerId, err := kobrautils.GetPrinterID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/%s/%s", printerModel, printerId), nil
}

func SendCommand(payload *structs.MqttPayload) error {
	// Generate a UUID
	msgID := uuid.New().String()

	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Add the msgid, timestamp, type, and action to the payload

	// Create the payload for the command

	payload.MsgID = msgID
	payload.Timestamp = timestamp

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Replace all occurrences of / with \/
	payloadStr := string(payloadBytes)
	payloadStr = strings.ReplaceAll(payloadStr, "/", "\\/")

	client := *GetMQTTClient()
	// Publish the message to the MQTT topic
	topic, err := getCommandTopic()
	if err != nil {
		return err
	}
	token := client.Publish(topic, 0, false, payloadStr)
	token.Wait()

	slog.Info("MQTT", "Published message to topic", topic)
	slog.Info("MQTT", "payload", payloadStr)

	return nil
}
