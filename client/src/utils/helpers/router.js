import {
  createRouter,
  createMemoryHistory,
  createWebHistory,
} from 'vue-router';

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

const createHistory = import.meta.env.SERVER
  ? createMemoryHistory
  : createWebHistory;

export const router = createRouter({
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      if (to.meta.delayScroll) {
        return new Promise((resolve) => {
          setTimeout(() => {
            resolve(savedPosition);
          }, to.meta.delayScroll);
        });
      } else {
        return savedPosition;
      }
    } else if (to.hash) {
      return { el: to.hash };
    } else {
      return { left: 0, top: 0 };
    }
  },
  routes: [],

  // Leave this as is and make changes in quasar.conf.js instead!
  // quasar.conf.js -> build -> vueRouterMode
  // quasar.conf.js -> build -> publicPath
  history: createHistory(import.meta.env.VUE_ROUTER_BASE),
});
