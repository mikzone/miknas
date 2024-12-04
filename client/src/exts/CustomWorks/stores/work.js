import { defineStore } from 'pinia';
import { useExtension } from '../extMain';
import { MikCall, gutil } from 'miknas/utils';

const extsObj = useExtension();

// 缓存文件数据
export const useWorkStore = defineStore('CustomWorks', {
    state: function () {
        return {
            plugins: {},
            spaces: {},
            pluginSimples: {},
            isLoadMainInfo: false,
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
        async refreshMainInfo(force) {
            if (!force && this.isLoadMainInfo) return;
            let iRet = await extsObj.mcpost('maininfo');
            if (!iRet.suc) {
                MikCall.alertRespErrMsg(iRet);
                return;
            }
            let result = iRet.ret;
            this.spaces = result.spaces;
            this.pluginSimples = result.plugins;
            this.isLoadMainInfo = true;
        },
        async reloadPluginDef(pluginId) {
            if (!pluginId) return;
            let isOk = await MikCall.coMakeConfirm('确定重载插件吗?\n(应该只在插件定义发生改变的时候才需要重载)');
            if (!isOk) return;
            let iRet = await extsObj.mcpost('reloadPluginDef', { Id: pluginId });
            if (!iRet.suc) {
                MikCall.alertRespErrMsg(iRet);
                return;
            }
            this.queryPluginDef(pluginId, true);
            MikCall.sendSuccTips('重载插件成功');
        },
    },
});
