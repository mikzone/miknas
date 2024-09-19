import { defineStore } from 'pinia';
import { MikCall } from 'miknas/utils';
import useExtension from '../extMain';
import SelectFolderDlg from '../components/dialogs/SelectFolderDlg.vue';

const extsObj = useExtension();

export const RootFolderId = 0; // 根文件夹id

export const useNoteStore = defineStore('note', {
  state: function () {
    return {
      isDrawerOpen: false,
      folderList: [], // 服务端传下来的folderlist
      noteList: [],
    };
  },

  getters: {
    noteDict(state) {
      let ret = {};
      for (let noteInfo of state.noteList) {
        ret[noteInfo.id] = noteInfo;
      }
      return ret;
    },
    folderDict(state) {
      let rootInfo = {
        id: RootFolderId,
        name: '我的文件夹',
        type: 'folder',
        gid: 'folder_0',
        children: {},
      };
      let folderDict = {};
      folderDict[RootFolderId] = rootInfo;
      let allIds = [];
      for (let folderInfo of state.folderList) {
        folderDict[folderInfo.id] = folderInfo;
        allIds.push(folderInfo.id);
        folderInfo.type = 'folder';
        folderInfo.gid = `folder_${folderInfo.id}`;
        folderInfo.children = {};
      }
      for (let folderId of allIds) {
        let folderInfo = folderDict[folderId];
        let parentId = folderInfo.parent;
        let parentInfo = folderDict[parentId];
        // 找不到父节点的就用root
        if (!parentInfo) parentInfo = rootInfo;
        parentInfo.children[folderId] = folderInfo;
      }
      for (let noteInfo of state.noteList) {
        let folderInfo = folderDict[noteInfo.folder];
        if (folderInfo) {
          noteInfo.type = 'markdown';
          noteInfo.gid = `note_${noteInfo.id}`;
          folderInfo.children[noteInfo.gid] = noteInfo;
        }
      }
      return folderDict;
    },
    rootFolder() {
      return this.folderDict[RootFolderId];
    },
  },

  actions: {
    setIsDrawerOpen(flag) {
      this.isDrawerOpen = flag;
    },
    toggleDrawerOpen() {
      this.isDrawerOpen = !this.isDrawerOpen;
    },
    async refreshMyFolders(loadingMgr) {
      let stateName = `正在加载`;
      if (loadingMgr) loadingMgr.addLoadingState(stateName);
      let iRet = await extsObj.mcpost('getUserFolder');
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        if (loadingMgr) loadingMgr.removeLoadingState(stateName);
        return;
      }
      let result = iRet.ret;
      this.folderList = result;
      if (loadingMgr) loadingMgr.removeLoadingState(stateName);
    },
    addFolder(parentId, loadingMgr) {
      if (!parentId && parentId != 0) {
        MikCall.sendErrorTips('非法父文件夹');
        return;
      }
      MikCall.makePrompt('新文件夹名称', '', async (folderName) => {
        if (!folderName || folderName.length <= 0) {
          MikCall.sendErrorTips('名称不能为空');
          return;
        }
        let stateName = `正在添加`;
        if (loadingMgr) loadingMgr.addLoadingState(stateName);

        let iRet = await extsObj.mcpost('addFolder', {
          name: folderName,
          parent: parentId,
        });
        if (!iRet.suc) {
          MikCall.alertRespErrMsg(iRet);
          if (loadingMgr) loadingMgr.removeLoadingState(stateName);
          return;
        }
        if (loadingMgr) loadingMgr.removeLoadingState(stateName);
        await this.refreshMyFolders(loadingMgr);
      });
    },
    async selectFolder(title) {
      return await MikCall.coCreateDialog({
        component: SelectFolderDlg,
        componentProps: {
          title: title,
        },
      });
    },
    async modifyFolder(modifyInfo) {
      if (!modifyInfo.id) {
        MikCall.sendErrorTips('非法文件夹id');
        return;
      }
      let iRet = await extsObj.mcpost('modifyFolder', modifyInfo);
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      this.refreshMyFolders();
      return iRet.ret;
    },
    async renameFolder(folderId) {
      if (!folderId) {
        MikCall.sendErrorTips('非法文件夹');
        return;
      }
      let folderInfo = this.folderDict[folderId];
      if (!folderInfo) {
        MikCall.sendErrorTips('不存在的文件夹');
        return;
      }
      let [isOk, folderName] = await MikCall.coMakePrompt(
        `文件夹新名称`,
        folderInfo.name
      );
      if (!isOk) return;
      if (!folderName || folderName.length <= 0) {
        MikCall.sendErrorTips('名称不能为空');
        return;
      }
      if (folderName == folderInfo.name) {
        return;
      }
      await this.modifyFolder({
        id: folderId,
        name: folderName,
      });
    },
    async deleteFolder(folderId, loadingMgr) {
      if (!folderId) {
        MikCall.sendErrorTips('非法文件夹');
        return;
      }
      let folderInfo = this.folderDict[folderId];
      if (!folderInfo) {
        MikCall.sendErrorTips('不存在的文件夹');
        return;
      }
      let isOk = await MikCall.coMakeConfirm(
        `是否确认删除文件夹(${folderInfo.name})?`
      );
      if (!isOk) return;
      let stateName = `正在删除`;
      if (loadingMgr) loadingMgr.addLoadingState(stateName);
      let iRet = await extsObj.mcpost('deleteFolder', {
        id: folderId,
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        if (loadingMgr) loadingMgr.removeLoadingState(stateName);
        return;
      }
      if (loadingMgr) loadingMgr.removeLoadingState(stateName);
      await this.refreshMyFolders(loadingMgr);
    },
    async moveFolder(folderId) {
      if (!folderId) {
        MikCall.sendErrorTips('非法文件夹');
        return;
      }
      let folderInfo = this.folderDict[folderId];
      if (!folderInfo) {
        MikCall.sendErrorTips('不存在的文件夹');
        return;
      }
      let [isOk, newFolderId] = await this.selectFolder('移动到');
      if (!isOk) return;
      if (!newFolderId) return;
      return await this.modifyFolder({
        id: folderId,
        parent: newFolderId,
      });
    },
    async refreshMyNotes(loadingMgr) {
      let stateName = `正在加载笔记`;
      if (loadingMgr) loadingMgr.addLoadingState(stateName);

      let iRet = await extsObj.mcpost('getUserItemBrief');
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        if (loadingMgr) loadingMgr.removeLoadingState(stateName);
        return;
      }
      let result = iRet.ret;
      this.noteList = result;
      if (loadingMgr) loadingMgr.removeLoadingState(stateName);
    },
    updateStoreNote(noteInfo) {
      if (!noteInfo.id) return;
      let prevNoteInfo = this.noteDict[noteInfo.id];
      if (!prevNoteInfo) this.noteList.push(noteInfo);
      prevNoteInfo.folder = noteInfo.folder;
      prevNoteInfo.title = noteInfo.title;
    },
    deleteStoreNote(noteid) {
      if (!noteid) return;
      let prevNoteInfo = this.noteDict[noteid];
      if (!prevNoteInfo) return;
      this.noteList = this.noteList.filter((obj) => obj.id !== noteid);
    },
    async addNote(folderId, loadingMgr) {
      if (!folderId) {
        MikCall.sendErrorTips('不能在此处新建笔记');
        return;
      }
      let [isOk, title] = await MikCall.coMakePrompt(`笔记标题`, '');
      if (!isOk) return;
      if (!title || title.length <= 0) {
        MikCall.sendErrorTips('名称不能为空');
        return;
      }
      let stateName = `正在添加`;
      if (loadingMgr) loadingMgr.addLoadingState(stateName);

      let iRet = await extsObj.mcpost('addItem', {
        title: title,
        folder: folderId,
        content: '',
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        if (loadingMgr) loadingMgr.removeLoadingState(stateName);
        return;
      }
      if (loadingMgr) loadingMgr.removeLoadingState(stateName);
      await this.refreshMyNotes(loadingMgr);
    },
    async saveNote(modifyInfo) {
      if (!modifyInfo.id) {
        MikCall.sendErrorTips('修改笔记必须携带id');
        return;
      }
      let iRet = await extsObj.mcpost('modifyItem', modifyInfo);
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.preDealAttachesInNote(result);
      this.updateStoreNote(result);
      return result;
    },
    async getNote(noteid) {
      let iRet = await extsObj.mcpost('getItem', { id: noteid });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.preDealAttachesInNote(result);
      this.updateStoreNote(result);
      return result;
    },
    async deleteNote(noteInfo) {
      let noteid = noteInfo.id;
      if (!noteid) return;
      let confirmRet = await MikCall.coMakeConfirm(
        `是否确认删除笔记(${noteInfo.title})?`
      );
      if (!confirmRet) return;
      let iRet = await extsObj.mcpost('deleteItem', { id: noteid });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.deleteStoreNote(noteid);
      return result;
    },
    async listNotes({ folder, search, pageNum, pageSize }) {
      let iRet = await extsObj.mcpost('listItems', {
        folder,
        search,
        pageNum,
        pageSize,
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      for (let noteInfo of result.notes) {
        this.preDealAttachesInNote(noteInfo);
      }
      return result;
    },
    preDealAttach(attachInfo) {
      attachInfo.jsonData = JSON.parse(attachInfo.data);
    },
    preDealAttachesInNote(noteInfo) {
      for (let attachInfo of noteInfo.noteAttachs) {
        this.preDealAttach(attachInfo);
      }
    },
    async addAttach(callArgs) {
      let iRet = await extsObj.mcpost('addAttach', callArgs);
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.preDealAttach(result);
      return result;
    },
    async modifyAttach(callArgs) {
      let iRet = await extsObj.mcpost('modifyAttach', callArgs);
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      this.preDealAttach(result);
      return result;
    },
    async deleteAttach(attachId) {
      let iRet = await extsObj.mcpost('deleteAttach', { id: attachId });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      let result = iRet.ret;
      return result;
    },
  },
});
