import {
  DataRule,
  FormTypes,
  coOpenFormDlg,
} from 'miknas/exts/Official/shares';
import { MikCall } from 'miknas/utils';
import useExtension from './extMain';

function getModifyBookmarkFormData(hintList) {
  return [
    {
      id: 'kind',
      title: '分类',
      component: FormTypes.MdcTextAutoComplete,
      componentProps: {
        filled: true,
      },
      default: '',
      rules: [DataRule.isNotEmptyString],
      hintList: hintList || [],
    },
    {
      id: 'name',
      title: '名称',
      component: FormTypes.MdcTextInput,
      componentProps: {
        filled: true,
      },
      default: '',
      rules: [DataRule.isNotEmptyString],
    },
    {
      id: 'url',
      title: '链接地址',
      component: FormTypes.MdcTextInput,
      componentProps: {
        filled: true,
      },
      default: '',
      rules: [DataRule.isUrl],
    },
    {
      id: 'icon',
      title: '图标地址',
      component: FormTypes.MdcTextInput,
      componentProps: {
        filled: true,
      },
      default: '',
      rules: [DataRule.emptyOr([DataRule.isUrl])],
      desc: '显示在页面的图标',
      helpActions: [
        {
          label: '使用favio.icon',
          func: (formData) => {
            let urlobj;
            try {
              urlobj = new URL(formData.url);
            } catch (error) {
              MikCall.sendErrorTips(`${formData.url}不是一个有效的链接地址`);
            }
            if (urlobj) formData.icon = `${urlobj.origin}/favicon.ico`;
          },
        },
      ],
    },
  ];
}
export async function openAddBookmark(hintList, sucCb) {
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: getModifyBookmarkFormData(hintList),
    initData: {},
    title: '添加书签',
  });
  if (!isOk) return;
  const extsObj = useExtension();
  let iRet = await extsObj.mcpost('add', newFormData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  if (sucCb) {
    sucCb(iRet.ret);
  }
}

export async function openModifyBookmark(initData, hintList, sucCb) {
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: getModifyBookmarkFormData(hintList),
    initData: initData,
    title: '修改书签',
  });
  if (!isOk) return;
  const extsObj = useExtension();
  newFormData.id = initData.id;
  let iRet = await extsObj.mcpost('modify', newFormData);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  if (sucCb) {
    sucCb(iRet.ret);
  }
}

export async function deleteBookmark(initData, sucCb) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除书签(${initData.name})?`);
  if (!isOk) return;
  let bmid = initData.id;
  if (!bmid) {
    MikCall.sendErrorTips('非法书签id');
    return;
  }
  const extsObj = useExtension();
  let iRet = await extsObj.mcpost('delete', { id: bmid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  if (sucCb) {
    sucCb(iRet.ret);
  }
}
