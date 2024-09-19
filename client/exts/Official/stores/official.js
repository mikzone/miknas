import { useExtension as useOfficialExtension } from 'miknas/exts/Official/extMain';
import { defineStore } from 'pinia';
import { getExtension } from 'miknas/utils';
import { gutil, MikCall } from 'miknas/utils';

export const useOfficialStore = defineStore('official', {
  state: function () {
    return {
      // 用户名
      uid: '',
      // userAuths
      userAuths: {},
      // serverConfigs
      serverConfigs: {},
      extids: {},
    };
  },

  getters: {
    logoutUrl(state) {
      let extsId = state.serverConfigs['MIKNAS_AUTH_EXTS'];
      if (!extsId) return undefined;
      let extsObj = getExtension(extsId);
      if (!extsObj) return undefined;
      return extsObj.serverUrl('logout');
    },

    loginUrl(state) {
      let extsId = state.serverConfigs['MIKNAS_AUTH_EXTS'];
      if (!extsId) return undefined;
      let extsObj = getExtension(extsId);
      if (!extsObj) return undefined;
      return extsObj.serverUrl('login');
    },

  },

  actions: {
    modifyStateDict(modifyDict) {
      for (const [k, v] of Object.entries(modifyDict)) {
        if (k in this) {
          this[k] = v;
        }
      }
    },

    // 客户端初始化好的时候
    async loadOnInit() {
      // 先去加载玩家数据
      let extsObj = useOfficialExtension()
      let iRet = await extsObj.mcpost('getClientInitInfo');
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.modifyStateDict(result);
    },

    mdClientUrl(subUrl) {
      subUrl = subUrl || '';
      let prefix = this.serverConfigs['MIKNAS_CLIENT_PREFIX'] || '/c';
      prefix = gutil.strip(prefix, '/', false, true);
      subUrl = gutil.strip(subUrl, '/', true, false);
      return `${prefix}/${subUrl}`;
    },

    extsClientUrl(extsId, extSubUrl) {
      extSubUrl = extSubUrl || '';
      let urlMap = this.serverConfigs['MIKNAS_CLIENT_URL_MAP'] || {};
      let name = urlMap[extsId];
      if (name === undefined) name = extsId;
      extSubUrl = gutil.strip(extSubUrl, '/', true, false);
      return this.mdClientUrl(`${name}/${extSubUrl}`);
    },

    mdServerUrl(subUrl) {
      let prefix = process.env.MIKNAS_SERVER_PREFIX || '/s';
      prefix = gutil.strip(prefix, '/', false, true);
      subUrl = gutil.strip(subUrl, '/', true, false);
      return `${prefix}/${subUrl}`;
    },

    canAccess(extsId, resid){
      let res = `${extsId}/${resid}`;
      return gutil.authCheck(res, this.userAuths);
    },
  },
});
