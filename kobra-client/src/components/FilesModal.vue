<script setup lang="ts">
import { FileList, MqttFileListRecord, MqttResponse } from '@/interfaces/mqtt';
import { useUserStore } from '@/stores/store';
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { DateTime, DateTimeFormatOptions } from 'luxon';
import FilesInspectModal from '@/components/FilesInspectModal.vue';
import CloseIcon from '~icons/carbon/close-large';
import DownloadButton from '~icons/carbon/download';
import ViewIcon from '~icons/carbon/view';
import PrintIcon from '~icons/mdi/printer-3d-nozzle';
import MoveUpIcon from '~icons/system-uicons/pull-up';
import MoveDownIcon from '~icons/system-uicons/pull-down';
import DeleteIcon from '~icons/ic/baseline-delete';

interface IFileList {
  records: MqttFileListRecord[];
  listType: string;
}

const showInspectModal = ref(false);
const isUSBConnected = ref(false);

const fileList = ref<IFileList[]>([]);

const ws = useUserStore().websock?.ws;

const messageHandler = (event: MessageEvent) => {
  if (event.data === 'pong') return;

  const data = JSON.parse(event.data);
  const message = atob(data.message);

  // Parse json
  const dataJson = JSON.parse(message);

  if (dataJson['usb_connected'] !== undefined) {
    isUSBConnected.value = dataJson['usb_connected'];
  }

  const mqttResponse = dataJson as MqttResponse;

  if (
    (mqttResponse.type === 'file' && mqttResponse.action === 'listLocal') ||
    mqttResponse.action === 'listUdisk'
  ) {
    const fileMqttResponse: FileList = mqttResponse;

    fileMqttResponse.data.records.forEach((record) => {
      record.file_location = fileMqttResponse.action;
    });

    if (fileList.value.length > 2) {
      fileList.value = [];
    }

    // If listType already exists, remove it and push the new one
    // Replace the old list with the new one
    const index = fileList.value.findIndex(
      (list) => list.listType === fileMqttResponse.action
    );

    if (index !== -1) {
      fileList.value.splice(index, 1);
    }

    fileList.value?.push({
      listType: fileMqttResponse.action,
      records: fileMqttResponse.data.records,
    });

    // Sort listtype Make sure listLocal is first
    fileList.value.sort((a, b) => {
      if (a.listType === 'listLocal') {
        return -1;
      }
      if (b.listType === 'listLocal') {
        return 1;
      }
      return 0;
    });
  }
};

ws?.addEventListener('message', messageHandler);

onBeforeUnmount(() => {
  ws?.removeEventListener('message', messageHandler);
});

const convertSize = (size: number) => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];

  let i = 0;
  while (size >= 1024) {
    size /= 1024;
    i++;
  }

  return `${size.toFixed()} ${units[i]}`;
};

const convertTimestamp = (timestamp: number) => {
  try {
    // Get browser locale
    const locale = navigator.language as string;

    // Convert timestamp to date object
    const dateTime = DateTime.fromMillis(timestamp);

    // Extract time using browser locale
    const time = dateTime.toLocaleString(
      {
        hour: 'numeric',
        minute: 'numeric',
        second: 'numeric',
      },
      {
        locale,
      }
    );

    // Format date example: "April 1, 2022"
    const format: DateTimeFormatOptions = {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    };
    // Combine formatted date and time with a space
    return `${dateTime.toLocaleString(format)} ${time}`;
  } catch (error) {
    // Handle errors
    console.error('Error formatting timestamp:', error);
    return 'Invalid timestamp';
  }
};

const selectedFile = ref<MqttFileListRecord | null>(null);

onMounted(async () => {
  await fetch('/api/files?pathType=listLocal');
  await fetch('/api/files?pathType=listUdisk');

  ws?.send(
    JSON.stringify({
      action: 'check-usb',
    })
  );
});

const emit = defineEmits(['close']);

const downloadFile = async (file: MqttFileListRecord) => {
  const response = await fetch(
    `/api/files/${file.file_location}/${file.filename}`
  );
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.style.display = 'none';
  a.href = url;
  a.download = file.filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  a.remove();
};
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
      <div class="flex justify-between">
        <h1 class="text-3xl font-bold">Files</h1>
        <button class="btn btn-primary btn-hover-danger" @click="emit('close')">
          <CloseIcon class="w-6 h-6" />
        </button>
      </div>
      <p>Here are your files</p>
      <div class="flex flex-col h-[90%] overflow-y-auto gap-y-2">
        <div
          v-for="(list, index) in fileList"
          :key="index"
          class="bg-neutral-100 dark:bg-neutral-800 p-4 rounded-lg flex flex-col space-y-2 flex-1"
          v-if="fileList"
        >
          <h2
            class="text-xl font-bold"
            v-if="list.listType === 'listUdisk' && isUSBConnected === false"
          >
            No USB connected
          </h2>
          <h2
            class="text-xl font-bold"
            v-else-if="list.listType === 'listLocal'"
          >
            Local
          </h2>
          <h2
            class="text-xl font-bold"
            v-else-if="list.listType === 'listUdisk'"
          >
            USB
          </h2>

          <ul class="flex flex-col space-y-2">
            <li
              v-for="record in list.records.filter(
                (record) => !record.is_dir && record.size > 0
              )"
              :key="record.filename"
              class="p-2 bg-neutral-200 dark:bg-neutral-700 rounded-lg flex flex-col md:flex-row justify-between items-center gap-y-2"
            >
              <div>
                <p class="font-semibold">{{ record.filename }}</p>
                <p>Size: {{ convertSize(record.size) }}</p>
                <p>{{ convertTimestamp(record.timestamp) }}</p>
              </div>
              <div class="flex gap-x-1">
                <button class="btn btn-primary">
                  <PrintIcon class="w-6 h-6" />
                </button>
                <button
                  class="btn btn-primary"
                  @click="
                    selectedFile = record;
                    showInspectModal = true;
                  "
                >
                  <ViewIcon class="w-6 h-6" />
                </button>
                <button
                  v-if="
                    list.listType === 'listLocal' ||
                    (list.listType === 'listUdisk' && isUSBConnected)
                  "
                  class="btn btn-primary"
                  @click="
                    list.listType === 'listLocal'
                      ? ws?.send(
                          JSON.stringify({
                            action: 'moveToUdisk',
                            filename: record.filename,
                          })
                        )
                      : ws?.send(
                          JSON.stringify({
                            action: 'moveToLocal',
                            filename: record.filename,
                          })
                        )
                  "
                >
                  <MoveDownIcon
                    class="w-6 h-6"
                    v-if="list.listType === 'listLocal'"
                  />
                  <MoveUpIcon class="w-6 h-6" v-else />
                </button>
                <button class="btn btn-primary">
                  <DownloadButton
                    class="w-6 h-6"
                    @click="downloadFile(record)"
                  />
                </button>
                <button
                  class="btn btn-hover-danger"
                  @click="
                    ws?.send(
                      JSON.stringify({
                        action: 'deleteFile',
                        filename: record.filename,
                        filelocation: list.listType,
                      })
                    )
                  "
                >
                  <DeleteIcon class="w-6 h-6" />
                </button>
              </div>
            </li>
            <p v-if="list.records.length === 0 && isUSBConnected">
              No files found
            </p>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
