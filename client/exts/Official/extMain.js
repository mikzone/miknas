import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'Official';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: 'MikNas',
  desc: '提供MikNas相关基础功能',
  icon: 'apartment',
  route: (extsObj) => {
    // const extsObj = useExtension();
    return {
      // component: () => import('./layouts/CustomLayout.vue'),
      // component: () => import('./shares').then((module)=>module['ExtensionPage']),
      children: [
        {
          path: '',
          name: extsObj.routeName('Index'),
          component: () => import('./pages/AllExtsPage.vue'),
        },
      ],
    }
  },
});
