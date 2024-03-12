<script setup lang="ts">
import { useUserStore } from '@/stores/store';
import { computed, onBeforeMount, onBeforeUnmount, onMounted, ref } from 'vue';
import FilesInspectModal from '@/components/FilesInspectModal.vue';
import CloseIcon from '~icons/carbon/close-large';
import ModeUpIcon from '~icons/system-uicons/pull-up';
import ModeDownIcon from '~icons/system-uicons/pull-down';
import DeleteIcon from '~icons/ph/trash-fill';

import { usePrintStore } from '@/stores/printer';
import { IFile } from '@/interfaces/printer';
import { convertSize, convertTimestamp } from '@/utils/utils';
import PrintIcon from '~icons/cbi/3dprinter-standby';
import DownloadIcon from '~icons/ph/download-fill';
import ViewIcon from '~icons/ph/eye-fill';
import ImagePreview from './ImagePreview.vue';

const showInspectModal = ref(false);

const userStore = useUserStore();
const printStore = usePrintStore();
const fileList = computed(() => printStore.getFileList);
const ws = computed(() => userStore.websock);
const isUsbConnected = computed(() => printStore.isUsbConnected);
const selectedFile = ref<IFile | null>(null);
const emit = defineEmits(['close', 'print']);

const handleDragOver = (e: DragEvent) => {
  e.preventDefault();
  //only allow gcode files
  if (e.dataTransfer?.items[0].type === 'text/x-gcode') {
    e.dataTransfer.dropEffect = 'copy';
  } else {
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'none';
    }
  }
};

const handleDrop = async (e: DragEvent, path: string) => {
  e.preventDefault();
  const file = e.dataTransfer?.files[0];
  if (file) {
    // only allow gcode files
    if (file.type === 'text/x-gcode') {
      const formData = new FormData();
      formData.append('file', file);
      const response = await fetch(`/api/files/${path}`, {
        method: 'POST',
        body: formData,
      });
      if (response.ok) {
        printStore.getFiles();
      } else {
        console.error('Failed to upload file');
      }
    }
  }
};

onBeforeMount(() => {
  printStore.getFiles();
});

onBeforeUnmount(() => {
  ws.value?.ws?.send(
    JSON.stringify({
      action: 'stopWatchUSB',
    })
  );
});

onMounted(async () => {
  ws.value?.ws?.send(
    JSON.stringify({
      action: 'watchUSB',
    })
  );
});
</script>

