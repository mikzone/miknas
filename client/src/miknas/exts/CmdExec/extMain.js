import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'CmdExec';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: '命令执行',
  desc: '实时执行脚本命令',
  icon: 'terminal',
  index: true,

  // 多页面要用嵌套路由，示例如下
  route: (extsObj) => {
    return {
      children: [
        {
          path: '',
          name: extsObj.routeName('Index'),
          component: () => import('./pages/IndexPage.vue'),
        },
      ],
    }
  },
});
