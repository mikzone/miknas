import { computed, nextTick, onActivated, onDeactivated } from 'vue';
import { date, Dialog, openURL } from 'quasar';
import { gutil, MikCall } from 'miknas/utils';
import { onMounted, reactive } from 'vue';
import { useExtension } from './extMain';
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute } from 'vue-router';
import { useViewStore } from './stores/view';
import { useCacheStore } from './stores/cache';
import {
  matIconArticle,
  matIconFolder,
  matIconImage,
  matIconMovie,
  matIconQuestionMark,
} from 'miknas/exts/Official/shares/mysvgs';
import MdcFolderSelectDlg from './components/MdcFolderSelectDlg.vue';
import MdcFileSelectDlg from './components/MdcFileSelectDlg.vue';
import MdcFileUploadDlg from './components/MdcFileUploadDlg.vue';
import MdcFileViewContainerDlg from './components/FileView/MdcFileViewContainerDlg.vue';
import { openMultOperateDlg } from '../Official/shares';

const viewStore = useViewStore();
const cacheStore = useCacheStore();

const PROCESS_COLUMNS = [
  {
    name: 'name',
    label: '文件名',
    field: 'name',
    sortable: false,
    align: 'left',
  },
  {
    name: 'modify',
    label: '修改时间',
    field: 'viewModify',
    sortable: false,
    align: 'left',
  },
  {
    name: 'size',
    label: '大小',
    field: 'size',
    sortable: false,
    align: 'right',
  },
];

const PathSep = '/';

export const FileUtil = {
  sizeUnits: ['B', 'K', 'M', 'G'],
  sizeStep: 1024,
  fileExts2Type: {
    apng: 'img',
    avif: 'img',
    gif: 'img',
    jpg: 'img',
    jpeg: 'img',
    jfif: 'img',
    pjpeg: 'img',
    pjp: 'img',
    png: 'img',
    svg: 'img',
    webp: 'img',
    bmp: 'img',
    ico: 'img',
    cur: 'img',
    tif: 'img',
    tiff: 'img',

    wmv: 'video',
    avi: 'video',
    mpeg: 'video',
    mpg: 'video',
    rm: 'video',
    rmvb: 'video',
    flv: 'video',
    mp4: 'video',
    '3gp': 'video',
    mov: 'video',
    divx: 'video',
    vob: 'video',
    mkv: 'video',
    fli: 'video',
    flc: 'video',
    f4v: 'video',
    m4v: 'video',
    mod: 'video',
    m2t: 'video',
    webm: 'video',
    mts: 'video',
    m2ts: 'video',
    '3g2': 'video',
    mpe: 'video',
    ts: 'video',
    div: 'video',
    lavf: 'video',
    dirac: 'video',

    txt: 'text',
    md: 'text',
    json: 'text',
    c: 'text',
    cpp: 'text',
    h: 'text',
    hpp: 'text',
    js: 'text',
    vue: 'text',
    html: 'text',
    htm: 'text',
    css: 'text',
    sql: 'text',
    log: 'text',
    ini: 'text',
    yml: 'text',
    yaml: 'text',
    py: 'text',
    go: 'text',
    java: 'text',
    php: 'text',
    xml: 'text',
  },
  typeDescs: {
    text: {
      name: '文本',
      icon: 'article',
      svg: matIconArticle,
    },
    video: {
      name: '视频',
      icon: 'movie',
      svg: matIconMovie,
    },
    img: {
      name: '图像',
      icon: 'image',
      svg: matIconImage,
    },
    unkown: {
      name: '未知',
      icon: 'question_mark',
      svg: matIconQuestionMark,
    },
    folder: {
      name: '文件夹',
      icon: 'folder',
      svg: matIconFolder,
    },
  },
  baseName(filePath) {
    let l = filePath.split(PathSep);
    return l[l.length - 1];
  },
  contactFolderName(name1, name2) {
    if (name1.length <= 0) return name2;
    return name1 + PathSep + name2;
  },
  getfileExtsion(filePath) {
    let name = FileUtil.baseName(filePath);
    let l = name.split('.');
    if (l.length <= 1) return '';
    return l[l.length - 1];
  },
  getFileType(filePath) {
    let ext = FileUtil.getfileExtsion(filePath);
    ext = ext.toLowerCase();
    return FileUtil.fileExts2Type[ext] || 'unkown';
  },
  calcRelativePath(curPath, rootPath) {
    if (!curPath.startsWith(rootPath)) return;
    let subPath = curPath.slice(rootPath.length);
    subPath = gutil.strip(subPath, '/', true, false);
    return subPath;
  },
  formatFloat(num) {
    return num.toLocaleString(undefined, {
      maximumFractionDigits: 2,
    });
  },
  formatSize(size) {
    for (let unit of FileUtil.sizeUnits) {
      if (size < FileUtil.sizeStep) {
        return `${FileUtil.formatFloat(size)}${unit}`;
      }
      size = (1.0 * size) / FileUtil.sizeStep;
    }
    size = size * FileUtil.sizeStep;
    return `${FileUtil.formatFloat(size)}${
      FileUtil.sizeUnits[FileUtil.sizeUnits.length - 1]
    }`;
  },
  formatTs(timeStamp) {
    return date.formatDate(timeStamp, 'YYYY-MM-DD HH:mm:ss');
  },
};

