import {
  DataRule,
  FormTypes,
  coOpenFormDlg,
} from 'miknas/exts/Official/shares';
import { gutil, MikCall, MyAes } from 'miknas/utils';
import useExtension from './extMain';

const FormData = [
  {
    id: 'name',
    title: '标题',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
    },
    default: '',
    rules: [DataRule.isNotEmptyString],
  },
  {
    id: 'content',
    title: '信息内容',
    component: FormTypes.MdcMarkdown,
    componentProps: {
      // byteMdOpts: {
      //   mode: 'tab',
      // },
    },
    default: '',
    desc: `该信息会使用下方的密码在前端进行加密后发送，服务端不会明文知道你的内容`,
    rules: [DataRule.isNotEmptyString],
  },
  {
    id: 'pwd',
    title: '密码',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
      type: 'password',
    },
    default: '',
    desc: `!!!注意: 你要自己记住好该密码，服务端不会记录该数据，因此如果忘记密码，那就没办法了`,
    rules: [DataRule.isNotEmptyString],
  },
  {
    id: 'hint',
    title: '与解密密码有关的提示信息',
    component: FormTypes.MdcTextInput,
    componentProps: {
      filled: true,
    },
    default: '',
    desc: `被分享的人也会看到密码提示信息，因此你可以和上面密码结合起来设计成一个问答`,
  },
  {
    id: 'intv',
    title: '消息有效期',
    component: FormTypes.MdcSelect,
    componentProps: {
      filled: true,
    },
    selectOptions: [
      { label: '5分钟', value: 5 * 60 },
      { label: '10分钟', value: 10 * 60 },
      { label: '30分钟', value: 30 * 60 },
      { label: '1小时', value: 1 * 60 * 60 },
      { label: '12小时', value: 12 * 60 * 60 },
      { label: '1天', value: 1 * 24 * 60 * 60 },
      { label: '3天', value: 3 * 24 * 60 * 60 },
      { label: '7天', value: 7 * 24 * 60 * 60 },
      { label: '14天', value: 14 * 24 * 60 * 60 },
      { label: '31天', value: 31 * 24 * 60 * 60 },
    ],
    default: 30 * 60,
    desc: `超过有效期的消息将无法读取`,
  },
];

export async function createNewSecret(initData) {
  let [isOk, secretObj] = await coOpenFormDlg({
    formConfs: FormData,
    initData: initData,
  });
  if (!isOk) return;
  let { content, pwd, ...others } = secretObj;
  let myaes = new MyAes(pwd);
  let newTxt = myaes.encryptEx(content);
  if (!newTxt) {
    MikCall.sendErrorTips('加密过程发生错误');
    return;
  }
  let reqData = { txt: newTxt, ...others };
  const extsObj = useExtension();
  let iRet = await extsObj.mcpost('addSecret', reqData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return await createNewSecret(secretObj);
  }
  return iRet.ret;
}

async function decrpytSecret(secretObj, canRetry) {
  let [isOk, formData] = await openUnlockSecretDlg(secretObj.hint);
  if (!isOk) return;
  let mypwd = formData.pwd;
  if (!mypwd) {
    MikCall.sendErrorTips('密码不能为空!');
    if (canRetry) return await decrpytSecret(secretObj, canRetry);
    else return;
  }

  let myaes = new MyAes(mypwd);
  let [decrpytMsg, decryptErr] = myaes.decryptEx(secretObj.txt);
  if (!decrpytMsg) {
    MikCall.sendErrorTips(decryptErr);
    if (canRetry) return await decrpytSecret(secretObj, canRetry);
    else return;
  }

  let { content, txt, pwd, ...others } = secretObj;
  return {
    content: decrpytMsg,
    pwd: mypwd,
    ...others
  }
}

export async function modifySecret(decryptData) {
  let [isOk, secretObj] = await coOpenFormDlg({
    formConfs: FormData,
    initData: decryptData,
  });
  if (!isOk) return;
  let { content, pwd, ...others } = secretObj;
  let myaes = new MyAes(pwd);
  let newTxt = myaes.encryptEx(content);
  if (!newTxt) {
    MikCall.sendErrorTips('加密过程发生错误');
    return await modifySecret(secretObj);
  }
  let reqData = { txt: newTxt, mid: decryptData.mid, ...others };
  const extsObj = useExtension();
  let iRet = await extsObj.mcpost('modifySecret', reqData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return await modifySecret(secretObj);
  }
  return iRet.ret;
}

export async function decryptAndModifySecret(initData) {
  let decryptData = await decrpytSecret(initData, true)
  if (!decryptData) return;
  return await modifySecret(decryptData);
}

export function viewSecretRoute(info, isfull) {
  const extsObj = useExtension();
  const routeLocate = extsObj.routePath(`view/${info.mid}`);
  if (!isfull) return routeLocate;
  return gutil.routeFullUrl(routeLocate);
}

export async function openUnlockSecretDlg(passHint) {
  let UnlockFormData = [
    {
      id: 'pwd',
      title: '解密密码提示信息:',
      component: FormTypes.MdcTextInput,
      componentProps: {
        filled: true,
        type: 'password',
      },
      default: '',
      desc: passHint || '无',
      rules: [DataRule.isNotEmptyString],
    },
  ];

  return await coOpenFormDlg({
    formConfs: UnlockFormData,
    title: '输入解密密码',
  });
}
