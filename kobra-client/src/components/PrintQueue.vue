<script setup lang="ts">
import { MqttResponse } from '@/interfaces/mqtt';
import { usePrintStore } from '@/stores/printer';
import { useUserStore } from '@/stores/store';
import { computed } from 'vue';

const printStore = usePrintStore();

const printJob = computed(() => printStore.printJob);

const ws = useUserStore().websock?.ws;

ws?.addEventListener('message', (e) => {
  if (e.data === 'pong') {
    return;
  }
  const data = JSON.parse(e.data);
  const message = atob(data.message);

  const mqttResponse: MqttResponse = JSON.parse(message);

  if (mqttResponse.type === 'print') {
    if (mqttResponse.action) {
      printStore.$patch({ printStatus: mqttResponse.data });
    }
  }
});

const cancelPrintJob = async () => {
  const response = await fetch(`/api/print/${printJob.value.taskid}/cancel`, {
    method: 'POST',
  });
  if (response.ok) {
    printStore.$patch({ printStatus: {} });
  }
};
</script>

<template>
  <div
    class="w-full bg-neutral-200 dark:bg-neutral-700 p-4 rounded-lg flex flex-col gap-y-2"
    v-if="printJob.taskid"
  >
    <div class="flex justify-between items-center">
      <h2 class="text-lg font-bold">Print Queue</h2>
    </div>
    <div class="flex flex-col gap-y-2">
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-bold">Current Print Job</h3>
        <button
          @click="cancelPrintJob"
          class="p-2 rounded-md bg-primary-500 text-white"
        >
          Cancel
        </button>
      </div>
      <div class="flex flex-col gap-y-2">
        <div class="flex justify-between items-center">
          <p class="text-lg font-bold">File Name</p>
          <p>{{ printJob.filename }}</p>
        </div>
        <div class="flex justify-between items-center">
          <p class="text-lg font-bold">Progress</p>
          <p>{{ printJob.progress }}%</p>
        </div>
        <div class="flex justify-between items-center">
          <p class="text-lg font-bold">Time Remaining</p>
          <p>{{ printJob.remain_time }}</p>
        </div>
        <div class="flex justify-between items-center">
          <p class="text-lg font-bold">Time Elapsed</p>
          <p>{{ printJob.progress }}</p>
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
