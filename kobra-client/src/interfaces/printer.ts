export interface PrinterState {
  state: string;
  currentNozzleTemp: number | undefined;
  currentBedTemp: number | undefined;
  targetNozzleTemp: number | undefined;
  targetBedTemp: number | undefined;
  fanSpeed: number | undefined;
  printSpeed: string | undefined;
  zComp: string | undefined;
}
