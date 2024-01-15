<template>
  <div class="card w-full md:w-2/3 lg:w-1/2 bg-base-100 shadow-xl">
    <div class="stats stats-horizontal shadow">
      <div class="stat">
        <div class="stat-title">Nozzle</div>
        <div class="stat-value">{{ printer.nozzle_temp }} °C</div>
        <div class="stat-desc">Target: {{ printer.target_nozzle_temp }} °C</div>
      </div>

      <div class="stat">
        <div class="stat-title">Hotbed</div>
        <div class="stat-value">{{ printer.hotbed_temp }} °C</div>
        <div class="stat-desc">Target: {{ printer.target_hotbed_temp }} °C</div>
      </div>
      <div class="stat">
        <div class="stat-title">Status</div>
        <div class="stat-value">{{ printer.state }}</div>
      </div>

    </div>
    <div class="card-body">
      <div v-if="isPrinting() || isPaused()" class="card card-bordered card-compact w-1/2">

        <div class="overflow-x-auto">
          <table class="table">
            <tbody>
            <tr>
              <th>File</th>
              <td>{{ printer.print_job.file }}</td>
            </tr>
            <tr>
              <th>Progress</th>
              <td>{{ printer.print_job.progress }}%</td>
            </tr>
            <tr>
              <th>Layers</th>
              <td>{{ printer.print_job.curr_layer }}/{{ printer.print_job.total_layers }}</td>
            </tr>
            <tr>
              <th>Print Time</th>
              <td>{{ printer.print_job.print_time.toFixed(2) }} min</td>
            </tr>
            <tr>
              <th>Time Remaining</th>
              <td>{{ printer.print_job.remaining_time.toFixed(2) }} min</td>
            </tr>
            <tr>
              <th>Material Used</th>
              <td>{{ printer.print_job.supplies_usage.toFixed(2) }} mm</td>
            </tr>
            <tr>
              <th>Print Speed</th>
              <td>{{ printSpeedMode() }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="card-actions justify-end">
        <div class="join">

          <button @click="selectFile" class="btn btn-primary join-item">Select File</button>

          <button v-if="!selectedFile" @click="upload_and_print" class="btn btn-primary join-item">Upload/Print</button>
          <div v-else class="tooltip tooltip-open tooltip-top"
               v-bind:data-tip="selectedFile.name">
            <button @click="upload_and_print" class="btn btn-primary join-item tooltip-open">Upload/Print</button>
          </div>
          <input type="file" ref="fileInput" class="hidden" @change="onFileSelected"/>
          <button @click="pausePrint" class="btn btn-primary join-item" v-if="isPrinting() && !isPaused()">Pause
          </button>
          <button @click="resumePrint" class="btn btn-primary join-item" v-if="isPrinting() && isPaused()">Resume
          </button>
          <button @click="stopPrint" class="btn btn-primary join-item" v-if="isPrinting()">Stop</button>
          <button @click="fetchFiles" class="btn btn-accent join-item">Fetch Files</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import socket from "../socket.js";

export default {
  props: ['printer'],
  data() {
    return {
      fanSpeed: 0,
      selectedFile: null,
    }
  },
  methods: {
    printSpeedMode() {
      // Return the print speed mode in human readable format
      let speed_mode = {
        1: "Unknown",
        1: "Slow",
        2: "Normal",
        3: "Fast",
      };
      if (this.printer.print_job.print_speed_mode in speed_mode) {
        return speed_mode[this.printer.print_job.print_speed_mode]
      } else {
        return "Unknown"
      }
    },
    selectFile() {
      this.$refs.fileInput.click();
    },
    onFileSelected(e) {
      this.selectedFile = e.target.files[0];
    },
    async upload_and_print() {
      if (this.selectedFile) {
        const formData = new FormData();
        formData.append('file', this.selectedFile);
        formData.append('printer_id', this.printer.id);
        try {
          const response = await axios.post('http://0.0.0.0:5000/api/print', formData);
          console.log('File uploaded:', response.data);
          this.selectedFile = null;
        } catch (error) {
          console.error('Error uploading file:', error);
        }
      }
    },
    isPrinting() {
      // return true if status one of: "busy", "printing", "preheating"
      return ["busy", "printing", "preheating", "paused"].includes(this.printer.state);
    },
    isPaused() {
      // return true if status one of: "paused"
      return this.printer.state === "paused"
    },
    setFanSpeed(speed) {
      socket.emit('set_fan', {"id": this.printer.id, "speed": this.fanSpeed});
    },
    pausePrint() {
      socket.emit('pause_print', {"id": this.printer.id});
    },
    resumePrint() {
      socket.emit('resume_print', {"id": this.printer.id});
    },
    stopPrint() {
      socket.emit('stop_print', {"id": this.printer.id});
    },
    async fetchFiles() {
      try {
        const response = await axios.get(`http://0.0.0.0:5000/api/printer/${this.printer.id}/files`);
        console.log(response.data);
      } catch (error) {
        console.error('Error fetching files:', error);
      }
    }
  }
};
</script>