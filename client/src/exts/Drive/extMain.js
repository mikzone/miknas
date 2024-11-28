import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'Drive';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: '文件管理',
  desc: '提供基础的文件管理功能',
  icon: 'folder',
  headerComponent: false,
  route: (extsObj) => {
    return {
      children: [
        {
          path: 'view/:fsid/:fspath(.*)?',
          name: extsObj.routeName('view'),
          component: () => import('./pages/ViewPage.vue'),
          meta: {
            fullCtrlLayout: true,
          },
          props: true,
        },
        {
          path: 's/:shareid',
          name: extsObj.routeName('viewShare'),
          meta: {
            fullCtrlLayout: true,
          },
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid, fsrela: '', kind: 'check' }),
        },
        {
          path: 's/:shareid/view/:routeSubPath(.*)?',
          name: extsObj.routeName('sview'),
          meta: {
            fullCtrlLayout: true,
          },
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid, fsrela: route.params.routeSubPath, kind: 'view' }),
        },
        {
          path: 's/:shareid/list/:routeSubPath(.*)?',
          name: extsObj.routeName('slist'),
          meta: {
            fsViewRouteName: extsObj.routeName('sview'),
            fullCtrlLayout: true,
          },
          component: () => import('./pages/ViewShare.vue'),
          props: route => ({ shareid: route.params.shareid, fsrela: route.params.routeSubPath, kind: 'list' }),
        },
      ],
    }
  },
});
