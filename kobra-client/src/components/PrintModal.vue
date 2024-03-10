<script setup lang="ts">
import { ref } from 'vue';
import CloseIcon from '~icons/carbon/close-large';
import PrintIcon from '~icons/cbi/3dprinter-standby';

const emit = defineEmits(['close', 'print']);

const fileInput = ref<HTMLInputElement | null>(null);

const selectedFile = ref<File | null>(null);

const handleDragOver = (e: DragEvent) => {
  e.preventDefault();

  //   Only allow .gcode text/x-gcode and show cursor
  if (e.dataTransfer?.items[0].type === 'text/x-gcode') {
    e.dataTransfer.dropEffect = 'copy';
  } else {
    e.dataTransfer ? (e.dataTransfer.dropEffect = 'none') : null;
  }
};

const handleDrop = (e: DragEvent) => {
  e.preventDefault();

  //   Only allow .gcode text/x-gcode
  if (e.dataTransfer?.items[0].type === 'text/x-gcode') {
    selectedFile.value = e.dataTransfer?.files[0];
  } else {
    selectedFile.value = null;
  }
};

const openFileInput = () => {
  fileInput.value?.click();
};

const handleFileInput = (e: Event) => {
  const target = e.target as HTMLInputElement;
  selectedFile.value = target.files?.[0] ?? null;
};
</script>

<template>
  <div
    class="z-20 absolute inset-0 flex items-center justify-center backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <div
      class="bg-neutral-200 dark:bg-neutral-700 p-4 rounded-lg grid grid-rows-3 gap-4 w-full h-full md:w-[65%] md:h-[90%] overflow-hidden"
    >
      <div class="title">
        <h1 class="header text-3xl font-semibold">New Print</h1>
        <button
          class="btn btn-primary btn-hover-danger close justify-self-end place-self-start"
          @click="emit('close')"
        >
          <CloseIcon class="w-6 h-6" />
        </button>
      </div>

      <!-- File drop -->
      <div
        class="flex flex-col flex-1 items-center justify-center border-2 border-dashed rounded-lg cursor-pointer bg-neutral-200 dark:bg-neutral-800"
        @dragover.prevent="handleDragOver"
        @drop.prevent="handleDrop"
        @click="openFileInput"
      >
        <input
          type="file"
          ref="fileInput"
          class="hidden"
          @change="handleFileInput"
          accept=".gcode"
        />
        <div v-if="!selectedFile">
          <h2 class="text-2xl font-semibold">Drop your gcode file here</h2>
          <p class="text-neutral-500 dark:text-neutral-400">
            or click to select a file
          </p>
        </div>
        <div v-else>
          <h2 class="text-2xl font-semibold">{{ selectedFile.name }}</h2>
          <p class="text-neutral-500 dark:text-neutral-400">
            Click to change file
          </p>
        </div>
      </div>

      <button
        class="btn btn-primary icon self-center justify-self-center"
        :disabled="!selectedFile"
        @click="
          emit('close');
          emit('print', selectedFile);
        "
      >
        <PrintIcon class="w-6 h-6" />
        Print
      </button>
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
</style>
