package structs

type MqttPayload struct {
	MsgID     string `json:"msgid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Action    string `json:"action"`
	// Add optional Data struct that can be anything
	Data interface{} `json:"data"`
}

//	{
//		"type":	"tempature",
//		"action":	"auto",
//		"msgid":	"77350b05-e3d9-4a10-9c12-502a439d7c3d",
//		"state":	"done",
//		"timestamp":	1709681948443,
//		"code":	200,
//		"msg":	"",
//		"data":	{
//			"curr_hotbed_temp":	27,
//			"curr_nozzle_temp":	27,
//			"target_hotbed_temp":	0,
//			"target_nozzle_temp":	0
//		}
//	}
type MqttResponse struct {
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	MsgID     string      `json:"msgid"`
	State     string      `json:"state"`
	Timestamp int64       `json:"timestamp"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

type MqttTempatureData struct {
	CurrHotbedTemp   int `json:"curr_hotbed_temp"`
	CurrNozzleTemp   int `json:"curr_nozzle_temp"`
	TargetHotbedTemp int `json:"target_hotbed_temp"`
	TargetNozzleTemp int `json:"target_nozzle_temp"`
}
