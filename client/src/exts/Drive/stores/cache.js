import { defineStore } from 'pinia';
import { useExtension } from '../extMain';
import { MikCall, gutil } from 'miknas/utils';
import { FileUtil } from 'miknas/exts/Drive/shares';

const extsObj = useExtension();

function genSortFunc(key) {
  return (desc) => {
    return (a, b) => {
      let order = (desc && -1) || 1;
      let fwa = (a.isFile && -1) || 1;
      let fwb = (b.isFile && -1) || 1;
      if (fwa != fwb) return fwb - fwa;
      return ((b[key] > a[key] && -1) || 1) * order;
    };
  };
}

// 缓存文件数据
export const useCacheStore = defineStore('driveCache', {
  state: function () {
    return {
      fsCache: {}, // {fsid: {fspath: [fileInfo, ...]}}
      fsRec: {}, // {fsid: {一些操作习惯的缓存}}
      sortConfs: {
        name: {
          label: '名称',
          sortFunc: genSortFunc('name'),
        },
        modify: {
          label: '修改时间',
          sortFunc: genSortFunc('modify'),
        },
        size: {
          label: '文件大小',
          sortFunc: genSortFunc('size'),
        },
      },
    };
  },

  getters: {
    getOneRec: (state) => {
      return (fsid, fspath, key, defv) => {
        return gutil.getDictValueByKeys(
          state.fsRec,
          [fsid, fspath, key],
          defv
        );
      };
    },
    getAllRec: (state) => {
      return (fsid, fspath) => {
        return gutil.getDictValueByKeys(
          state.fsRec,
          [fsid, fspath],
          {}
        );
      };
    },
  },

  actions: {
    async tryRefreshFiles(fsid, fspath, withCache) {
      if (!fsid) return;
      let fsidCache = this.fsCache[fsid];
      if (!fsidCache) {
        fsidCache = {};
        this.fsCache[fsid] = fsidCache;
      }
      if (withCache && fsidCache[fspath]) return fsidCache[fspath];
      // await MikCall.coDelay(20000)
      let needFolderSize = this.getNeedFolderSize(fsid, fspath);
      if (needFolderSize === undefined) {
        let parentFspath = FileUtil.dir(fspath);
        let parentNeedFolderSize = this.getNeedFolderSize(fsid, parentFspath);
        if (parentNeedFolderSize == true) {
          this.updateNeedFolderSize(fsid, fspath, true);
          needFolderSize = true;
        }
      }
      let beginTs = Date.now();
      let iRet = await extsObj.mcpost('listFiles', { fsid, fspath, needFolderSize });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let endTs = Date.now();
      if (endTs - beginTs > 2000) {
        // 耗时太久的话，下次就不计算文件夹大小了
        this.updateNeedFolderSize(fsid, fspath, false);
      }
      let ret = iRet.ret;
      fspath = ret.fspath;
      for (let fileInfo of ret.files) {
        if (fileInfo.size) fileInfo.viewSize = FileUtil.formatSize(fileInfo.size);
        else fileInfo.viewSize = '';
        fileInfo.viewModify = FileUtil.formatTs(fileInfo.modify * 1000);
        if (fileInfo.isFile) {
          fileInfo.fileType = FileUtil.getFileType(fileInfo.name);
        } else {
          fileInfo.fileType = 'folder';
        }
        fileInfo.icon = FileUtil.typeDescs[fileInfo.fileType].icon;
      }
      if (!this.fsCache[fsid]) this.fsCache[fsid] = {};
      this.fsCache[fsid][fspath] = ret.files;
      this.updateRec(fsid, fspath, 'hasFolderSize', needFolderSize);
      this.sortFiles(fsid, fspath);
      return this.fsCache[fsid][fspath];
    },
    getSortSave(fsid, fspath) {
      let sortSave = this.getOneRec(fsid, fspath, 'sort');
      if (sortSave) {
        return sortSave;
      }
      this.updateRec(fsid, fspath, 'sort', {
        method: 'modify',
        desc: true,
      });
      return this.getOneRec(fsid, fspath, 'sort');
    },
    sortFiles(fsid, fspath, sortMethod, desc) {
      if (!sortMethod || desc === undefined) {
        let sortSave = this.getSortSave(fsid, fspath);
        if (!sortMethod) sortMethod = sortSave.method;
        if (desc === undefined) desc = sortSave.desc;
      }
      let fsidInfo = this.fsCache[fsid];
      if (!fsidInfo) return;
      let files = fsidInfo[fspath];
      if (!files) return;
      let methodConf = this.sortConfs[sortMethod];
      if (!methodConf) return;
      this.updateRec(fsid, fspath, 'sort', {
        method: sortMethod,
        desc,
      });

      files.sort(methodConf.sortFunc(desc));
    },
    updateRec(fsid, fspath, key, value) {
      gutil.setDictValueByKeys(this.fsRec, [fsid, fspath, key], value);
    },
    updateFsPosition(fsid, fspath, top) {
      this.updateRec(fsid, fspath, 'top', top);
    },
    getFsPosition(fsid, fspath) {
      return this.getOneRec(fsid, fspath, 'top');
    },
    updateNeedFolderSize(fsid, fspath, needFolderSize) {
      this.updateRec(fsid, fspath, 'needFolderSize', needFolderSize);
    },
    getNeedFolderSize(fsid, fspath) {
      return this.getOneRec(fsid, fspath, 'needFolderSize');
    },
  },
});
