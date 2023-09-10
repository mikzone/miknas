import {
  DataRule,
  FormTypes,
  coOpenFormDlg,
} from 'miknas/exts/Official/shares';
import { MikCall, MyAes } from 'miknas/utils';
import SecretAttach from './components/attaches/SecretAttach.vue';
import { useNoteStore } from './stores/note';

const noteStore = useNoteStore();
// 所有附件其实能改的也就是它自己data的那块数据，因为这是个json对象，所以我们可以直接用jsonData(已经预处理了)
// formData 是传到form那边的数据
export const attachCfgs = {
  recret: {
    name: '加密数据',
    component: SecretAttach,
    getAttachName: (jsonData) => jsonData.name,
    genFormData: async (jsonData) => {
      let ret = {
        name: jsonData.name,
        hint: jsonData.hint,
        pwd: '',
        txt: '',
      };
      if (jsonData.encrypt) {
        let [isOk, pwd] = await MikCall.coMakePrompt(
          `提示: ${jsonData.hint || '无'}`,
          '',
          '请先输入解密密码',
          'password'
        );
        if (!isOk) {
          return;
        }
        if (!pwd) {
          MikCall.sendErrorTips('密码不能为空');
          return;
        }
        let myaes = new MyAes(pwd);
        let [decrpytMsg, decryptErr] = myaes.decryptEx(jsonData.encrypt);
        if (!decrpytMsg) {
          MikCall.sendErrorTips(decryptErr);
          return;
        }
        ret.pwd = pwd;
        ret.txt = decrpytMsg;
      }
      return ret;
    },
    genJsonData: (formData) => {
      if (formData.encrypt && !formData.pwd) {
        MikCall.sendErrorTips('密码不能为空');
        return;
      }
      let ret = {
        name: formData.name,
        hint: formData.hint,
      };
      if (formData.txt) {
        let myaes = new MyAes(formData.pwd);
        let newTxt = myaes.encryptEx(formData.txt);
        if (!newTxt) {
          MikCall.sendErrorTips('加密过程发生错误');
          return;
        }
        ret.encrypt = newTxt;
      }
      return ret;
    },
    getFormConfs: () => {
      return [
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
          id: 'txt',
          title: '待加密内容',
          component: FormTypes.MdcMarkdown,
          componentProps: {
            // filled: true,
            // type: 'textarea',
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
          title: '提示信息',
          component: FormTypes.MdcTextInput,
          componentProps: {
            filled: true,
          },
          default: '',
          desc: `你可以记录与解密密码或者加密数据有关的任何提示信息`,
        },
      ];
    },
  },
};

export async function openAddAttach(noteInfo, type, createFormArgs, sucCb) {
  let noteid = noteInfo.id;
  if (!noteid) return;
  let attachConf = attachCfgs[type];
  if (!attachConf) return;
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: attachConf.getFormConfs(createFormArgs),
    initData: {},
    title: `添加${attachConf.name}`,
  });
  if (!isOk) return;

  const newJsonData = await attachConf.genJsonData(newFormData);
  if (!newJsonData) return;
  const callArgs = {
    itemId: noteid,
    type: type,
    data: JSON.stringify(newJsonData),
  };
  let ret = await noteStore.addAttach(callArgs);
  if (ret) {
    noteInfo.noteAttachs.push(ret);
    if (sucCb) {
      sucCb(ret);
    }
  }
}

export async function openModifyAttach(initData, createFormArgs, sucCb) {
  if (!initData || !initData.id) return;
  let attachConf = attachCfgs[initData.type];
  if (!attachConf) return;
  let formData = await attachConf.genFormData(initData.jsonData);
  if (!formData) return;
  let [isOk, newFormData] = await coOpenFormDlg({
    formConfs: attachConf.getFormConfs(createFormArgs),
    initData: formData,
    title: `修改${attachConf.name}`,
  });
  if (!isOk) return;
  const newJsonData = await attachConf.genJsonData(newFormData);
  if (!newJsonData) return;
  const callArgs = {
    id: initData.id,
    itemId: initData.itemId,
    type: initData.type,
    data: JSON.stringify(newJsonData),
  };
  let ret = await noteStore.modifyAttach(callArgs);
  if (ret && sucCb) {
    sucCb(ret);
  }
}

export function deleteAttach(initData, sucCb) {
  if (!initData || !initData.id) return;
  let attachConf = attachCfgs[initData.type];
  if (!attachConf) return;
  let attachName = attachConf.getAttachName(initData.jsonData);
  MikCall.makeConfirm(`是否确认删除附件(${attachName})?`, async () => {
    let attachId = initData.id;
    if (!attachId) {
      MikCall.sendErrorTips('非法附件id');
      return;
    }
    let ret = await noteStore.deleteAttach(attachId);
    if (ret && sucCb) {
      sucCb(ret);
    }
  });
}

export function useNoteInfoUtil(computeNoteInfo) {
  const noteUtil = {
    onAttachChange(id, attachInfo) {
      let noteInfo = computeNoteInfo.value;
      let newAttaches = noteInfo.noteAttachs.filter((obj) => obj.id !== id);
      if (attachInfo) {
        newAttaches.push(attachInfo);
      }
      noteInfo.noteAttachs = newAttaches;
    },
    async startAddAttach(type) {
      let noteInfo = computeNoteInfo.value;
      await openAddAttach(noteInfo, type, null);
    },
    async startModifyAttach(attachInfo) {
      await openModifyAttach(attachInfo, null, (ret) => {
        noteUtil.onAttachChange(attachInfo.id, ret);
      });
    },
    async startDeleteAttach(attachInfo) {
      await deleteAttach(attachInfo, () => {
        noteUtil.onAttachChange(attachInfo.id, null);
      });
    },
  };
  return noteUtil;
}

export async function moveNote(noteInfo) {
  let [isOk, folderId] = await noteStore.selectFolder('移动到');
  if (!isOk) return;
  if (!folderId) {
    MikCall.sendErrorTips('文件夹id错误');
    return;
  }
  if (folderId == noteInfo.folder) {
    MikCall.sendErrorTips('文件夹未改变');
    return;
  }
  let newNoteInfo = await noteStore.saveNote({
    id: noteInfo.id,
    folder: folderId,
  });
  return newNoteInfo;
}
