import {
  DataRule,
  FormTypes,
  coOpenFormDlg,
} from 'miknas/exts/Official/shares';
import useExtension from './extMain';
import { JSEncrypt } from 'jsencrypt';
import { MikCall } from 'miknas/utils';

export function rsaEncrypt(pubKey, msg) {
  let encrypt = new JSEncrypt();
  encrypt.setPublicKey(pubKey);
  return encrypt.encrypt(msg);
}

const ModifyNicknameFormData = [
  {
    id: 'name',
    title: '昵称',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
    },
    default: '',
    desc: `用户昵称`,
    rules: [DataRule.isNotEmptyString],
  },
];

export async function modifyNickname(prevVal, sucCb) {
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: ModifyNicknameFormData,
    initData: {
      name: prevVal || '',
    },
    title: '修改昵称',
  });
  if (!isOk) return;
  const extsObj = useExtension();
  let iRet = await extsObj.mcpost('modifyNickname', newFormData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let result = iRet.ret;
  if (sucCb) {
    sucCb(result);
  }
}

const ModifyPwdFormData = [
  {
    id: 'oldPwd',
    title: '旧密码',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
      type: 'password',
    },
    default: '',
    rules: [DataRule.minLenght(3)],
  },
  {
    id: 'newPwd1',
    title: '新密码',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
      type: 'password',
    },
    default: '',
    rules: [DataRule.minLenght(3)],
  },
  {
    id: 'newPwd2',
    title: '再次输入新密码',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
      type: 'password',
    },
    default: '',
    rules: [DataRule.minLenght(3)],
  },
];

export async function modifyPassword(sucCb, initData) {
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: ModifyPwdFormData,
    initData: initData || {},
    title: '修改密码',
  });
  if (!isOk) return;
  const extsObj = useExtension();
  if (newFormData.newPwd1 != newFormData.newPwd2) {
    MikCall.sendErrorTips('两次新密码输入不一致');
    modifyPassword(sucCb, newFormData);
    return;
  }

  let iRet = await extsObj.mcpost('querySecretToken', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let ret = iRet.ret;

  let token = ret.token;
  let postData = {
    oldPwd: rsaEncrypt(token, newFormData.oldPwd),
    newPwd: rsaEncrypt(token, newFormData.newPwd1),
  };

  iRet = await extsObj.mcpost('modifyPassword', postData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    modifyPassword(sucCb, newFormData);
    return;
  }
  ret = iRet.ret;
  if (sucCb) {
    sucCb(ret);
  }
}
