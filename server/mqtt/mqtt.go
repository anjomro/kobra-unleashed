package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/anjomro/kobra-unleashed/server/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile(utils.GetEnv("MQTT_CAFILE", "samplecerts/ca-crt.pem"))
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(utils.GetEnv("MQTT_CLIENT_CERT", "samplecerts/client-crt.pem"), utils.GetEnv("MQTT_CLIENT_KEY", "samplecerts/client-key.pem"))
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

func getCommandTopic(cmdType string, action string) string {
	// Returns the topic where a command should be published
	return fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/20021/%s/%s/%s", utils.GetPrinterID(), cmdType, action)
}

func SendCommand(cmdType string, action string, payload map[string]interface{}) {
	// Generate a UUID
	msgID := uuid.New().String()

	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Add the msgid, timestamp, type, and action to the payload
	payload["msgid"] = msgID
	payload["timestamp"] = timestamp
	payload["type"] = cmdType
	payload["action"] = action

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Replace all occurrences of / with \/
	payloadStr := string(payloadBytes)
	payloadStr = strings.ReplaceAll(payloadStr, "/", "\\/")

	client := *GetMQTTClient()
	// Publish the message to the MQTT topic
	topic := getCommandTopic(cmdType, action)
	token := client.Publish(topic, 0, false, payloadStr)
	token.Wait()
}

func Print(filename string, filePath string) {
	// Seed the random number generator
	// Generate a random task id
	taskID := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000)

	// Create the data payload
	data := map[string]interface{}{
		"filename":  filename,
		"filepath":  filePath,
		"taskid":    taskID,
		"task_mode": 1,
		"filetype":  1,
	}

	// Create the payload for the command
	payload := map[string]interface{}{
		"data": data,
	}

	// Send the command
	SendCommand("print", "start", payload)
}

func SendPrintAction(taskID string, action string) {
	// Create the data payload
	data := map[string]interface{}{
		"taskid": taskID,
	}

	// Create the payload for the command
	payload := map[string]interface{}{
		"data": data,
	}

	// Send the command
	SendCommand("print", action, payload)
}
