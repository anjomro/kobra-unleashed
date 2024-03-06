<script setup lang="ts">
import StatusCard from '@/components/StatusCard.vue';
import { MqttResponse, Temperature } from '@/interfaces/mqtt';
import { useStorage } from '@vueuse/core';

const isDev = import.meta.env.DEV;

interface PrinterState {
  state: string;
  currentNozzleTemp: number | undefined;
  currentBedTemp: number | undefined;
  targetNozzleTemp: number | undefined;
  targetBedTemp: number | undefined;
}

// const printState = ref<PrinterState>({
//   state: 'offline',
//   currentNozzleTemp: undefined,
//   currentBedTemp: undefined,
//   targetNozzleTemp: undefined,
//   targetBedTemp: undefined,
// });

const PrinterState = useStorage<PrinterState>('printer-state', {
  state: 'offline',
  currentNozzleTemp: undefined,
  currentBedTemp: undefined,
  targetNozzleTemp: undefined,
  targetBedTemp: undefined,
});

const wsURL = isDev ? 'ws://localhost:3000/ws/info' : 'ws://localhost/ws/info';

const ws = new WebSocket(wsURL);

ws.onerror = (err) => {
  console.error('WebSocket Error:', err);
};

ws.onopen = () => {
  console.log('WebSocket Client Connected');
};

ws.onmessage = (e) => {
  const data = JSON.parse(e.data);
  const message = atob(data.message);

  const mqttResponse: MqttResponse = JSON.parse(message);
  console.log('Received:', mqttResponse);

  if (mqttResponse.type === 'status' && mqttResponse.action === 'workReport') {
    PrinterState.value.state = mqttResponse.state;
  }

  if (mqttResponse.type === 'tempature' && mqttResponse.action === 'report') {
    // Set mqttResponse to Temperature interface
    const temp: Temperature = mqttResponse;
    PrinterState.value.currentBedTemp = temp.data.curr_hotbed_temp;
    PrinterState.value.currentNozzleTemp = temp.data.curr_nozzle_temp;
    PrinterState.value.targetBedTemp = temp.data.target_hotbed_temp;
    PrinterState.value.targetNozzleTemp = temp.data.target_nozzle_temp;
  }
};
</script>

<template>
  <div class="page p-4">
    <div class="flex items-center justify-between">
      <div>
        <h1>Kobra Unleashed</h1>
        <p>Welcome to your dashboard</p>
      </div>
      <RouterLink to="/logout">Logout</RouterLink>
    </div>

    <div class="flex flex-wrap gap-4">
      <StatusCard
        title="Nozzle"
        :message="PrinterState.currentNozzleTemp?.toString()"
      />
      <StatusCard
        title="Hotbed"
        :message="PrinterState.currentBedTemp?.toString()"
      />
      <StatusCard
        title="Printer Status"
        :message="PrinterState.state"
        :displaysubmessage="false"
      />
      <h2></h2>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
