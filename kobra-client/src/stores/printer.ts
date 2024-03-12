// import { defineStore } from 'pinia';
// import { ref } from 'vue';
// import { IFile, IPrintJob, PrinterState } from '@/interfaces/printer';

import { IPrintJob, PrinterState, IFile } from '@/interfaces/printer';
import { defineStore } from 'pinia';

export const usePrintStore = defineStore('printer', {
  // arrow function recommended for full type inference
  state: () => ({
    printJob: {} as IPrintJob,
    printStatus: {} as PrinterState,
    files: [] as IFile[],
    isUsbConnected: false as boolean,
  }),

  getters: {
    getFileList: (state) => {
      return state.files;
    },
  },
  actions: {
    async moveFileUp(file: IFile) {
      // Find file in list and move it to local by changing the path

      // Move file to local
      const response = await fetch(
        `/api/files/${file.path}/${file.name}/local`,
        {
          // Get
          method: 'GET',
        }
      );

      if (!response.ok) {
        console.error('Error moving file to local');
      }

      this.getFiles();
    },

    async moveFileDown(file: IFile) {
      // Find file in list and move it to usb by changing the path

      // Move file to usb
      const response = await fetch(`/api/files/${file.path}/${file.name}/usb`, {
        // Get
        method: 'GET',
      });

      if (!response.ok) {
        console.error('Error moving file to usb');
      }

      this.getFiles();
    },

    async deleteFile(file: IFile) {
      const response = await fetch(`/api/files/${file.path}/${file.name}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        console.error('Error deleting file');
      }

      // Delete it from the list
      this.getFiles();
    },

    downloadFile(file: IFile) {
      // Send filename to printer
      // Download file
      const link = document.createElement('a');
      link.style.display = 'none';
      link.href = `/api/files/${file.path}/${file.name}`;
      link.download = file.name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },

    async getFiles() {
      const response = await fetch('/api/files');
      const data: IFile[] = await response.json();

      this.files = data;
    },
  },
});
