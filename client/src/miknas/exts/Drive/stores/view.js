import { defineStore } from 'pinia';

export const useViewStore = defineStore('driveView', {
  state: function () {
    return {
      fsid: '',
      fileList: [],  // 缓存的文件列表
    };
  },

  getters: {
  },

  actions: {
    cacheFileList(fsid, fileList) {
      this.fsid = fsid;
      this.fileList = fileList;
    },
  },
});
