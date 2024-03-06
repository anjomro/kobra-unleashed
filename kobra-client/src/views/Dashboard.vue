<script setup lang="ts">
import StatusCard from '@/components/StatusCard.vue';
import { MqttResponse, PrintUpdate, Temperature } from '@/interfaces/mqtt';
import { useUserStore } from '@/stores/store';
import { ref } from 'vue';

const isDev = import.meta.env.DEV;

const userStore = useUserStore();

interface PrinterState {
  state: string;
  currentNozzleTemp: number | undefined;
  currentBedTemp: number | undefined;
  targetNozzleTemp: number | undefined;
  targetBedTemp: number | undefined;
  fanSpeed: number | undefined;
  printSpeed: string | undefined;
  zComp: string | undefined;
}

const PrinterState = ref<PrinterState>({
  state: 'offline',
  currentNozzleTemp: undefined,
  currentBedTemp: undefined,
  targetNozzleTemp: undefined,
  targetBedTemp: undefined,
  fanSpeed: undefined,
  printSpeed: undefined,
  zComp: undefined,
});

const wsURL = isDev ? 'ws://localhost:3000/ws/info' : 'ws://localhost/ws/info';

const ws = new WebSocket(wsURL);

userStore.registerWebSocket(ws);

ws.onerror = (err) => {
  console.error('WebSocket Error:', err);
};

ws.onopen = () => {
  console.log('WebSocket Client Connected');
};

ws.onmessage = (e) => {
  // return if pong
  if (e.data === 'pong') {
    console.log('Pong received');

    return;
  }
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

  if (mqttResponse.type === 'print' && mqttResponse.action === 'update') {
    const temp: PrintUpdate = mqttResponse;
    console.log('Print Update:', temp);

    PrinterState.value.currentBedTemp = temp.data.curr_hotbed_temp;
    PrinterState.value.currentNozzleTemp = temp.data.curr_nozzle_temp;
    PrinterState.value.targetBedTemp = temp.data.settings.target_hotbed_temp;
    PrinterState.value.targetNozzleTemp = temp.data.settings.target_nozzle_temp;
    PrinterState.value.fanSpeed = temp.data.settings.fan_speed_pct;

    switch (temp.data.settings.print_speed_mode) {
      case 1:
        PrinterState.value.printSpeed = 'Slow';
        break;
      case 2:
        PrinterState.value.printSpeed = 'Normal';
        break;
      case 3:
        PrinterState.value.printSpeed = 'Fast';
        break;
      default:
        PrinterState.value.printSpeed = 'N/A';
    }

    // PrinterState.value.printSpeed = temp.data.settings.print_speed_mode;
    PrinterState.value.zComp = temp.data.settings.z_comp;
  }
};
</script>

<template>
  <div class="page p-4">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">Kobra Unleashed</h1>
        <p>Welcome to your dashboard</p>
      </div>
      <RouterLink to="/logout">Logout</RouterLink>
    </div>
    <!-- take all width. Only 1 col -->
    <div class="card-container">
      <StatusCard
        title="Nozzle"
        :message="PrinterState.currentNozzleTemp?.toString().concat(' °C')"
        :submessage="`Target: ${
          PrinterState.targetNozzleTemp?.toString().concat(' °C') ?? 'N/A'
        }`"
      />
      <StatusCard
        title="Hotbed"
        :message="PrinterState.currentBedTemp?.toString().concat(' °C')"
        :submessage="`Target: ${
          PrinterState.targetBedTemp?.toString().concat(' °C') ?? 'N/A'
        }`"
      />
      <StatusCard
        title="Printer Status"
        :message="PrinterState.state"
        :displaysubmessage="false"
      />
      <StatusCard
        title="Speed Mode"
        :message="PrinterState.printSpeed?.toString() ?? 'N/A'"
      />
      <StatusCard
        title="Fan Speed"
        :message="`${PrinterState.fanSpeed?.toString().concat('%') ?? 'N/A'}`"
      />
      <StatusCard
        title="Z Compensation"
        :message="PrinterState.zComp?.toString() ?? 'N/A'"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.card-container {
  @apply grid gap-4 mt-4;
  // grid-auto-flow: column; Do this but if too small make it stack vertically
  grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
}
</style>
