import TextCopyDlg from './components/dialogs/TextCopyDlg.vue';
import MultOperateDlg from './components/dialogs/MultOperateDlg.vue';
import { MikCall } from 'miknas/utils';

export function openTextCopyDlg(props) {
  return MikCall.coCreateDialog({
    component: TextCopyDlg,
    componentProps: props,
  })
}

export function openMultOperateDlg(props) {
  return MikCall.coCreateDialog({
    component: MultOperateDlg,
    componentProps: props,
  })
}
