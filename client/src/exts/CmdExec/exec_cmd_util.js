import { Dialog } from 'quasar';
import MdcCmdExecResultDlg from 'miknas/exts/CmdExec/components/MdcCmdExecResultDlg.vue'

export function fetchResult({ jobId, aceLang }) {
  return Dialog.create({
    component: MdcCmdExecResultDlg,
    componentProps: {
      jobId: jobId,
      aceLang: aceLang,
    }
  })
}
