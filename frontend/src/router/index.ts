import { createRouter, createWebHistory } from 'vue-router';
import DashboardView from '@/views/DashboardView.vue';
import RecordsView from '@/views/RecordsView.vue';
import SearchView from '@/views/SearchView.vue';
import SettingsView from '@/views/SettingsView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView,
      meta: { title: 'Dashboard' }
    },
    {
      path: '/records',
      name: 'records',
      component: RecordsView,
      meta: { title: 'Records' }
    },
    {
      path: '/records/:id',
      name: 'record-detail',
      component: () => import('@/views/RecordDetailView.vue'),
      meta: { title: 'Record Details' }
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView,
      meta: { title: 'Search Images' }
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
      meta: { title: 'Settings' }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { title: 'Page Not Found' }
    }
  ]
});

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - Image RAG` : 'Image RAG';
  next();
});

export default router;