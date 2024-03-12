<script setup lang="ts">
import { usePrintStore } from '@/stores/printer';
import { storeToRefs } from 'pinia';
import { WebGLPreview, init } from 'gcode-preview';

import { onBeforeUnmount, onMounted, ref, watchEffect } from 'vue';

const printStore = usePrintStore();

const gcodePreview = ref<HTMLCanvasElement | undefined>(undefined);
const displayCanvas = ref(false);

// Make a function that calculates a color based on the progress red to green
const calculateColor = (progress: number) => {
  const r =
    progress < 50 ? 255 : Math.floor(255 - ((progress * 2 - 100) * 255) / 100);
  const g = progress > 50 ? 255 : Math.floor((progress * 2 * 255) / 100);
  return `rgb(${r}, ${g}, 0)`;
};

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

let preview: WebGLPreview | undefined = undefined;

const handleResize = () => {
  if (preview) {
    preview.resize();
  }
};

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});

onMounted(() => {
  watchEffect(async () => {
    // return if no print job progress or curr_layer return

    if (printJob.value?.filename) {
      window.addEventListener('resize', handleResize);
      const response = await fetch(
        `/api/files/local/${printJob.value.filename}`
      );
      if (response.ok) {
        const gcode = await response.text();

        if (
          printJob.value.progress === undefined ||
          printJob.value.curr_layer === undefined
        ) {
          return;
        }

        preview = init({
          canvas: gcodePreview.value,
          extrusionColor: calculateColor(printJob.value.progress),
          initialCameraPosition: [200, 200, 200],
          renderTubes: true,
          buildVolume: {
            x: 410,
            y: 400,
            z: 400,
            r: 400,
            i: 400,
            j: 400,
          },
        });

        preview.endLayer = printJob.value.curr_layer;

        preview.processGCode(gcode);

        preview.render();

        console.log('ok');

        displayCanvas.value = true;
      }
    }
  });
});
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
        <canvas
          v-show="displayCanvas"
          ref="gcodePreview"
          class="w-full h-[20rem] rounded-lg"
        ></canvas>
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
async
