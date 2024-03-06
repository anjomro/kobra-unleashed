import { useUserStore } from '@/stores/store';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';

// Import your views/components here

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Dashboard',
    meta: { requiresAuth: true },
    component: () => import('@/views/Dashboard.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    meta: { requiresAuth: false },
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/logout',
    name: 'Logout',
    meta: { requiresAuth: true },
    beforeEnter: () => {
      // Send a logout request to the server here
      fetch('/api/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      }).then((data) => {
        if (data.ok) {
          const userStore = useUserStore();
          userStore.$reset();
          router.push({ name: 'Login' });
        } else {
          router.push({ name: 'Login', query: { error: 'logout' } });
        }
      });
    },
    redirect: '/login',
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    meta: { requiresAuth: false },
    component: () => import('@/views/NotFound.vue'),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

// Auth guard
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore();
  // If not logged in and not on the login page, redirect to login
  if (!userStore.auth && to.meta.requiresAuth) {
    next('/login');
  }
  // If logged in and on the login page, redirect to home
  else if (userStore.auth && !to.meta.requiresAuth) {
    next('/');
  } else {
    next();
  }
});
export default router;
