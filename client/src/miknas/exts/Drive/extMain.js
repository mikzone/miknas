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
      children: [
        // {
        //   path: '',
        //   name: extsObj.routeName('Index'),
        //   redirect: extsObj.routePath('ws'),
        // },
        {
          path: 'view/:fsid/:fspath(.*)?',
          name: extsObj.routeName('view'),
          component: () => import('./pages/ViewPage.vue'),
          props: true,
        },
      ],
    }
  },
});