export function useFileView(fsid, rootPath, extraConf) {
  rootPath = rootPath || '';
  extraConf = extraConf || {};

  const fileState = reactive({
    fsid: fsid,
    filterTxt: '',
    loadingStates: {},
    processColumns: PROCESS_COLUMNS,
    curPath: '',
    curFiles: [],
    isSelectMode: false,
    curSelected: [],
    lastShowMode: 'list',
    showModeConfs: {
      list: {
        desc: '列表显示',
        icon: 'view_list',
      },
      grid: {
        desc: '网格显示',
        icon: 'grid_view',
      },
      blog: {
        desc: '单列原图显示',
        icon: 'width_full',
      },
    },
  });

  function contactFolderName(name1, name2) {
    if (name1.length <= 0) return name2;
    return name1 + PathSep + name2;
  }

  const fileGetter = {
    fileBC: computed(() => {
      let ret = [];
      // let rootPath = rootPath;
      let curPath = fileState.curPath;
      if (!curPath.startsWith(rootPath))
        return [{ name: `非法路径(${curPath})` }];
      ret.push({
        icon: 'home',
        name: fsid,
        path: rootPath,
      });
      let subPath = curPath.slice(rootPath.length);
      let subFodlers = subPath.split(PathSep);
      let genPath = rootPath;
      for (let folderName of subFodlers) {
        if (folderName) {
          genPath = contactFolderName(genPath, folderName);
          ret.push({
            name: folderName,
            path: genPath,
          });
        }
      }
      // 最后一个不能点击
      delete ret[ret.length - 1]['path'];
      return ret;
    }),
    loadingLabel: computed(() => {
      return Object.keys(fileState.loadingStates).join(',');
    }),
    isloading: computed(() => {
      return fileGetter.loadingLabel.value.length > 0;
    }),
    grid: computed(() => {
      return fileGetter.showMode.value != 'list';
    }),
    gridClass: computed(() => {
      switch (fileGetter.showMode.value) {
        case 'grid':
          return 'col-xs-4 col-sm-3 col-md-2 col-lg-1';
        case 'blog':
          return 'col-12';
        default:
          return '';
      }
    }),
    imgRatio: computed(() => {
      switch (fileGetter.showMode.value) {
        case 'grid':
          return 1;
        default:
          return undefined;
      }
    }),
    showMode: computed(() => {
      let mode = fileOp.getShowMode();
      if (mode && fileState.showModeConfs[mode]) {
        fileState.lastShowMode = mode;
        return mode;
      }
      mode = fileState.lastShowMode;
      fileOp.updateShowMode(mode);
      return fileOp.getShowMode();
    }),
    sortConfs: computed(() => {
      let ret = cacheStore.sortConfs;
      return ret;
    }),
    curSortSave: computed(() => {
      return cacheStore.getSortSave(fsid, fileState.curPath);
    }),
    selectAllStatusIcon: computed(()=>{
      let selectLen = fileState.curSelected.length;
      if (selectLen <= 0) return 'check_box_outline_blank';
      else if (selectLen == fileState.curFiles.length) return 'check_box';
      else return 'indeterminate_check_box';
    }),
  };

  const extsObj = useExtension();

  const fileOp = {
    makeConfirm: MikCall.makeConfirm,

    addLoadingState(stateName) {
      fileState.loadingStates[stateName] = true;
    },

    removeLoadingState(stateName) {
      if (fileState.loadingStates[stateName])
        delete fileState.loadingStates[stateName];
    },

    async gotoPath(path) {
      if (path === undefined || path === null) return;
      if (path == fileState.curPath) return;
      if (extraConf.openDirFn) {
        // 通过路由跳转
        let subPath = FileUtil.calcRelativePath(path, rootPath);
        await extraConf.openDirFn(subPath);
      } else {
        await fileOp.tryRefreshFiles(path);
      }
    },

    async tryRefreshFiles(path, withCache) {
      // 该方法直接加载path，如果是打开文件之类会影响route的，请用gotoPath
      if (path === undefined) {
        path = fileState.curPath;
        if (path === undefined) return;
      }
      let stateName = '加载目录中';
      fileOp.addLoadingState(stateName);
      let fileInfos = await cacheStore.tryRefreshFiles(fsid, path, withCache);
      fileState.curPath = path;
      fileState.curFiles = fileInfos;
      fileOp.quitSelectMode();
      if (extraConf.refreshCb) {
        extraConf.refreshCb();
      }
      if (withCache && extraConf.isRememberPos) {
        // 请求有缓存的，那就说明是通过路由变化来的，尝试恢复定位
        let prevTop = cacheStore.getFsPosition(fsid, path);
        if (prevTop) {
          nextTick(() => {
            window.scrollTo({ top: prevTop });
          });
        }
      }
      fileOp.removeLoadingState(stateName);
    },

    tryOpenSubfoder(folderName) {
      let newPath = contactFolderName(fileState.curPath, folderName);
      fileOp.gotoPath(newPath);
    },

    updateShowMode(newMode) {
      if (fileState.showModeConfs[newMode]) {
        cacheStore.updateRec(fsid, fileState.curPath, 'showMode', newMode);
      }
    },

    getShowMode() {
      return cacheStore.getOneRec(fsid, fileState.curPath, 'showMode');
    },

    updateSort(sortMethod) {
      let sortSave = fileGetter.curSortSave.value;
      if (sortSave.method == sortMethod) {
        cacheStore.sortFiles(
          fsid,
          fileState.curPath,
          sortMethod,
          !sortSave.desc
        );
        return;
      }
      cacheStore.sortFiles(fsid, fileState.curPath, sortMethod, false);
    },

    abs(relPath) {
      return contactFolderName(fileState.curPath, relPath);
    },

    relToRoot(fullpath) {
      return FileUtil.calcRelativePath(fullpath, rootPath);
    },

    openFile(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('文件名不能为空');
        return;
      }
      let initFilePath = contactFolderName(fileState.curPath, fileName);
      let fileList = [];
      for (let fileInfo of fileState.curFiles) {
        if (fileInfo.isFile && fileInfo.name) {
          let filePath = contactFolderName(fileState.curPath, fileInfo.name);
          fileList.push(filePath);
        }
      }
      viewStore.cacheFileList(fsid, fileList);
      if (extraConf.viewFn) {
        let subPath = FileUtil.calcRelativePath(initFilePath, rootPath);
        extraConf.viewFn(subPath);
        return;
      }
      viewFile(fileState.fsid, initFilePath, fileList);
    },

    clickOpen(fileInfo) {
      if (fileState.isSelectMode) {
        fileOp.toggleSelect(fileInfo);
        return;
      }
      if (fileInfo.isFile) fileOp.openFile(fileInfo.name);
      else fileOp.tryOpenSubfoder(fileInfo.name);
    },

    downloadFile(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('文件名不能为空');
        return;
      }
      let path = contactFolderName(fileState.curPath, fileName);
      openURL(fileOp.getDownloadUrl(path));
    },

    async removeFile(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('文件名不能为空');
        return;
      }
      let isOk = await MikCall.coMakeConfirm(`确认删除 ${fileName} ?`);
      if (!isOk) return;
      let iRet = await extsObj.mcpost('removeFile', {
        fsid: fileState.fsid,
        fspath: contactFolderName(fileState.curPath, fileName),
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      await fileOp.tryRefreshFiles();
    },

    async reqCopyFile(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('文件名不能为空');
        return;
      }
      let [isOk, toPath] = await openSelectFsidFolderDlg(
        fsid,
        rootPath,
        '复制到此处'
      );
      if (!isOk) return;

      fileOp.addLoadingState('复制中');
      let iRet = await extsObj.mcpost('copyFile', {
        fsid: fileState.fsid,
        fspath: contactFolderName(fileState.curPath, fileName),
        topath: toPath,
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        fileOp.removeLoadingState('复制中');
        return;
      }
      await fileOp.tryRefreshFiles();
      fileOp.removeLoadingState('复制中');
      // 强制刷新一下目标目录
      await cacheStore.tryRefreshFiles(fsid, toPath);
    },

    async reqMvFile(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('文件名不能为空');
        return;
      }
      let [isOk, toPath] = await openSelectFsidFolderDlg(
        fsid,
        rootPath,
        '移动到此处'
      );
      if (!isOk) return;
      fileOp.addLoadingState('移动中');
      let iRet = await extsObj.mcpost('mvFile', {
        fsid: fileState.fsid,
        fspath: contactFolderName(fileState.curPath, fileName),
        topath: toPath,
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        fileOp.removeLoadingState('移动中');
        return;
      }
      fileOp.removeLoadingState('移动中');
      await fileOp.tryRefreshFiles();
      // 强制刷新一下目标目录
      await cacheStore.tryRefreshFiles(fsid, toPath);
    },

    async rename(fileName) {
      if (!fileName) {
        MikCall.sendErrorTips('要先选择移动的文件');
        return;
      }
      let [isOk, newName] = await MikCall.coMakePrompt(`新文件名称`, fileName);
      if (!isOk) return;
      if (!newName) {
        MikCall.sendErrorTips('新文件名不能为空');
        return;
      }
      let iRet = await extsObj.mcpost('renameFile', {
        fsid: fileState.fsid,
        fspath: contactFolderName(fileState.curPath, fileName),
        toname: newName,
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      await fileOp.tryRefreshFiles();
    },

    openUploadDlg() {
      Dialog.create({
        component: MdcFileUploadDlg,
        componentProps: {
          factoryInfo: {
            url: extsObj.serverUrl('uploadFiles'),
            method: 'POST',
            formFields: [
              { name: 'fspath', value: fileState.curPath },
              { name: 'fsid', value: fileState.fsid },
            ],
          },
        },
      }).onOk(() => {
        fileOp.tryRefreshFiles();
      });
    },

    async newFolder() {
      let [isOk, folderName] = await MikCall.coMakePrompt(
        `请输入新的文件夹名称`,
        '',
        '新建文件夹'
      );
      if (!isOk) return;
      if (!folderName || folderName.length <= 0) {
        MikCall.sendErrorTips('文件夹名称不能为空');
        return;
      }

      let iRet = await extsObj.mcpost('createFolder', {
        fsid: fileState.fsid,
        fspath: contactFolderName(fileState.curPath, folderName),
      });
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        return;
      }
      MikCall.sendSuccTips('创建成功');
      await fileOp.tryRefreshFiles();
    },

    getViewUrl(path) {
      return extsObj.serverUrl(`view/${fsid}/${path}`);
    },

    getDownloadUrl(path) {
      return extsObj.serverUrl(`download/${fsid}/${path}`);
    },

    getTypeSvgData(fileType) {
      let svgData = FileUtil.typeDescs[fileType].svg;
      return `data:image/svg+xml;utf8,${svgData}`;
    },

    startSelectMode(fileInfo) {
      fileState.curSelected = [fileInfo.name];
      fileState.isSelectMode = true;
    },

    quitSelectMode() {
      if (fileState.isSelectMode) {
        fileState.curSelected = [];
        fileState.isSelectMode = false;
      }
    },

    SelectAll() {
      let ret = [];
      for (let fileInfo of fileState.curFiles) {
        ret.push(fileInfo.name);
      }
      fileState.curSelected = ret;
    },

    SelectNone() {
      fileState.curSelected = [];
    },

    toggleSelectAll() {
      let selectLen = fileState.curSelected.length;
      if (selectLen == fileState.curFiles.length) fileOp.SelectNone();
      else fileOp.SelectAll();
    },

    toggleSelect(fileInfo) {
      let fileName = fileInfo.name;
      if (fileState.curSelected.includes(fileName)) {
        fileState.curSelected = fileState.curSelected.filter(
          (word) => word != fileName
        );
      } else {
        fileState.curSelected.push(fileName);
      }
    },

    async multThumb() {
      if (fileState.curSelected.length <= 0) {
        MikCall.sendErrorTips('请先选择文件');
        return;
      }
      let [isOk, inputVal] = await MikCall.coMakePrompt(`请输入最大长宽`, '1920');
      if (!isOk) return;
      let maxSize = parseInt(inputVal);
      if (!maxSize) {
        MikCall.sendErrorTips('非法长宽');
        return;
      }
      let curPath = fileState.curPath;
      await openMultOperateDlg({
        title: `批量生成缩略图(生成结果存在对应文件夹的.mnthumbs目录下)`,
        jobKeys: fileState.curSelected,
        jobDescFn: (fileName) => `生成缩略图: ${fileName}`,
        jobProcFn: async function (fileName) {
          if (!fileName) return MikCall.failRet('文件名不能为空');
          return await extsObj.mcpost('genThumb', {
            fsid: fsid,
            fspath: contactFolderName(curPath, fileName),
            maxSize: maxSize,
          });
        },
      });
      await fileOp.tryRefreshFiles();
    },
    async multCopy() {
      if (fileState.curSelected.length <= 0) {
        MikCall.sendErrorTips('请先选择文件');
        return;
      }
      let [isOk, toPath] = await openSelectFsidFolderDlg(
        fsid,
        rootPath,
        '复制到这里'
      );
      if (!isOk) return;
      let curPath = fileState.curPath;
      await openMultOperateDlg({
        title: `批量复制`,
        jobKeys: fileState.curSelected,
        jobDescFn: (fileName) => `复制 ${fileName}`,
        jobProcFn: async function (fileName) {
          if (!fileName) return MikCall.failRet('文件名不能为空');
          // await MikCall.coDelay(2000);
          // return ['手动错误', null];
          return await extsObj.mcpost('copyFile', {
            fsid: fsid,
            fspath: contactFolderName(curPath, fileName),
            topath: toPath,
          });
        },
      });
      await fileOp.tryRefreshFiles();
      // // 强制刷新一下目标目录
      // await cacheStore.tryRefreshFiles(fsid, toPath);
    },
    async multMove() {
      if (fileState.curSelected.length <= 0) {
        MikCall.sendErrorTips('请先选择文件');
        return;
      }
      let [isOk, toPath] = await openSelectFsidFolderDlg(
        fsid,
        rootPath,
        '移动到这里'
      );
      if (!isOk) return;
      let curPath = fileState.curPath;
      await openMultOperateDlg({
        title: `批量移动`,
        jobKeys: fileState.curSelected,
        jobDescFn: (fileName) => `移动 ${fileName}`,
        jobProcFn: async function (fileName) {
          if (!fileName) return MikCall.failRet('文件名不能为空');
          return await extsObj.mcpost('mvFile', {
            fsid: fsid,
            fspath: contactFolderName(curPath, fileName),
            topath: toPath,
          });
        },
      });
      await fileOp.tryRefreshFiles();
      // // 强制刷新一下目标目录
      // await cacheStore.tryRefreshFiles(fsid, toPath);
    },
    async multRemove() {
      if (fileState.curSelected.length <= 0) {
        MikCall.sendErrorTips('请先选择文件');
        return;
      }
      let fileNames = fileState.curSelected.join('、');
      let confirmRet = await MikCall.coMakeConfirm(
        `是否确认删除(${fileNames})?`
      );
      if (!confirmRet) return;
      let curPath = fileState.curPath;
      await openMultOperateDlg({
        title: `批量删除`,
        jobKeys: fileState.curSelected,
        jobDescFn: (fileName) => `删除 ${fileName}`,
        jobProcFn: async function (fileName) {
          if (!fileName) return MikCall.failRet('文件名不能为空');
          // await MikCall.coDelay(2000);
          // return ['手动错误', null];
          return await extsObj.mcpost('removeFile', {
            fsid: fsid,
            fspath: contactFolderName(curPath, fileName),
          });
        },
      });
      await fileOp.tryRefreshFiles();
    },
  };

  function OpenRouteSubPath(routeSubPath) {
    routeSubPath = routeSubPath || '';
    let newPath = contactFolderName(rootPath, routeSubPath);
    fileOp.tryRefreshFiles(newPath, true);
  }

  onMounted(() => {
    // console.log('FileHelper onMounted');
    let newPath = rootPath;
    if (extraConf.isListenRouteChange) {
      const route = useRoute();
      OpenRouteSubPath(route.params.routeSubPath);
    } else {
      if (extraConf.initFsrela) {
        newPath = contactFolderName(newPath, extraConf.initFsrela);
      }
      fileOp.tryRefreshFiles(newPath, true);
    }
  });

  let _scrollY;
  onDeactivated(() => {
    // console.log('onDeactivated');
    _scrollY = window.window.scrollY;
  });
  onActivated(() => {
    // console.log('onActivated');
    if (_scrollY) {
      setTimeout(() => {
        window.scrollTo({ top: _scrollY });
      }, 500);
    }
    if (fileState.curPath) {
      fileOp.tryRefreshFiles(undefined, false);
    }
  });

  if (extraConf.isListenRouteChange) {
    onBeforeRouteUpdate((to, from) => {
      if (extraConf.isRememberPos) {
        cacheStore.updateFsPosition(fsid, fileState.curPath, window.scrollY);
      }
      if (to.params.routeSubPath != from.params.routeSubPath) {
        OpenRouteSubPath(to.params.routeSubPath);
      }
    });
    onBeforeRouteLeave(() => {
      if (extraConf.isRememberPos) {
        cacheStore.updateFsPosition(fsid, fileState.curPath, window.scrollY);
      }
    });
  }

  return { fileState, fileGetter, fileOp };
}

// ----------------------------- 预览文件 ---------------------

export function viewFile(fsid, initFilePath, fileList) {
  Dialog.create({
    component: MdcFileViewContainerDlg,
    componentProps: {
      fsid,
      initFilePath,
      fileList,
    },
  });
}

export function downloadFile(fsid, fspath) {
  const extsObj = useExtension();
  openURL(extsObj.serverUrl(`download/${fsid}/${fspath}`));
}

export async function openSelectFsidFolderDlg(fsid, rootPath, confirmLabel) {
  return await MikCall.coCreateDialog({
    component: MdcFolderSelectDlg,
    componentProps: {
      fsid: fsid,
      rootPath: rootPath || '',
      confirmLabel: confirmLabel,
    },
  });
}

export async function openSelectFsidFileDlg(fsid, rootPath, confirmLabel) {
  return await MikCall.coCreateDialog({
    component: MdcFileSelectDlg,
    componentProps: {
      fsid: fsid,
      rootPath: rootPath || '',
      confirmLabel: confirmLabel,
    },
  });
}
