// Extensions相关工具

import { useOfficialStore } from 'miknas/exts/Official/stores/official.js';
import { MikCall, gutil } from './official_utils';
import { defineAsyncComponent } from 'vue';
import { gPinia } from './instance';

class Extension {

  constructor(extsConf) {
    this.extsConf = extsConf;
    this.id = extsConf.id;
    this.desc = extsConf.desc;
    this.title = extsConf.title;
    this.icon = extsConf.icon;
    this.route = extsConf.route;
    if (extsConf.headerComponent === false) {
      this.headerComponent = undefined;
    }
    else {
      let com = extsConf.headerComponent || (() => import('miknas/exts/Official/shares').then((module) => module['LayoutHeader']));
      this.headerComponent = defineAsyncComponent(com);
    }
    this.routeNamePrefix = `miknas_exts_${this.id}`;
    this.alias = undefined;
  }

  res(extResId) {
    return `${this.id}/${extResId}`;
  }

  hasAuth(extResId) {
    let resid = this.res(extResId);
    return useOfficialStore(gPinia).canAccess(resid);
  }

  serverUrl(extsSubUrl, param, isfull) {
    let url = useOfficialStore(gPinia).mdServerUrl(`${this.id}/${extsSubUrl}`);
    let href = MikCall.genUrlWithParam(url, param);
    if (isfull) {
      return gutil.genFullUrl(href);
    }
    return href
  }

  mcpost(extsSubUrl, postData, extraConf) {
    let url = this.serverUrl(extsSubUrl);
    return MikCall.mcpost(url, postData, extraConf);
  }

  mcget(extsSubUrl, param, extraConf) {
    let url = this.serverUrl(extsSubUrl);
    return MikCall.mcget(url, param, extraConf);
  }

  routePath(extsSubUrl) {
    return useOfficialStore(gPinia).extsClientUrl(this.id, extsSubUrl);
  }

  routeName(subName) {
    if (!subName) return this.routeNamePrefix;
    return `${this.routeNamePrefix}#${subName}`;
  }

  getIndex() {
    let convRoute = this.convRoute;
    if (!convRoute) return undefined;
    let IndexName = this.routeName('Index');
    let oneRoute = dfsFindOneRoute(convRoute, IndexName);
    if (!oneRoute) return undefined;
    let meta = oneRoute.meta;
    if (!meta) return undefined;
    // 判断它是否能访问
    let [suc] = useOfficialStore(gPinia).canAccessByMeta(meta);
    if (!suc) {
      return undefined;
    }
    return { name: IndexName };
  }
}

var G_EXTS_INSTS = {}

export function getExtension(extsId) {
  return G_EXTS_INSTS[extsId];
}

export function getAllExtensions() {
  return G_EXTS_INSTS;
}

function dfsFindOneRoute(oneRoute, name) {
  if (oneRoute.name == name) return oneRoute;
  if (oneRoute.children) {
    for (let tmpRoute of oneRoute.children) {
      let result = dfsFindOneRoute(tmpRoute, name);
      if (result) return result;
    }
  }
}

function dfsHandleRoute(oneRoute, parentMeta) {
  // 递归的处理route
  // 主要是处理meta，聚合一下信息
  let curMeta = oneRoute.meta;
  if (!curMeta) {
    curMeta = {};
    oneRoute.meta = curMeta;
  }
  for (let [k, v] of Object.entries(parentMeta)) {
    if (curMeta[k] === undefined) {
      // 和vue router同样的update方式
      curMeta[k] = v;
    }
  }
  if (oneRoute.children) {
    for (let tmpRoute of oneRoute.children) {
      dfsHandleRoute(tmpRoute, curMeta);
    }
  }
  return true
}

function calcOrigRoute(extsObj) {
  let extsRoute = extsObj.route;
  if (typeof extsRoute == 'function') extsRoute = extsRoute(extsObj);
  return extsRoute;
}

function verifyRouteValid(extsId, oneRoute) {
  // 这里校验的是原始的route，不是convRoute，因为在校验时，还没搜到服务端下发的客户端前缀
  let path = oneRoute.path;
  if (path) { // 根节点是自动填充上去的，不用检查，规则不一样
    if (path.startsWith('/')) throw `扩展(${extsId})里含有非法路由路径(${path}), 不能以'/'开头)`;
    let firstLetter = path.charAt(0);
    if (firstLetter >= 'A' && firstLetter <= 'Z') throw `扩展(${extsId})里含有非法路由路径(${path}), 为了方便指定默认扩展,不能以大写字母开头)`;
  }
  if (oneRoute.children) {
    for (let tmpRoute of oneRoute.children) {
      let suc = verifyRouteValid(extsId, tmpRoute);
      if (!suc) return false;
    }
  }
  return true
}

function verifyExtension(extsObj) {
  // 校验一下扩展
  let firstLetter = extsObj.id.charAt(0);
  if (firstLetter < 'A' || firstLetter > 'Z') throw `扩展名必须以大写字母开头,当前扩展名不符合${extsObj.id}`;
  let extsRoute = calcOrigRoute(extsObj);
  if (extsRoute) {
    verifyRouteValid(extsObj.id, extsRoute);
  }
}

// 注意：defineExtension这个应该只在扩展的 extMain.js 中使用，其余时候由 getExtension 来获取
export function defineExtension(extsConf) {
  function GetExtensionInst() {
    let extsId = extsConf.id;
    if (!extsId) {
      throw 'defineExtension Error: Id require';
    }
    let extsObj = G_EXTS_INSTS[extsId];
    if (extsObj) {
      if (extsObj.extsConf === extsConf) return extsObj;
      else throw `defineExtension Error: dumplicat extension Id, ${extsId} is already exist`;
    }
    extsObj = new Extension(extsConf);
    verifyExtension(extsObj);
    G_EXTS_INSTS[extsId] = extsObj;
    return extsObj;
  }
  return GetExtensionInst;
}

export function registerExtensions(ctx) {
  // 扫描所有的扩展
  let { router, extsObjMap } = ctx;
  for (let extsObj of Object.values(extsObjMap)) {
    // 注册路由
    if (extsObj && extsObj.route) {
      let extsRoute = calcOrigRoute(extsObj);

      let extsId = extsObj.id;
      extsRoute.path = extsObj.routePath('');
      extsRoute.name = extsObj.routeName('');
      extsRoute.meta = extsRoute.meta || {};
      extsRoute.meta.extsId = extsId;
      dfsHandleRoute(extsRoute, {});

      extsObj.convRoute = extsRoute;
      // console.log('convRoute', extsRoute);
      router.addRoute('miknas_exts', extsRoute);
    }
  }
}

export function scanAllExtension(extsObjMap, modules) {
  // 扫描所有的扩展
  for (let module of modules) {
    let useExtension = module.useExtension;
    if (useExtension) {
      let extsObj = useExtension();
      let extsId = extsObj.id;
      extsObjMap[extsId] = extsObj;
    }
  }
}

export function calcValidExtsInfos(excludeIds) {
  let ret = {};
  let allExtsObjs = getAllExtensions();
  for (let [extsId, extsObj] of Object.entries(allExtsObjs)) {
    if (excludeIds && excludeIds.includes(extsId)) continue;
    let indexTo = extsObj.getIndex();
    if (!indexTo) continue;
    let info = {
      id: extsObj.id,
      desc: extsObj.desc,
      title: extsObj.title,
      icon: extsObj.icon || 'extension',
      indexTo: indexTo,
    };
    ret[extsObj.id] = info;
  }
  return ret;
}