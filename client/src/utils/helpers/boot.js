// import something here
import '../../css/app.scss'
import { gutil, MikCall } from './official_utils'
import { registerExtensions } from './exts_utils'
import { useOfficialStore } from 'miknas/exts/Official/stores/official';

// "async" is optional;
// more info on params: https://quasar.dev/quasar-cli/boot-files
export const boot = async (ctx) => {
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
    component: () => import('miknas/exts/Official/shares').then((module) => module['SimpleNestRouterView']),
    meta: { needLogined: true },
    name: 'miknas_exts',
    children: [],
  });

  // 注册所有的extension
  registerExtensions(ctx);

  router.addRoute('miknas_exts', {
    path: ':catchAll(.*)*',
    component: () => import('miknas/exts/Official/pages/ExtsNotFound.vue')
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
