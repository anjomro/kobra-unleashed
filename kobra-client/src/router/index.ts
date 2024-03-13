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
  // {
  //   path: '/settings',
  //   name: 'Settings',
  //   meta: { requiresAuth: true },
  //   component: () => import('@/views/Settings.vue'),
  // },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
  },
  // Add any other routes you have here
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

// Auth guard
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore();

  // If not logged in and not on the login page, redirect to login
  if (!userStore.auth && to.meta.requiresAuth === true) {
    return next('/login');
  }
  // If logged in and on the login page, redirect to dashboard
  else if (userStore.auth === true && to.name === 'Login') {
    return next('/');
  }

  next();
});

export default router;
