import { Dialog } from 'quasar';
import MdcCmdExecResultDlg from './components/cmd/MdcCmdExecResultDlg.vue';

export function fetchResult({ jobId, aceLang }) {
  return Dialog.create({
    component: MdcCmdExecResultDlg,
    componentProps: {
      jobId: jobId,
      aceLang: aceLang,
    }
  })
}
