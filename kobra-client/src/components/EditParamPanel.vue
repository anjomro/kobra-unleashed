<script setup lang="ts">
import { PrinterState } from '@/interfaces/printer';
import { usePrintStore } from '@/stores/printer';
import { computed, ref, watchEffect } from 'vue';
import CloseIcon from '~icons/icon-park-solid/back';

const emit = defineEmits(['close']);

const printState = usePrintStore();

interface IPrinter {
  taskid: string;
  target_nozzle_temp?: number;
  target_hotbed_temp?: number;
  fan_speed_pct?: number;
  print_speed_mode?: number;
  z_comp?: string;
}

const currentPrinterState = printState.printStatus;

const newPrinterState = ref<PrinterState>({});

const saveNewParams = async (newState: PrinterState) => {
  const tempPrinter: IPrinter = {
    taskid: '0',
  };

  // Only set if not empty
  if (newState.targetNozzleTemp !== undefined) {
    tempPrinter.target_nozzle_temp = newState.targetNozzleTemp;
  } else if (newState.targetBedTemp !== undefined) {
    tempPrinter.target_hotbed_temp = newState.targetBedTemp;
  } else if (newState.fanSpeed !== undefined) {
    tempPrinter.fan_speed_pct = newState.fanSpeed;
  } else if (newState.printSpeed !== undefined) {
    tempPrinter.print_speed_mode = newState.printSpeed;
  } else if (newState.zComp !== undefined) {
    tempPrinter.z_comp = newState.zComp;
  } else {
    return;
  }

  const response = await fetch('/api/printer/settings', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(tempPrinter),
  });

  if (response.ok) {
    console.log('Settings updated');
  } else {
    console.log('Settings update failed');
  }
};

const buttonDisable = ref(false);
const currentPrintSpeed = computed(() => {
  switch (currentPrinterState?.printSpeed) {
    case 1:
      return 'Slow';
    case 2:
      return 'Normal';
    case 3:
      return 'Fast';
    default:
      return 'Unchanged';
  }
});

watchEffect(() => {
  buttonDisable.value = Object.values(newPrinterState.value).every(
    (value) => value === undefined || value === ''
  );
});
</script>

<template>
  <div
    class="z-20 absolute inset-0 flex items-center justify-center backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <div
      class="bg-neutral-200 dark:bg-neutral-700 p-4 rounded-lg flex flex-col gap-y-2 w-full h-full md:w-[65%] md:h-[90%] overflow-hidden"
    >
      <div class="title">
        <h1 class="header text-3xl font-semibold">Edit Parameters</h1>
        <button
          class="btn btn-primary btn-hover-danger self-end close justify-self-end"
          @click="emit('close')"
        >
          <CloseIcon class="w-6 h-6" />
        </button>
      </div>

      <div class="flex flex-col flex-1 overflow-y-auto">
        <div class="flex flex-col gap-y-4">
          <form
            class="flex flex-col gap-y-2"
            @submit.prevent="saveNewParams(newPrinterState)"
          >
            <div class="flex flex-col gap-y-2">
              <label for="targetNozzleTemp"
                >Target Nozzle Temp ({{
                  currentPrinterState?.targetNozzleTemp
                }})</label
              >
              <input
                v-model="newPrinterState.targetNozzleTemp"
                type="number"
                id="targetNozzleTemp"
                min="-1"
                class="input"
              />
            </div>
            <div class="flex flex-col gap-y-2">
              <label for="targetBedTemp"
                >Target Bed Temp ({{
                  currentPrinterState?.targetBedTemp
                }})</label
              >
              <input
                v-model="newPrinterState.targetBedTemp"
                type="number"
                id="targetBedTemp"
                min="-1"
                class="input"
              />
            </div>
            <div class="flex flex-col gap-y-2">
              <label for="fanSpeed"
                >Fan Speed ({{ currentPrinterState?.fanSpeed }}%)</label
              >
              <input
                v-model="newPrinterState.fanSpeed"
                type="number"
                id="fanSpeed"
                min="-1"
                max="100"
                class="input"
              />
            </div>
            <div class="flex flex-col gap-y-2">
              <!-- Dropdown -->
              <label for="printSpeed"
                >Print Speed ({{ currentPrintSpeed }})</label
              >
              <select
                v-model.number="newPrinterState.printSpeed"
                id="printSpeed"
                class="input"
              >
                <option value="">Unchanged</option>
                <option value="1">Slow</option>
                <option value="2">Normal</option>
                <option value="3">Fast</option>
              </select>
            </div>
            <div class="flex flex-col gap-y-2">
              <label for="zComp"
                >Z Compensation ({{ currentPrinterState?.zComp }})</label
              >
              <input
                v-model="newPrinterState.zComp"
                type="number"
                id="zComp"
                class="input"
              />
            </div>
            <p class="text-sm text-neutral-500 dark:text-neutral-400">
              *Leave empty if no change
            </p>
            <button
              class="btn btn-primary self-start"
              :disabled="buttonDisable"
            >
              Save
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.title {
  @apply grid items-center place-items-center;
  grid-template-areas: '. title close';
  grid-template-columns: 1fr 2fr 1fr;
}

.close {
  grid-area: close;
}

.header {
  grid-area: title;
}

.input {
  @apply w-full p-2 rounded-lg border bg-neutral-200 dark:bg-neutral-600;
}
</style>
