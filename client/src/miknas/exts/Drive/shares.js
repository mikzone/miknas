import { defineAsyncComponent } from 'vue';

export {
  FileUtil,
  useFileView,
  viewFile,
  downloadFile,
  openSelectFsidFolderDlg,
  openSelectFsidFileDlg,
} from './FileHelper';

export { default as MdcFileViewText } from './components/FileView/MdcFileViewText.vue'
export const MdcFileViewVideo = defineAsyncComponent(() => import('./components/FileView/MdcFileViewVideo.vue'))
export const MdcFileViewImg = defineAsyncComponent(() => import('./components/FileView/MdcFileViewImg.vue'))

export { default as usePreviewView } from './helpers/usePreviewView';

export { default as MdcDriveAliveView } from './components/MdcDriveAliveView.vue';
export { default as CommonExplorerPage } from './pages/CommonExplorerPage.vue';
export { default as MySharesPage } from './pages/MyShares.vue';
