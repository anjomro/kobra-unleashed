import { defineStore } from 'pinia';
import { ref } from 'vue';
import { IFile, IPrintJob, PrinterState } from '@/interfaces/printer';

export const usePrintStore = defineStore('printer', () => {
  const printJob = ref<IPrintJob>({});
  const printStatus = ref<PrinterState | null>(null);
  const files = ref<IFile[]>([]);
  const isUsbConnected = ref<boolean>(false);

  const printFile = (file: IFile | undefined) => {
    console.log('Printing file', file);
  };

  const moveFileUp = async (file: IFile) => {
    // Find file in list and move it to local by changing the path

    // Move file to local
    const response = await fetch(`/api/files/${file.path}/${file.name}/local`, {
      // Get
      method: 'GET',
    });

    if (!response.ok) {
      console.error('Error moving file to local');
    }

    getFiles();
  };

  const moveFileDown = async (file: IFile) => {
    // Find file in list and move it to usb by changing the path

    // Move file to usb
    const response = await fetch(`/api/files/${file.path}/${file.name}/usb`, {
      // Get
      method: 'GET',
    });

    if (!response.ok) {
      console.error('Error moving file to usb');
    }

    getFiles();
  };

  const deleteFile = async (file: IFile) => {
    const response = await fetch(`/api/files/${file.path}/${file.name}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      console.error('Error deleting file');
    }

    // Delete it from the list
    getFiles();
  };

  const downloadFile = (file: IFile) => {
    // Send filename to printer
    // Download file
    const link = document.createElement('a');
    link.style.display = 'none';
    link.href = `/api/files/${file.path}/${file.name}`;
    link.download = file.name;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const getFiles = async () => {
    const response = await fetch('/api/files');
    const data: IFile[] = await response.json();

    files.value = data;
  };

  return {
    printJob,
    printStatus,
    files,
    isUsbConnected,
    printFile,
    moveFileUp,
    moveFileDown,
    deleteFile,
    getFiles,
    downloadFile,
  };
});
