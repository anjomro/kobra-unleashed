import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';
import { ref } from 'vue';

interface IWebSocket {
  ws: WebSocket;
  pingInterval: number | null;
}

const isDev = import.meta.env.DEV;

export const useUserStore = defineStore('user', {
  // arrow function recommended for full type inference
  state: () => ({
    auth: useStorage<boolean>('auth', false),
    authExpiryDate: useStorage<number>('authExpiryDate', 0),
    websock: ref<IWebSocket | null>(null),
    username: ref('N/A'),
    wsURL: isDev
      ? 'ws://localhost:3000/ws/info'
      : `ws://${location.host}/ws/info`,
  }),
  actions: {
    createWebSocket() {
      const ws = new WebSocket(this.wsURL);
      this.registerWebSocket(ws);
    },
    registerWebSocket(ws: WebSocket) {
      // Ping server every 20 seconds to keep connection alive
      const pingInterval = setInterval(() => {
        if (ws.readyState === ws.OPEN) {
          ws.send('ping');
        } else {
          clearInterval(pingInterval);
        }
        console.log('Ping sent');
      }, 20000);

      ws.addEventListener('close', (e) => {
        // Reconnect if the connection is closed for unexpected reasons
        if (e.code !== 1006) {
          console.log('Websocket closed unexpectedly, reconnecting');
          if (this.auth) this.registerWebSocket(new WebSocket(this.wsURL));
        }
      });

      // Set the websocket and pingInterval
      this.websock = {
        ws,
        pingInterval,
      };
    },
  },
});
