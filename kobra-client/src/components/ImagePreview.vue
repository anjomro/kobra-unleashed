<template>
  <img
    v-if="image"
    :src="image"
    alt="Image preview"
    class="w-[100%] h-[100%] md:h-14 md:w-14 object-cover"
  />
</template>

<script setup lang="ts">
import { IFile } from '@/interfaces/printer';
import { PropType, onBeforeMount, ref } from 'vue';

const image = ref<string>('');

const props = defineProps({
  file: {
    type: Object as PropType<IFile>,
    required: true,
  },
});

onBeforeMount(async () => {
  const file = props.file;
  const response = await fetch(`/api/files/${file.path}/${file.name}`);
  const data = await response.text();

  const lines = data.split('\n');
  // Find the line that starts with '; thumbnail begin'
  // Skit it
  let i = 0;
  for (; i < lines.length; i++) {
    if (lines[i].startsWith('; thumbnail begin')) {
      break;
    }
  }

  // Start from the next line
  i++;
  // Create a new array to store the base64 string
  const base64 = [];
  // Iterate through the lines
  for (; i < lines.length; i++) {
    // If the line starts with '; thumbnail end', break the loop
    if (lines[i].startsWith('; thumbnail end')) {
      break;
    }
    // Otherwise, push the line to the base64 array
    const currentLine = lines[i];
    // Remove the leading and trailing whitespaces and semicolons
    base64.push(currentLine.trim().replace(';', ''));
  }

  // Join the base64 array to create the base64 string
  const base64String = base64.join('');
  // Set the base64 string to the image ref
  image.value = `data:image/jpeg;base64,${base64String}`;
});
</script>
