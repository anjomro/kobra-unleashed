<template>
  <button class="btn btn-circle btn-outline btn-sm" onclick="name_modal.showModal()">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
    </svg>
  </button>
  <dialog id="name_modal" class="modal">
    <div class="modal-box">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
      </form>
      <h3 class="font-bold text-lg">Set New Printer Name</h3>
      <p class="py-4">ID: {{ printer.id }}</p>
      <form @submit="savePrinterName" method="dialog" class="join flex">
        <input v-model="newPrinterName" class="input input-bordered join-item flex-1">
        <button class="btn join-item">Save</button>
      </form>
    </div>
  </dialog>
</template>


<script>
export default {
  props: ['printer'],
  data() {
    return {
      newPrinterName: '',
      // ...
    };
  },
  methods: {
    savePrinterName() {
      this.$emit('set-printer-name', this.printer.id, this.newPrinterName);
    },
  }
};
</script>


<style scoped>

</style>