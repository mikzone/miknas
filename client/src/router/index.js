import { route } from 'quasar/wrappers';
import {
  createRouter,
  createMemoryHistory,
  createWebHistory,
  createWebHashHistory,
} from 'vue-router';
import routes from './routes';

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default route(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : process.env.VUE_ROUTER_MODE === 'history'
    ? createWebHistory
    : createWebHashHistory;

  const Router = createRouter({
    // scrollBehavior: () => ({ left: 0, top: 0 }),
    scrollBehavior(to, from, savedPosition) {
      if (savedPosition) {
        // console.log('xxxxxxxxxx', 'scrollBehavior', to, from, savedPosition);
        if (to.meta.delayScroll) {
          return new Promise((resolve) => {
            setTimeout(() => {
              resolve(savedPosition)
            }, to.meta.delayScroll)
          })
        } else {
          return savedPosition;
        }
      } else if (to.hash) {
        return { el: to.hash };
      } else {
        return { left: 0, top: 0 };
      }
    },
    routes,

    // Leave this as is and make changes in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    history: createHistory(process.env.VUE_ROUTER_BASE),
  });

  return Router;
});
