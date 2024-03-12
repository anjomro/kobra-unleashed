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

export interface ITempColor {
  // Nozzle and bed temp
  nozzle: string;
  bed: string;
  // Fan speed
  fan: string;
  status: string;
  zComp: string;
}

export interface IPrintJob {
  taskid?: string;
  remain_time?: number;
  filename?: string;
  print_time?: number;
  progress?: number;
  supplies_usage?: number;
  total_layers?: number;
  curr_layer?: number;
}

export interface IFile {
  name: string;
  size: number;
  modified_at: string;
  path: string;
}

export interface IPrinterFiles {
  files: IFile[];
}
