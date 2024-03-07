import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  const auth = useStorage<boolean>('auth', false);
  const authExpiryDate = useStorage<number>('authExpiryDate', 0);
  const username = ref('N/A');

  // Make onlogout callback that takes in a websocket and closes it

  interface IWebSocket {
    client: WebSocket;
    pingInterval: number;
  }

  const webSockets = ref<IWebSocket[]>([]);

  const registerWebSocket = (ws: WebSocket) => {
    // Ping server every 20 seconds to keep connection alive
    const pingInterval = setInterval(() => {
      ws.send('ping');
      console.log('Ping sent');
    }, 20000);

    webSockets.value.push({ client: ws, pingInterval });
  };

  async function logout(callback: Function) {
    // Disconnect from ws server

    if (webSockets.value.length > 0) {
      webSockets.value.forEach((ws) => {
        console.log('Closing websocket', ws.client.url);
        ws.client.close();
        clearInterval(ws.pingInterval);
      });
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

    callback({ name: 'Login' });
  }

  return { auth, authExpiryDate, username, logout, registerWebSocket };
});
