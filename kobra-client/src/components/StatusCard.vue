<script setup lang="ts">
const props = defineProps({
  title: { type: String, required: true },
  message: String,
  submessage: String,
  displaysubmessage: { type: Boolean, default: true },
});
</script>

<template>
  <div class="status-card">
    <p>{{ props.title }}</p>
    <h2 class="font-bold text-2xl capitalize">
      {{ props.message ?? 'N/A' }}
    </h2>
    <h3 v-if="displaysubmessage" class="text-sm">
      {{ props.submessage }}
    </h3>
  </div>
</template>

<style lang="scss" scoped>
.status-card {
  @apply shadow-md rounded-lg p-4 dark:bg-neutral-600 bg-neutral-200;
}
</style>
