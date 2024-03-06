<script setup lang="ts">
import { ref } from 'vue';
import { useUserStore } from '@/stores/store';
import { useRouter } from 'vue-router';
import { LoginResponse } from '@/interfaces/loginResponse';

const username = ref('');
const password = ref('');
const loginError = ref('');
const userStore = useUserStore();
const router = useRouter();

const doLogin = async () => {
  loginError.value = '';
  //   Use fetch to send a POST request to the server on /api/login
  //   with the username and password as the body

  const response = await fetch('/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Required for cookies to be sent with the request
    body: JSON.stringify({
      username: username.value,
      password: password.value,
    }),
  });

  if (response.status === 200 && response.ok) {
    //   Redirect to the dashboard
    // Get "expires" from the response and set it as the authExpiryDate

    // Decode json response
    const data: LoginResponse = await response.json();

    userStore.$patch({ auth: true, authExpiryDate: data.expires });
    router.push({ name: 'Dashboard' });
  } else {
    loginError.value = 'Invalid username or password';
  }
};
</script>

<template>
  <div class="page grid place-items-center">
    <div class="self-end text-center">
      <h1 class="text-3xl">Kobra Unleashed</h1>
      <p class="text-gray-500">Login to your dashboard</p>
    </div>

    <div>
      <form @submit.prevent="doLogin">
        <input
          type="text"
          placeholder="Username"
          v-model="username"
          autocomplete="username"
        />
        <input
          type="password"
          placeholder="Password"
          v-model="password"
          autocomplete="current-password"
        />
        <button type="submit">Login</button>
      </form>
      <p v-if="loginError" class="text-red-500">{{ loginError }}</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.page {
  grid-template-columns: 1fr;
  grid-template-rows: auto auto auto;
}

form {
  @apply flex flex-col gap-y-2;

  input {
    @apply p-2 border border-gray-300 dark:bg-transparent rounded-md;
  }

  button {
    @apply p-2 bg-blue-500 text-white rounded-md;
  }
}
</style>
