import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'Official';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: 'Official',
  desc: '提供MikNas相关基础功能',
  icon: 'apartment',
  route: (extsObj) => {
    // const extsObj = useExtension();
    return {
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
