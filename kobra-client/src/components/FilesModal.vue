<script setup lang="ts">
import { useUserStore } from '@/stores/store';
import { computed, onBeforeMount, onBeforeUnmount, onMounted, ref } from 'vue';
import FilesInspectModal from '@/components/FilesInspectModal.vue';
import CloseIcon from '~icons/carbon/close-large';
import ModeUpIcon from '~icons/system-uicons/pull-up';
import ModeDownIcon from '~icons/system-uicons/pull-down';

import { usePrintStore } from '@/stores/printer';
import { IFile } from '@/interfaces/printer';
import { convertSize, convertTimestamp } from '@/utils/utils';
import PrintIcon from '~icons/cbi/3dprinter-standby';
import DownloadIcon from '~icons/ph/download-fill';
import ViewIcon from '~icons/ph/eye-fill';

const showInspectModal = ref(false);

const ws = useUserStore().websock?.ws;
const printStore = usePrintStore();
const fileList = computed(() => printStore.files);
const isUsbConnected = computed(() => printStore.isUsbConnected);

onBeforeMount(() => {
  printStore.getFiles();
});

onBeforeUnmount(() => {
  ws?.send(
    JSON.stringify({
      action: 'stopWatchUSB',
    })
  );
});

const selectedFile = ref<IFile | null>(null);

const emit = defineEmits(['close', 'print']);

onMounted(async () => {
  ws?.send(
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
        v-if="!fileList?.files.length"
        class="flex items-center justify-center w-full h-full"
      >
        <div class="spinner"></div>
        <p class="text-lg font-bold ml-2">Loading...</p>
      </div>
      <div v-else class="flex flex-col gap-y-2 w-full h-full overflow-y-auto">
        <div
          class="flex flex-col bg-neutral-200 dark:bg-neutral-800 p-2 rounded-lg"
        >
          <h2 class="text-lg font-bold p-2">Local</h2>
          <ul
            class="flex flex-col gap-y-2 bg-neutral-100 dark:bg-neutral-700 p-2 rounded-lg"
          >
            <li
              v-for="file in fileList?.files.filter((f) => f.path === 'local')"
            >
              <h2 class="text-lg font-bold">{{ file.name }}</h2>
              <p class="text-sm">{{ convertSize(file.size) }}</p>

              <div class="flex justify-between items-center">
                <p class="text-sm">{{ convertTimestamp(file.modified_at) }}</p>
                <div class="flex gap-x-2">
                  <button class="btn btn-primary">
                    <PrintIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    @click="
                      showInspectModal = true;
                      selectedFile = file;
                    "
                  >
                    <ViewIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    :disabled="!isUsbConnected"
                    @click="printStore.moveFileDown(file)"
                  >
                    <ModeDownIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    @click="printStore.downloadFile(file)"
                  >
                    <DownloadIcon class="w-6 h-6" />
                  </button>
                </div>
              </div>
            </li>
            <li
              v-if="!fileList.files.filter((f) => f.path === 'local').length"
              class="flex items-center justify-center w-full h-full"
            >
              <p class="text-lg font-bold">No files found</p>
            </li>
          </ul>
        </div>
        <div
          class="flex flex-col bg-neutral-200 dark:bg-neutral-800 p-2 rounded-lg"
        >
          <h2 class="text-lg font-bold p-2">USB</h2>
          <ul
            class="flex flex-col gap-y-2 bg-neutral-100 dark:bg-neutral-700 p-2 rounded-lg"
          >
            <li v-for="file in fileList?.files.filter((f) => f.path === 'usb')">
              <h2 class="text-lg font-bold">{{ file.name }}</h2>
              <p class="text-sm">{{ convertSize(file.size) }}</p>

              <div class="flex justify-between items-center">
                <p class="text-sm">{{ convertTimestamp(file.modified_at) }}</p>
                <div class="flex gap-x-2">
                  <button class="btn btn-primary">
                    <PrintIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    @click="
                      showInspectModal = true;
                      selectedFile = file;
                    "
                  >
                    <ViewIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    @click="printStore.moveFileUp(file)"
                  >
                    <ModeUpIcon class="w-6 h-6" />
                  </button>
                  <button
                    class="btn btn-primary"
                    @click="printStore.downloadFile(file)"
                  >
                    <DownloadIcon class="w-6 h-6" />
                  </button>
                </div>
              </div>
            </li>
            <li
              v-if="
                !fileList.files.filter((f) => f.path === 'usb').length &&
                isUsbConnected
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
          </ul>
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
</style>
