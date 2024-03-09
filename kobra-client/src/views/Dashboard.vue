<template>
  <div class="page p-4">
    <Teleport to="body">
      <FilesModal v-if="showFilesModal" @close="showFilesModal = false" />
      <EditParamPanel
        v-if="showEditParamPanel"
        @close="showEditParamPanel = false"
        :printerState="PrinterState"
      />
    </Teleport>
    <div
      class="flex items-center justify-between flex-col md:flex-row gap-y-2 md:gap-y-0"
    >
      <div>
        <h1 class="text-3xl font-bold">Kobra Unleashed</h1>
        <p>Welcome to your dashboard. Logged in as {{ userStore.username }}</p>
      </div>
      <div
        class="flex gap-x-2 flex-col md:flex-row gap-y-2 md:gap-y-0 w-full md:w-auto"
      >
        <button class="btn btn-primary icon">
          <PrintIcon class="w-8 h-8" />
          <p>New Print</p>
        </button>
        <button class="btn btn-primary icon" @click="showFilesModal = true">
          <FileIcon class="w-8 h-8" />
          <p>Files</p>
        </button>

        <button class="btn icon btn-hover-danger" @click="userStore.logout">
          <LogoutIcon class="w-8 h-8" />
          <p>Logout</p>
        </button>
      </div>
    </div>
    <!-- take all width. Only 1 col -->
    <div
      class="card-container cursor-pointer"
      @click="PrinterState ? (showEditParamPanel = true) : null"
    >
      <StatusCard
        title="Nozzle"
        :message="PrinterState.currentNozzleTemp?.toString().concat(' °C')"
        :displaysubmessage="true"
        :submessage="`Target: ${
          PrinterState.targetNozzleTemp?.toString().concat(' °C') ?? 'N/A'
        }`"
        :bgcolor="tempColor?.nozzle"
      />
      <StatusCard
        title="Hotbed"
        :message="PrinterState.currentBedTemp?.toString().concat(' °C')"
        :displaysubmessage="true"
        :submessage="`Target: ${
          PrinterState.targetBedTemp?.toString().concat(' °C') ?? 'N/A'
        }`"
        :bgcolor="tempColor?.bed"
      />
      <StatusCard
        title="Printer Status"
        :message="PrinterState.state"
        :bgcolor="tempColor?.status"
      />
      <StatusCard
        title="Speed Mode"
        :message="
          PrinterState.printSpeed === 1
            ? 'Slow'
            : PrinterState.printSpeed === 2
            ? 'Normal'
            : PrinterState.printSpeed === 3
            ? 'Fast'
            : PrinterState.printSpeed?.toString() ?? 'N/A'
        "
      />
      <StatusCard
        title="Fan Speed"
        :message="`${PrinterState.fanSpeed?.toString().concat('%') ?? 'N/A'}`"
        :bgcolor="tempColor?.fan"
      />
      <StatusCard
        title="Z Compensation"
        :message="PrinterState.zComp?.toString() ?? 'N/A'"
        :bgcolor="tempColor?.zComp"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { MqttResponse, PrintUpdate, Temperature } from '@/interfaces/mqtt';
import { useUserStore } from '@/stores/store';
import { ITempColor, PrinterState } from '@/interfaces/printer';
import { onMounted, ref, watchEffect, Teleport } from 'vue';
import StatusCard from '@/components/StatusCard.vue';
import EditParamPanel from '@/components/EditParamPanel.vue';
import LogoutIcon from '~icons/carbon/logout';
import FileIcon from '~icons/carbon/volume-file-storage';
import PrintIcon from '~icons/cbi/3dprinter-standby';
import FilesModal from '@/components/FilesModal.vue';

const userStore = useUserStore();

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

const tempColor = ref<ITempColor>({
  nozzle: '',
  bed: '',
  fan: '',
  status: '',
  zComp: '',
});

userStore.createWebSocket();

const ws = userStore.websock?.ws;

if (ws) {
  ws.onerror = (err) => {
    console.error('WebSocket Error:', err);
  };

  ws.onopen = () => {
    console.log('WebSocket Client Connected');
  };

  ws.addEventListener('message', (e) => {
    // return if pong
    if (e.data === 'pong') {
      console.log('Pong received');

      return;
    }
    const data = JSON.parse(e.data);
    const message = atob(data.message);

    const mqttResponse: MqttResponse = JSON.parse(message);

    if (
      mqttResponse.type === 'status' &&
      mqttResponse.action === 'workReport'
    ) {
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

      PrinterState.value.currentBedTemp = temp.data.curr_hotbed_temp;
      PrinterState.value.currentNozzleTemp = temp.data.curr_nozzle_temp;
      PrinterState.value.targetBedTemp = temp.data.settings.target_hotbed_temp;
      PrinterState.value.targetNozzleTemp =
        temp.data.settings.target_nozzle_temp;
      PrinterState.value.fanSpeed = temp.data.settings.fan_speed_pct;
      PrinterState.value.printSpeed = temp.data.settings.print_speed_mode;
      PrinterState.value.zComp = temp.data.settings.z_comp;
    }
  });
}

const showFilesModal = ref(false);
const showEditParamPanel = ref(false);

// Get username
onMounted(async () => {
  const response = await fetch('/api/user');
  const data = await response.json();
  // get username if ok
  if (response.ok) {
    userStore.$patch({ username: data.username });
  }

  watchEffect(() => {
    // Watch nozzle temp and change color from blue to green between 0 and target temp
    const fanSpeed = PrinterState.value.fanSpeed;

    if (fanSpeed !== undefined) {
      if (fanSpeed === 0) {
        // Set to blue
        tempColor.value.fan = 'rgb(0, 0, 100)';
      } else {
        // Calculate color for fan. from blue to green gradient
        const blueValue = 255 - Math.round((fanSpeed / 100) * 255); // Decreasing blue with fan speed
        const greenValue = Math.round((fanSpeed / 100) * 255); // Increasing green with fan speed

        // Start from blue to green, with values adjusted dynamically
        // tempColor.value.fan = `rgb(0, ${greenValue}, ${blueValue})`;

        // Clamp values to 100 because I don't want it to be so bright
        tempColor.value.fan = `rgb(0, ${Math.min(greenValue, 100)}, ${Math.min(
          blueValue,
          100
        )})`;
      }
    }
    // Nozzle. Blue to red gradient. If target temp is 0, set to blue
    const nozzleTemp = PrinterState.value.currentNozzleTemp;
    const targetNozzleTemp = PrinterState.value.targetNozzleTemp;

    if (nozzleTemp && targetNozzleTemp !== undefined) {
      if (targetNozzleTemp === 0) {
        tempColor.value.nozzle = 'rgb(0, 0, 100)'; // Set color to blue
      } else {
        // Calculate color for nozzle. from blue to red gradient

        const blueValue =
          255 - Math.round((nozzleTemp / targetNozzleTemp) * 255); // Decreasing blue with nozzle temp
        const redValue = Math.round((nozzleTemp / targetNozzleTemp) * 255); // Increasing red with nozzle temp

        // Clamp values to 100
        const clampedBlueValue = Math.min(blueValue, 100);
        const clampedRedValue = Math.min(redValue, 100);

        // Start from blue to red, with values adjusted dynamically
        tempColor.value.nozzle = `rgb(${clampedRedValue}, 0, ${clampedBlueValue})`;
      }
    }

    // Bed. Blue to red gradient. If target temp is 0, set to blue
    const bedTemp = PrinterState.value.currentBedTemp;
    const targetBedTemp = PrinterState.value.targetBedTemp;

    if (bedTemp && targetBedTemp !== undefined) {
      if (targetBedTemp === 0) {
        tempColor.value.bed = 'rgb(0, 0, 100)'; // Set color to blue
      } else {
        // Calculate color for bed. from blue to red gradient

        const blueValue = 255 - Math.round((bedTemp / targetBedTemp) * 255); // Decreasing blue with bed temp
        const redValue = Math.round((bedTemp / targetBedTemp) * 255); // Increasing red with bed temp

        // Clamp values to 100
        const clampedBlueValue = Math.min(blueValue, 100);
        const clampedRedValue = Math.min(redValue, 100);

        // Start from blue to red, with values adjusted dynamically
        tempColor.value.bed = `rgb(${clampedRedValue}, 0, ${clampedBlueValue})`;
      }
    }

    // Printer status. Green for free, yellow for printing, red for error
    const status = PrinterState.value.state;

    if (status) {
      switch (status) {
        case 'offline':
          tempColor.value.status = 'rgb(100, 0, 0)';
          break;
        case 'printing':
          tempColor.value.status = 'rgb(100, 100, 0)';
          break;
        case 'free':
          tempColor.value.status = 'rgb(0, 100, 0)';
          break;
        default:
          tempColor.value.status = 'rgb(100, 100, 100)';
      }
    }

    const zComp = PrinterState.value.zComp;

    // Gradient from red to green
    if (zComp !== undefined) {
      // If zComp is negative, make it red

      try {
        const parsedzComp = parseFloat(zComp);
        if (parsedzComp === 0) {
          // Set to blue
          tempColor.value.zComp = 'rgb(0, 0, 100)';
        } else {
          const greenValue = Math.min(255, Math.round((parsedzComp / 2) * 255));
          const redValue = Math.min(100, Math.round((-parsedzComp / 2) * 100));

          tempColor.value.zComp = `rgb(${redValue}, ${greenValue}, 0)`;
        }
      } catch (e) {
        console.error('Error parsing zComp:', e);
      }
    }
  });
});
</script>

<style lang="scss" scoped>
.card-container {
  @apply grid gap-4 mt-4;
  // grid-auto-flow: column; Do this but if too small make it stack vertically
  grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
}
</style>
