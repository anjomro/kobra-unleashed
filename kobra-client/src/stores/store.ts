import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  const auth = useStorage<boolean>('auth', false);
  const authExpiryDate = useStorage<number>('authExpiryDate', 0);

  // Make onlogout callback that takes in a websocket and closes it

  const webSockets = ref<WebSocket[]>([]);

  const registerWebSocket = (ws: WebSocket) => {
    webSockets.value.push(ws);
  };

  async function logout(callback: Function) {
    // Disconnect from ws server

    if (webSockets.value.length > 0) {
      webSockets.value.forEach((ws) => {
        console.log('Closing websocket', ws.url);
        ws.close();
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

  return { auth, authExpiryDate, logout, registerWebSocket };
});