<template>
  <Teleport to="body">
    <FilesInspectModal
      v-if="showInspectModal && selectedFile"
      :file="selectedFile"
      @close="showInspectModal = false"
    />
  </Teleport>

  <div
    class="w-full h-full absolute inset-0 z-10 backdrop-blur-sm flex items-center justify-center"
    ref="modal"
    @click.self="emit('close')"
  >
    <div
      class="w-full h-full md:w-[75%] md:rounded-lg md:h-[90%] bg-neutral-200 dark:bg-neutral-700 p-4 flex flex-col gap-y-2 absolute overflow-hidden"
    >
      <div class="titlebar">
        <h1 class="text-3xl font-bold title">Files</h1>
        <button
          class="btn btn-primary btn-hover-danger close place-self-end"
          @click="emit('close')"
        >
          <CloseIcon class="w-6 h-6" />
        </button>
      </div>
      <div
        v-if="!fileList.length"
        class="flex items-center justify-center w-full h-full"
      >
        <div class="spinner"></div>
        <p class="text-lg font-bold ml-2">Loading...</p>
      </div>
      <div v-else class="flex flex-col gap-y-2 w-full h-full overflow-y-auto">
        <div
          class="flex flex-col bg-neutral-200 dark:bg-neutral-800 p-2 rounded-lg gap-y-2"
          @dragover.prevent="handleDragOver"
          @drop.prevent="handleDrop($event, 'local')"
        >
          <h2 class="text-lg font-bold p-2">Local</h2>

          <li
            v-for="localfile in fileList.filter((f) => f.path === 'local')"
            :key="localfile.name"
            class="files-container"
          >
            <ImagePreview :file="localfile" />
            <h2 class="text-lg font-bold">{{ localfile.name }}</h2>
            <p class="text-sm">{{ convertSize(localfile.size) }}</p>
            <p class="text-sm">{{ convertTimestamp(localfile.modified_at) }}</p>
            <div class="files-container">
              <div class="flex gap-x-2">
                <button
                  class="btn btn-primary"
                  @click="printStore.printFile(localfile)"
                >
                  <PrintIcon class="w-6 h-6" />
                </button>
                <button
                  class="btn btn-primary"
                  @click="
                    showInspectModal = true;
                    selectedFile = localfile;
                  "
                >
                  <ViewIcon class="w-6 h-6" />
                </button>
                <button
                  class="btn btn-primary"
                  :disabled="!isUsbConnected"
                  @click="printStore.moveFileDown(localfile)"
                >
                  <ModeDownIcon class="w-6 h-6" />
                </button>
                <button
                  class="btn btn-primary"
                  @click="printStore.downloadFile(localfile)"
                >
                  <DownloadIcon class="w-6 h-6" />
                </button>
                <button
                  class="btn btn-hover-danger"
                  @click="printStore.deleteFile(localfile)"
                >
                  <DeleteIcon class="w-6 h-6" />
                </button>
              </div>
            </div>
          </li>
          <li
            v-if="!fileList.filter((f) => f.path === 'local').length"
            class="flex items-center justify-center w-full h-full"
          >
            <p class="text-lg font-bold">No files found</p>
          </li>
        </div>
        <div
          class="flex flex-col bg-neutral-200 dark:bg-neutral-800 p-2 rounded-lg gap-y-2"
          @dragover.prevent="handleDragOver"
          @drop.prevent="handleDrop($event, 'sdcard')"
        >
          <h2 class="text-lg font-bold p-2">USB</h2>

          <li
            v-for="usbfile in fileList.filter((f) => f.path === 'usb')"
            :key="usbfile.name"
            class="files-container"
          >
            <ImagePreview :file="usbfile" />
            <h2 class="text-lg font-bold">{{ usbfile.name }}</h2>
            <p class="text-sm">{{ convertSize(usbfile.size) }}</p>
            <p class="text-sm">{{ convertTimestamp(usbfile.modified_at) }}</p>

            <div class="flex gap-x-2">
              <button
                class="btn btn-primary"
                @click="printStore.printFile(usbfile)"
              >
                <PrintIcon class="w-6 h-6" />
              </button>
              <button
                class="btn btn-primary"
                @click="
                  showInspectModal = true;
                  selectedFile = usbfile;
                "
              >
                <ViewIcon class="w-6 h-6" />
              </button>
              <button
                class="btn btn-primary"
                @click="printStore.moveFileUp(usbfile)"
              >
                <ModeUpIcon class="w-6 h-6" />
              </button>
              <button
                class="btn btn-primary"
                @click="printStore.downloadFile(usbfile)"
              >
                <DownloadIcon class="w-6 h-6" />
              </button>
              <button
                class="btn btn-hover-danger"
                @click="printStore.deleteFile(usbfile)"
              >
                <DeleteIcon class="w-6 h-6" />
              </button>
            </div>
          </li>
          <li
            v-if="
              !fileList.filter((f) => f.path === 'usb').length && isUsbConnected
            "
            class="flex items-center justify-center w-full h-full"
          >
            <p class="text-lg font-bold">No files found</p>
          </li>
          <li
            class="flex items-center justify-center w-full h-full"
            v-else-if="!isUsbConnected"
          >
            <p class="text-lg font-bold">No USB connected</p>
          </li>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.titlebar {
  @apply grid place-items-center;
  grid-template-areas: '. title close';
  grid-template-columns: 1fr 2fr 1fr;
  .title {
    grid-area: title;
  }

  .close {
    grid-area: close;
  }
}

.files-container {
  @apply flex flex-col md:flex-row items-center justify-between gap-y-2 md:gap-x-2 bg-neutral-100 dark:bg-neutral-700 p-2 rounded-lg;
}
</style>
