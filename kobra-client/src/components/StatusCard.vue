<script setup lang="ts">
const props = defineProps({
  title: { type: String, required: true },
  message: String,
  submessage: String,
  displaysubmessage: { type: Boolean, default: true },
  bgcolor: String,
});
</script>

<template>
  <div
    class="status-card text-white shadow-md rounded-xl p-4 dark:bg-neutral-600 bg-neutral-500"
    :style="{ backgroundColor: props.bgcolor ?? '' }"
  >
    <p>{{ props.title }}</p>
    <h2 class="font-bold text-2xl capitalize">
      {{ props.message ?? 'N/A' }}
    </h2>
    <h3 v-if="displaysubmessage" class="text-sm">
      {{ props.submessage }}
    </h3>
  </div>
</template>

<style lang="scss" scoped></style>
