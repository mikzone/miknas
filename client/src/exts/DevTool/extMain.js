
import { defineExtension } from 'miknas/utils';

// 此处定义扩展id，要和文件目录名保持一致
const EXTS_ID = 'DevTool';

// WARNING: 注意: 此处只供exts_util.js扫描注册使用，你在代码中不应该使用它
export const useExtension = defineExtension({
  id: EXTS_ID,
  title: '开发管理工具',
  desc: '用于开发时期各种调试测试',
  icon: 'construction',
  // 定义扩展是否有首页，boolean值，缺省则为true
  index: true,

  // 定义扩展路由, 格式应该是 undefined 或者是 VueRouter对应的RouteLocationRaw格式
  // 因为会被当作是被嵌套的路由，因此对于扩展顶层的route有一定限制
  // 首层的 path 和 name 都由 miknas 直接分配指定
  // children 下的请使用相对链接，如果要指定name，请使用 extsObj.routeName(扩展下自定义名称)

  // 多页面要用嵌套路由，示例如下
  route: (extsObj) => {
    return {
      component: () => import('./layouts/CustomLayout.vue'),
      children: [
        {
          path: '',
          name: extsObj.routeName('Index'),
          redirect: extsObj.routePath('ws'),
        },
        {
          path: 'db',
          component: () => import('./pages/ViewDbPage.vue'),
        },
        {
          path: 'view/:routeSubPath(.*)?',
          name: extsObj.routeName('view'),
          component: () => import('../Drive/shares').then((module)=>module['CommonExplorerPage']),
          // component: () => import('./pages/IndexPage.vue'),
          props: route => ({ fsid: 'Ws', fsrela: route.params.routeSubPath, kind: 'view' }),
        },
        {
          path: 'ws/:routeSubPath(.*)?',
          name: extsObj.routeName('ws'),
          meta: {
            fsViewRouteName: extsObj.routeName('view'),
          },
          component: () => import('../Drive/shares').then((module)=>module['CommonExplorerPage']),
          props: route => ({ fsid: 'Ws', fsrela: route.params.routeSubPath, kind: 'list' }),
        },
      ],
    }
  },

});

export default useExtension;
