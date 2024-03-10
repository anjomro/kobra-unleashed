import { defineStore } from 'pinia';
import { ref } from 'vue';
import {
  IFile,
  IPrintJob,
  IPrinterFiles,
  PrinterState,
} from '@/interfaces/printer';

export const usePrintStore = defineStore('printer', () => {
  const printJob = ref<IPrintJob>({});
  const printStatus = ref<PrinterState | null>(null);
  const files = ref<IPrinterFiles | null>(null);
  const isUsbConnected = ref<boolean>(false);

  const printFile = (file: IFile | undefined) => {
    // Send filename to printer
  };

  const moveFileUp = async (file: IFile | undefined) => {
    // Move file to local
    const response = await fetch(
      `/api/files/${file?.path}/${file?.name}/local`,
      {
        // Get
        method: 'GET',
      }
    );

    if (response.ok) {
      console.log('File moved to local');
    } else {
      console.error('Error moving file to local');
    }

    getFiles();
  };

  const moveFileDown = async (file: IFile | undefined) => {
    // Move file to usb
    const response = await fetch(`/api/files/${file?.path}/${file?.name}/usb`, {
      // Get
      method: 'GET',
    });

    if (response.ok) {
      console.log('File moved to usb');
    } else {
      console.error('Error moving file to usb');
    }

    getFiles();
  };

  const deleteFile = (file: IFile | undefined) => {
    // Send filename to printer
  };

  const downloadFile = (file: IFile | undefined) => {
    // Send filename to printer
    // Download file
    const link = document.createElement('a');
    link.style.display = 'none';
    link.href = `/api/files/${file?.path}/${file?.name}`;
    link.download = file?.name || '';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const getFiles = async () => {
    const response = await fetch('/api/files');
    const data: IPrinterFiles = await response.json();
    files.value = data;
  };
  return {
    printJob,
    printStatus,
    files,
    printFile,
    moveFileUp,
    moveFileDown,
    deleteFile,
    getFiles,
    downloadFile,
    isUsbConnected,
  };
});
