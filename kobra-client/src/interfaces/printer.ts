import { MqttFileListRecord } from '@/interfaces/mqtt';

export interface PrinterState {
  state?: string;
  currentNozzleTemp?: number;
  currentBedTemp?: number;
  targetNozzleTemp?: number;
  targetBedTemp?: number;
  fanSpeed?: number;
  printSpeed?: number;
  zComp?: string;
}

export interface IFileList {
  records: MqttFileListRecord[];
  listType: string;
}

export interface ITempColor {
  // Nozzle and bed temp
  nozzle: string;
  bed: string;
  // Fan speed
  fan: string;
  status: string;
  zComp: string;
}
