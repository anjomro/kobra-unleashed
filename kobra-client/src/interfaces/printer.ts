import { MqttFileListRecord } from '@/interfaces/mqtt';

export interface PrinterState {
  state: string;
  currentNozzleTemp: number | undefined;
  currentBedTemp: number | undefined;
  targetNozzleTemp: number | undefined;
  targetBedTemp: number | undefined;
  fanSpeed: number | undefined;
  printSpeed: number | undefined;
  zComp: string | undefined;
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
