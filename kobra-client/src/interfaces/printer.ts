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

// {
// 	"files": [
// 		{
// 			"name": "Calibration Cube-C0.2-19m-2024-3-10.gcode",
// 			"size": 389980,
// 			"modified_at": "2024-03-11T03:43:53.47356424+08:00",
// 			"path": "local"
// 		},
// 		{
// 			"name": "Rounded-Mini-Scraper-v01-C0.2-9m-2024-3-10.gcode",
// 			"size": 899394,
// 			"modified_at": "2024-03-11T03:34:40+08:00",
// 			"path": "usb"
// 		}
// 	]
// }

export interface IFile {
  name: string;
  size: number;
  modified_at: string;
  path: string;
}

export interface IPrinterFiles {
  files: IFile[];
}
