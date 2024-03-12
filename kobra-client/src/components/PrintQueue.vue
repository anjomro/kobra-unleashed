<script setup lang="ts">
import { usePrintStore } from '@/stores/printer';
import { storeToRefs } from 'pinia';

const printStore = usePrintStore();

const { printJob, printStatus } = storeToRefs(printStore);

const cancelPrintJob = async () => {
  if (!printJob.value) {
    return;
  }
  const response = await fetch(`/api/print/${printJob.value.taskid}/cancel`, {
    method: 'POST',
  });
  if (!response.ok) {
    console.error('Error canceling print job');
  }
};

const pausePrintJob = async () => {
  if (!printJob.value) {
    return;
  }
  const response = await fetch(`/api/print/${printJob.value.taskid}/pause`, {
    method: 'POST',
  });
  if (!response.ok) {
    console.error('Error pausing print job');
  }
};

const resumePrintJob = async () => {
  if (!printJob.value) {
    return;
  }
  const response = await fetch(`/api/print/${printJob.value.taskid}/resume`, {
    method: 'POST',
  });
  if (!response.ok) {
    console.error('Error resuming print job');
  }
};
</script>

<template>
  <div
    class="w-full bg-neutral-200 dark:bg-neutral-700 p-4 rounded-lg flex flex-col gap-y-2"
    v-if="printJob?.filename"
  >
    <div class="flex justify-between items-center">
      <h2 class="text-lg font-bold">Print Queue</h2>
    </div>
    <div class="flex flex-col gap-y-2">
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-bold">Current Print Job</h3>
        <div class="flex gap-x-2">
          <button
            v-if="printStatus.state === 'printing'"
            class="btn btn-primary"
            @click="pausePrintJob"
          >
            Pause
          </button>
          <button
            v-if="printStatus.state === 'paused'"
            class="btn btn-primary"
            @click="resumePrintJob"
          >
            Resume
          </button>
          <button @click="cancelPrintJob" class="btn btn-danger">Cancel</button>
        </div>
      </div>
      <div class="flex flex-col gap-y-2">
        <div
          class="flex flex-col md:items-center md:flex-row gap-y-2 md:gap-x-2"
        >
          <p class="text-lg font-bold">File Name</p>
          <p>{{ printJob.filename }}</p>
        </div>
        <div class="flex justify-between items-center">
          <!-- div that moves with progrss -->
          <div class="bg-neutral-200 dark:bg-neutral-500 rounded-lg w-full">
            <div
              class="w-full bg-purple-500 rounded-lg flex justify-center items-center text-white font-semibold p-2"
              :style="{ width: printJob.progress + '%' }"
            >
              {{ printJob.progress }}%
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div
    class="w-full bg-neutral-200 dark:bg-neutral-700 p-4 rounded-lg flex flex-col gap-y-2"
    v-else
  >
    <h2 class="text-lg font-bold">Print Queue</h2>
    <p>No print job in queue</p>
  </div>
</template>

<style scoped></style>
