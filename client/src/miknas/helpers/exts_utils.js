// Extensions相关工具

import { useOfficialStore } from 'miknas/exts/Official/stores/official';
import { MikCall } from './official_utils';

class Extension {

  constructor(extsConf) {
    this.extsConf = extsConf;
    this.id = extsConf.id;
    this.desc = extsConf.desc;
    this.title = extsConf.title;
    this.icon = extsConf.icon;
    this.route = extsConf.route;
    if (extsConf.index  === undefined) extsConf.index = true;
    this.index = extsConf.index;
    this.hasIndex = false;
    this.routeNamePrefix = `miknas_exts_${this.id}`;
    this.alias = undefined;
  }

  hasAuth(resid){
    return useOfficialStore().canAccess(this.id, resid);
  }

  serverUrl(extsSubUrl, param) {
    let url = useOfficialStore().mdServerUrl(`${this.id}/${extsSubUrl}`);
    return MikCall.genUrlWithParam(url, param);
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
    return useOfficialStore().extsClientUrl(this.id, extsSubUrl);
  }

  routeName(subName) {
    if (!subName) return this.routeNamePrefix;
    return `${this.routeNamePrefix}#${subName}`;
  }

  getIndex() {
    if (!this.index) return undefined;
    return this.routePath();
  }
}

var G_EXTS_INSTS = {}

export function getExtension(extsId) {
  return G_EXTS_INSTS[extsId];
}

export function getAllExtensions() {
  return G_EXTS_INSTS;
}

function getExtensionRoute(extsObj) {
  if (extsObj && extsObj.route) {
    let extsRoute = extsObj.route;
    if (typeof extsRoute == 'function') extsRoute = extsRoute(extsObj);
    return extsRoute;
  }
}

function verifyRouteValid(extsId, oneRoute) {
  let path = oneRoute.path;
  if (path) {
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
  let extsRoute = getExtensionRoute(extsObj);
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

export function scanAllExtension(ctx, files) {
  // 扫描所有的扩展
  let { router } = ctx;
  for (let module of Object.values(files)) {
    let useExtension = module.useExtension;
    if (useExtension) {
      let extsObj = useExtension();
      // 注册路由
      if (extsObj && extsObj.route) {
        let extsId = extsObj.id;
        let extsRoute = extsObj.route;
        if (typeof extsRoute == 'function') extsRoute = extsRoute(extsObj);
        extsRoute.path = extsObj.routePath('');
        extsRoute.name = extsObj.routeName('');
        extsRoute.component = extsRoute.component || (() => import('../exts/Official/shares').then((module)=>module['ExtensionPage']));
        extsRoute.meta = extsRoute.meta || {};
        extsRoute.meta.extsId = extsId;
        router.addRoute('miknas_exts', extsRoute);
      }
    }
  }
}
