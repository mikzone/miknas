import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'Drive';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: '文件管理',
  desc: '提供基础的文件管理功能',
  icon: 'folder',
  index: true,
  route: (extsObj) => {
    return {
      // component: () => import('./layouts/CustomLayout.vue'),
      component: (() => import('../Official/shares').then((module)=>module['SimpleNestRouterView'])),
      children: [
        {
          path: 'view/:fsid/:fspath(.*)?',
          name: extsObj.routeName('view'),
          component: () => import('./pages/ViewPage.vue'),
          props: true,
        },
        {
          path: 's/:shareid',
          name: extsObj.routeName('viewShare'),
          meta: {needLogined: false},
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid ,fsrela: '', kind: 'check' }),
        },
        {
          path: 's/:shareid/view/:routeSubPath(.*)?',
          name: extsObj.routeName('sview'),
          meta: {needLogined: false},
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid ,fsrela: route.params.routeSubPath, kind: 'view' }),
        },
        {
          path: 's/:shareid/list/:routeSubPath(.*)?',
          name: extsObj.routeName('slist'),
          meta: {
            needLogined: false,
            fsViewRouteName: extsObj.routeName('sview'),
          },
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid ,fsrela: route.params.routeSubPath, kind: 'list' }),
        },
      ],
    }
  },
});
