<script setup lang="ts">
import { MqttFileListRecord } from '@/interfaces/mqtt';
import { onMounted, ref, toRefs } from 'vue';
import CloseIcon from '~icons/icon-park-solid/back';
import { EditorView, basicSetup } from 'codemirror';

const prop = defineProps({
  file: {
    type: Object as () => MqttFileListRecord,
    required: true,
  },
});

const props = toRefs(prop);
const editorRef = ref<HTMLDivElement | null>(null);
const emit = defineEmits(['close']);
const data = ref('');

onMounted(async () => {
  const response = await fetch(
    `/api/files/${props.file.value.file_location}/${props.file.value.filename}`
  );
  data.value = await response.text();

  const editor = new EditorView({
    extensions: [basicSetup],
    parent: editorRef.value!,
  });

  editor.dispatch({
    changes: {
      from: 0,
      to: 0,
      insert: data.value,
    },
  });
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
      <div class="flex items-center justify-between">
        <div class="flex">
          <h1 class="text-lg font-semibold">{{ props.file.value.filename }}</h1>
        </div>
        <button
          class="btn btn-primary btn-hover-danger self-end"
          @click="emit('close')"
        >
          <CloseIcon class="w-6 h-6" />
        </button>
      </div>
      <div class="flex flex-col h-[90%] overflow-y-auto gap-y-2">
        <p>File content</p>
        <div class="h-full" ref="editorRef"></div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
