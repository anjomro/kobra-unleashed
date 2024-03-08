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
  curr_hotbed_temp: number;
  curr_nozzle_temp: number;
  target_hotbed_temp: number;
  target_nozzle_temp: number;
}

export interface MqttPrintUpdateData {
  taskid: string;
  localtask: string;
  curr_hotbed_temp: number;
  curr_nozzle_temp: number;
  settings: {
    target_nozzle_temp: number;
    target_hotbed_temp: number;
    fan_speed_pct: number;
    print_speed_mode: number;
    z_comp: number;
  };
}

// Temperature interface that extends the MqttResponse interface
export interface Temperature extends MqttResponse {
  data: MqttTempatureData;
}

export interface PrintUpdate extends MqttResponse {
  data: MqttPrintUpdateData;
}
