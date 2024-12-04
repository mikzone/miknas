import MdcCmdExecResultDlg from './components/cmd/MdcCmdExecResultDlg.vue';
import { MikCall } from 'miknas/utils';

export async function coFetchResult({ jobId, aceLang }) {
  return await MikCall.coCreateDialog({
    component: MdcCmdExecResultDlg,
    componentProps: {
      jobId: jobId,
      aceLang: aceLang,
    }
  })
}
