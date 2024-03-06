<script setup lang="ts">
import { MqttResponse } from '@/interfaces/mqtt';
import { ref } from 'vue';
const isDev = import.meta.env.DEV;

const printState = ref('NO_DATA');

const wsURL = isDev ? 'ws://localhost:3000/ws/info' : 'ws://localhost/ws/info';

const ws = new WebSocket(wsURL);

ws.onerror = (err) => {
  console.error('WebSocket Error:', err);
};

ws.onopen = () => {
  console.log('WebSocket Client Connected');
};

ws.onmessage = (e) => {
  const data = JSON.parse(e.data);
  const message = atob(data.message);

  const mqttResponse: MqttResponse = JSON.parse(message);
  console.log('Received:', mqttResponse);

  if (mqttResponse.type === 'status' && mqttResponse.action === 'workReport') {
    printState.value = mqttResponse.state;
  }
};
</script>

<template>
  <div class="page">
    <h1>Kobra Unleashed</h1>
    <p>Welcome to your dashboard</p>

    <div>
      <p>Status: {{ printState }}</p>
      <h2></h2>
    </div>
  </div>
</template>

<style lang="scss" scoped></style>
