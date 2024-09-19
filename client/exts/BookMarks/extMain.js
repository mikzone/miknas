
import { defineExtension } from 'miknas/utils';

// 此处定义扩展id，要和文件目录名保持一致
const EXTS_ID = 'BookMarks';

// WARNING: 注意: 此处只供exts_util.js扫描注册使用，你在代码中不应该使用它
export const useExtension = defineExtension({
  id: EXTS_ID,
  title: '书签',
  desc: '常用链接收藏夹',
  icon: 'bookmark',
  // 定义扩展是否有首页，boolean值，缺省则为true
  index: true,

  // 定义扩展路由, 格式应该是 undefined 或者是 VueRouter对应的RouteLocationRaw格式
  // 因为会被当作是被嵌套的路由，因此对于扩展顶层的route有一定限制
  // 首层的 path 和 name 都由 miknas 直接分配指定
  // children 下的请使用相对链接，如果要指定name，请使用 extsObj.routeName(扩展下自定义名称)

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

export default useExtension;
