import { defineStore } from 'pinia';
import { useStorage } from '@vueuse/core';

export const useUserStore = defineStore('user', () => {
  const auth = useStorage<boolean>('auth', false);
  const authExpiryDate = useStorage<number>('authExpiryDate', 0);

  async function logout(callback: Function) {
    const response = await fetch('/api/logout', {
      method: 'POST',
      credentials: 'include',
    });

    if (response.ok && response.status === 200) {
      auth.value = false;
      authExpiryDate.value = 0;
      callback({ name: 'Login' });
    } else {
      console.error('Failed to log out', response);
    }
  }

  return { auth, authExpiryDate, logout };
});
