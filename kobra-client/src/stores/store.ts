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

  // Make onlogout callback that takes in a websocket and closes it

  const registerWebSocket = (ws: WebSocket) => {
    // Ping server every 20 seconds to keep connection alive
    const pingInterval = setInterval(() => {
      ws.send('ping');
      console.log('Ping sent');
    }, 20000);

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
  };
});
