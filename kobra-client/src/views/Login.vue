<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useUserStore } from '@/stores/store';
import { useRouter } from 'vue-router';
import { LoginResponse } from '@/interfaces/loginResponse';
import LoginIcon from '~icons/carbon/login';

const username = ref('');
const password = ref('');
const loginError = ref('');
const remember = ref(false);
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
      remember: remember.value,
    }),
  });

  if (response.status === 200 && response.ok) {
    //   Redirect to the dashboard
    // Get "expires" from the response and set it as the authExpiryDate

    // Decode json response
    const data: LoginResponse = await response.json();

    userStore.$patch({ auth: true, authExpiryDate: data.expires });
    router.replace('/');
  } else {
    if (response.status === 401) {
      loginError.value = 'Invalid username or password';
    } else loginError.value = response.statusText;
  }
};

onMounted(() => {
  // Check if query parameter "logout" exists and is equal to "true"
  const query = router.currentRoute.value.query;
  if (query.logout === 'true') {
    loginError.value = 'You have been logged out';
  }
});
</script>

<template>
  <div class="page grid place-items-center">
    <div class="self-end text-center">
      <h1 class="text-3xl">Kobra Unleashed</h1>
      <p class="text-gray-500">Login to your dashboard</p>
    </div>

    <div>
      <form @submit.prevent="doLogin">
        <label for="username">Username</label>
        <input
          type="text"
          id="username"
          placeholder="Username"
          v-model="username"
          autocomplete="username"
        />
        <label for="password">Password</label>
        <input
          type="password"
          id="password"
          placeholder="Password"
          v-model="password"
          autocomplete="current-password"
        />
        <div class="flex items-center gap-x-1">
          <input type="checkbox" id="remember" v-model="remember" />
          <label for="remember">Remember me</label>
        </div>
        <button
          type="submit"
          class="flex items-center gap-x-2 btn btn-primary icon"
        >
          <LoginIcon class="w-6 h-6" />
          Login
        </button>
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
}
</style>
