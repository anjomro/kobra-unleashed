// type MqttResponse struct {
// 	Type      string      `json:"type"`
// 	Action    string      `json:"action"`
// 	MsgID     string      `json:"msgid"`
// 	State     string      `json:"state"`
// 	Timestamp int64       `json:"timestamp"`
// 	Code      int         `json:"code"`
// 	Msg       string      `json:"msg"`
// 	Data      interface{} `json:"data"`
// }

// type MqttTempatureData struct {
// 	CurrHotbedTemp   int `json:"curr_hotbed_temp"`
// 	CurrNozzleTemp   int `json:"curr_nozzle_temp"`
// 	TargetHotbedTemp int `json:"target_hotbed_temp"`
// 	TargetNozzleTemp int `json:"target_nozzle_temp"`
// }

export interface MqttResponse {
  type: string;
  action: string;
  msgID: string;
  state: string;
  timestamp: number;
  code: number;
  msg: string;
  data: any;
}

export interface MqttTempatureData {
  currHotbedTemp: number;
  currNozzleTemp: number;
  targetHotbedTemp: number;
  targetNozzleTemp: number;
}

// Temperature interface that extends the MqttResponse interface
export interface Temperature extends MqttResponse {
  data: MqttTempatureData;
}
