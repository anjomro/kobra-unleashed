import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

interface IWebSocket {
  ws: WebSocket;
  pingInterval: number | null;
}

export const useUserStore = defineStore('user', () => {
  const auth = useStorage<boolean>('auth', false);
  const authExpiryDate = useStorage<number>('authExpiryDate', 0);
  const websock = ref<IWebSocket | null>(null);
  const username = ref('N/A');
  const router = useRouter();
  const isDev = import.meta.env.DEV;

  // Make onlogout callback that takes in a websocket and closes it

  const createWebSocket = () => {
    const wsURL = isDev
      ? 'ws://localhost:3000/ws/info'
      : 'ws://localhost/ws/info';
    const ws = new WebSocket(wsURL);
    registerWebSocket(ws);
  };

  const registerWebSocket = (ws: WebSocket) => {
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
        registerWebSocket(new WebSocket('ws://localhost:8080'));
      }
    });

    // Set the websocket and pingInterval
    websock.value = {
      ws,
      pingInterval,
    };
  };

  async function logout() {
    // Disconnect from ws server

    if (websock.value) {
      if (websock.value.pingInterval) {
        clearInterval(websock.value.pingInterval);
      }

      websock.value.ws.close();
      console.log('Websocket closed');
    }

    await fetch('/api/logout', {
      method: 'POST',
      credentials: 'include',
    });

    auth.value = false;
    authExpiryDate.value = 0;

    // Delete cookies
    document.cookie =
      'session_id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';

    // Redirect to login
    router.replace({ name: 'Login', query: { logout: 'true' } });
  }

  return {
    auth,
    authExpiryDate,
    username,
    logout,
    registerWebSocket,
    websock,
    createWebSocket,
  };
});
