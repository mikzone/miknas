// import something here

import { gutil, MikCall, scanAllExtension } from './utils'
import { useOfficialStore } from './exts/Official/stores/official';

// "async" is optional;
// more info on params: https://quasar.dev/quasar-cli/boot-files
export default async (ctx) => {
  // something to do
  let { app, router } = ctx;

  window.gutil = gutil;
  gutil.setCacheData('app', app);
  gutil.setCacheData('router', router);

  const officialStore = useOfficialStore();
  await officialStore.loadOnInit();

  // 注册路由
  router.addRoute({
    path: officialStore.mdClientUrl('/'),
    component: () => import('./exts/Official/shares').then((module)=>module['SimpleNestRouterView']),
    meta: {needLogined: true},
    name: 'miknas_exts',
    children: [],
  });

  // 注册所有的extension
  scanAllExtension(ctx, import.meta.globEager('./exts/*/extMain.js'));

  router.addRoute('miknas_exts', {
    path: ':catchAll(.*)*',
    component: () => import('./exts/Official/pages/ExtsNotFound.vue')
  });

  router.beforeEach((to) => {
    let needLogined = to.meta.needLogined;
    if (needLogined) {
      if (!officialStore.uid) {
        MikCall.sendErrorTips('您尚未登录');
        let loginUrl = officialStore.loginUrl;
        if (loginUrl) window.location.href = officialStore.loginUrl;
        return false;
      }
    }
  })

}
