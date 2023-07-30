// 共享给外部可用的变量或组件
import { defineAsyncComponent } from 'vue';

export { openTextCopyDlg, openMultOperateDlg } from './helpers';

export {
  FormTypes,
  DataRule,
  coOpenFormDlg,
  // ConfsHelper,
} from './helpers/FormHelper';

export { default as useMikLoading } from './compositions/useMikLoading';
export { default as useFormView } from './compositions/useFormView';

export { default as PageMenuItem } from './components/PageMenuItem.vue';
export { default as PageSubMenu } from './components/PageSubMenu.vue';
export { default as SimpleNestRouterView } from './layouts/SimpleNestRouterView.vue';
export { default as ExtensionPage } from './layouts/ExtensionPage.vue';
export { default as ResizeDrawer } from './layouts/ResizeDrawer.vue';

export const ByteMdInput = defineAsyncComponent(() =>
  import('miknas/exts/Official/components/markdown/ByteMdInput.vue')
);
export const ByteMdView = defineAsyncComponent(() =>
  import('miknas/exts/Official/components/markdown/ByteMdView.vue')
);
export const VideoPlayer = defineAsyncComponent(() =>
  import('miknas/exts/Official/components/videojs/VideoPlayer.vue')
);
export { default as MdcAceEditor } from './components/ace/MdcAceEditor.vue'
