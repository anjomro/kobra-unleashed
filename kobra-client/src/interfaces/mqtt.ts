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

interface MqttTempatureData {
  curr_hotbed_temp: number;
  curr_nozzle_temp: number;
  target_hotbed_temp: number;
  target_nozzle_temp: number;
}

interface MqttPrintUpdateData {
  taskid: string;
  filename: string;
  print_time: number;
  progress: number;
  supplies_usage: number;
  total_layers: number;
  curr_layer: number;
  remain_time: number;
}

// Temperature interface that extends the MqttResponse interface
export interface Temperature extends MqttResponse {
  data: MqttTempatureData;
}

export interface PrintUpdate extends MqttResponse {
  data: MqttPrintUpdateData;
}
