<template>
  <div class="card bordered w-full md:w-2/3 lg:w-1/2">
    <div class="overflow-x-auto">
      <table class="table">
        <tbody>
        <!-- row 1 -->
        <tr v-for="file in (printer && printer.files ? printer.files[1] : [])" :key="file" class="hover">
          <td>
            <div class="flex items-center gap-3">
              <div>
                <div class="font-bold">{{ file.filename }}</div>
                <div class="text-sm opacity-50">{{ displayFileSize(file.size) }}</div>

              </div>
            </div>
          </td>
          <th>
            <button @click="printFile(file)" class="btn btn-ghost btn-xs">Print</button>
          </th>
        </tr>
        </tbody>

      </table>
    </div>
  </div>
</template>

<script>

import socket from "../socket.js";

export default {
  props: ['printer'],
  methods: {
    printFile(file) {
      socket.emit('print_file', {printerId: this.printer.id, file: file.filename});
    },

    displayFileSize(size) {
      // If size is less than 1 MB, display in KB
      // Otherwise, display in MB
      if (size < 1024 * 1024) {
        size = size / 1024;
        size = Number((size).toFixed(2));
        return `${size} KB`;
      } else {
        size = size / (1024 * 1024);
        size = Number((size).toFixed(2));
        return `${size} MB`;
      }
    },
  }
};
</script>