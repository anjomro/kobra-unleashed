<template>
  <div class="p-4">
    <div class="navbar bg-base-100">
      <div class="w-full md:w-2/3 lg:w-1/2 justify-between flex gap-5">

        <div class="flex-1">
          <a class="btn btn-ghost text-xl" href="https://github.com/anjomro/kobra-unleashed" target="_blank">Kobra Unleashed</a>
        </div>

        <div class="flex-none gap-2">
          <PrinterNaming v-if="selectedPrinter" :printer="selectedPrinter" @set-printer-name="addPrinterName"></PrinterNaming>
        </div>
        <div class="flex-none">
          <select v-model="selectedPrinterId" @change="selectPrinterId" class="select select-bordered w-full max-w-xs">
            <option disabled="disabled" selected="selected" v-if="Object.keys(printers).length == 0">No Printer
              available
            </option>
            <option v-for="printerId in Object.keys(printers)" :key="printerId" :value="printerId">
              {{ printerNames[printerId] || printerId }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <PrinterDetails v-if="selectedPrinter" :printer="selectedPrinter" class="mt-4"/>
    <FileList v-if="selectedPrinter" :printer="selectedPrinter" class="mt-4"/>
  </div>
</template>
<script>
import PrinterDetails from './PrinterDetails.vue';
import FileList from './FileList.vue';
import PrinterNaming from './PrinterNaming.vue';
import axios from "axios";
import socket from "../socket.js";
import VueCookies from 'vue-cookies';

export default {
  components: {
    PrinterNaming,
    PrinterDetails,
    FileList,
  },
  data() {
    return {
      printers: {}, // dict of printers id -> printer
      selectedPrinter: null,
      selectedPrinterId: null,
      isPrinting: false,
      printerNames: {},
    };
  },
  methods: {
    getPrinterIdList() {
      // make a list of printer ids from printers array
      return this.printers.map(printer => printer.id);
    },
    selectPrinterId() {
      this.selectedPrinter = this.printers[this.selectedPrinterId];
    },

    updatePrinter(printerDict) {
      console.log('Updating printer:', printerDict.id);
      // updates a printer from a dict {id: printerId, printer: { ... printer data ... }}
      this.printers[printerDict.id] = printerDict.printer;
      // update selected printer if it is the changed printer
      if (this.selectedPrinterId === printerDict.id) {
        this.selectedPrinter = printerDict.printer;
      }
    },

    loadPrinterNames() {
      return VueCookies.get('printerNames') || {};
    },
    addPrinterName(id, name) {
      this.loadPrinterNames();
      this.printerNames[id] = name;
      VueCookies.set('printerNames', this.printerNames);
    },
    resolvePrinterName(id) {
      if (!this.printerNames) {
        this.printerNames = this.loadPrinterNames();
      }
      return this.printerNames[id] || id;
    },
  },
  created() {
    console.log('Dashboard created');
    this.printerNames = this.loadPrinterNames();
    socket.on('connect', () => {
      console.log('Connected to server');
    });
    socket.on('disconnect', () => {
      console.log('Disconnected from server');
    });
    socket.on('error', (error) => {
      console.error('Error:', error);
    });
    socket.on('printer_updated', (data) => {
      // data: {id: printerId, printer: { ... printer data ... }}
      console.log('Printer state changed: ', data);
      this.updatePrinter(data);

      if (this.selectedPrinterId === undefined || this.selectedPrinterId === null) {
        this.selectedPrinterId = data.id;
        this.selectPrinterId();
      }
    });
    socket.on('printer_list', (data) => {
      this.printers = data;
      console.log('Printer list:', data);

      if (this.selectedPrinterId === undefined || this.selectedPrinterId === null) {
        this.selectedPrinterId = Object.keys(data)[0];
        this.selectPrinterId();
      }
    });

    socket.on('files_updated', (data) => {
      console.log('File state changed');
      this.updatePrinter(data);
    });

    // Fetch printers from socket.io API
    socket.emit('get_printer_list');

  }
};
</script>