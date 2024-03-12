import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

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
    wsURL: isDev ? 'ws://localhost:3000/ws/info' : 'ws://localhost/ws/info',
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
        if (e.code !== 1000) {
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
    async logout() {
      // Disconnect from ws server

      if (this.websock) {
        if (this.websock.pingInterval) {
          clearInterval(this.websock.pingInterval);
        }

        this.websock.ws.close();
        console.log('Websocket closed');
      }

      await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include',
      });

      this.auth = false;
      this.authExpiryDate = 0;

      // Delete cookies
      document.cookie =
        'session_id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
      const router = useRouter();

      // Redirect to login
      router.replace({ name: 'Login', query: { logout: 'true' } });
    },
  },
});
