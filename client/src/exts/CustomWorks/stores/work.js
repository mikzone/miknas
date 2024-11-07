import { defineStore } from 'pinia';
import { useExtension } from '../extMain';
import { MikCall, gutil } from 'miknas/utils';

const extsObj = useExtension();

// 缓存文件数据
export const useWorkStore = defineStore('CustomWorks', {
    state: function () {
        return {
            plugins: {},
        };
    },

    getters: {
    },

    actions: {
        async queryPluginDef(pluginId, force) {
            if (!pluginId) return;
            if (!force && this.plugins[pluginId]) return this.plugins[pluginId];
            let iRet = await extsObj.mcpost('queryPluginDef', { Id: pluginId });
            if (!iRet.suc) {
                MikCall.alertRespErrMsg(iRet);
                return;
            }
            let result = iRet.ret;
            result._jobMap = gutil.list2map(result.Jobs, 'Id');
            this.plugins[pluginId] = result;
        },
    },
});
