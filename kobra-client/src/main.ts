import { createPinia } from 'pinia';
import { createApp } from 'vue';
import '@/scss/main.scss';
import App from '@/App.vue';
import router from './router';
import { cloneDeep } from 'lodash';

const app = createApp(App);

const store = createPinia();

store.use(({ store }) => {
  const initialState = cloneDeep(store.$state);
  store.$reset = () => {
    store.$patch(($state) => {
      Object.assign($state, initialState);
    });
  };
});

app.use(store);

app.use(router);

app.mount('#app');
