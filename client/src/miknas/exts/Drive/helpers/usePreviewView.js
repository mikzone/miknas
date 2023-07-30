import { useExtension } from 'miknas/exts/Drive/extMain';
import { onMounted, reactive } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { downloadFile, FileUtil } from 'miknas/exts/Drive/shares';
import { computed } from 'vue';
import { gutil, MikCall } from 'miknas/utils';
import { useViewStore } from '../stores/view';

const viewStore = useViewStore();

export default function usePreviewView({fsid, initFilePath}) {

  if (!fsid) return;
  if (!initFilePath) return;

  let fileList = [];
  if (viewStore.fsid == fsid) {
    fileList = gutil.jsonCopy(viewStore.fileList)
  }

  const extsObj = useExtension();

  const viewState = reactive({
    curFilePath: '',
    mode: 'unkown',
    curFileStat: {},
    modeDesc: gutil.jsonCopy(FileUtil.typeDescs),
    loadingMgr: useMikLoading(),
  });

  async function updateCurFile(filePath, mode) {
    if (!filePath) return;
    if (!mode) mode = FileUtil.getFileType(filePath);
    if (filePath == viewState.curFilePath && mode == viewState.mode) return;
    viewState.curFileStat = {};
    viewState.curFilePath = filePath;
    viewState.mode = mode;
    let stateName = `正在请求${filePath}`;
    viewState.loadingMgr.addLoadingState(stateName);
    // await MikCall.coDelay(2000);
    let iRet = await extsObj.mcpost('queryFileInfo', { fsid: fsid, fspath: filePath });
    if (!iRet.suc) {
      MikCall.alertRespErrMsg(iRet);
      viewState.loadingMgr.removeLoadingState(stateName);
      return;
    }
    let ret = iRet.ret;

    ret.viewSize = FileUtil.formatSize(ret.size);
    ret.viewModify = FileUtil.formatTs(ret.modify * 1000);
    viewState.curFileStat = ret;
    viewState.curFilePath = filePath;
    viewState.mode = mode;
    viewState.loadingMgr.removeLoadingState(stateName);
  }

  onMounted(() => {
    updateCurFile(initFilePath);
  });

  function GetNextFile(idx, step, fileTypes) {
    let nextIdx = idx + step;
    let nextFile = fileList[nextIdx];
    while (nextFile) {
      if (!fileTypes) return nextFile;
      let ftype = FileUtil.getFileType(nextFile);
      if (fileTypes.includes(ftype)) return nextFile;
      nextIdx = nextIdx + step;
      nextFile = fileList[nextIdx];
    }
  }

  const viewGetter = {
    curModeIcon: computed(() => {
      let mode = viewState.mode;
      let desc = FileUtil.typeDescs[mode];
      if (!desc) return 'question_mark';
      return desc.icon || 'preview';
    }),
    picIdx: computed(() => {
      let ret = 0;
      for (let url of fileList) {
        if (url == viewState.curFilePath) return ret;
        ret += 1;
      }
      return null;
    }),
  }

  const viewOp = {
    chooseNext(step, fileTypes) {
      let idx = viewGetter.picIdx.value;
      if (idx == null) {
        MikCall.sendErrorTips('文件定位失败');
        return;
      }
      // let newFile = props.fileList[idx + step];
      let newFile = GetNextFile(idx, step, fileTypes)
      if (!newFile) {
        MikCall.sendErrorTips('已到尽头');
        return;
      }
      updateCurFile(newFile);
    },
    downloadCurrent() {
      downloadFile(fsid, viewState.curFilePath);
    },
  };

  return { viewState, viewGetter, viewOp }

}
