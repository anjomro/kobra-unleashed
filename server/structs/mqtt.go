package structs

type MqttPayload struct {
	MsgID     string `json:"msgid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Action    string `json:"action"`
	// Add optional Data struct that can be anything
	Data interface{} `json:"data"`
}
