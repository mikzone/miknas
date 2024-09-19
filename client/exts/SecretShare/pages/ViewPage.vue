<template>
  <q-page>
    <q-layout
      v-if="state.decrpytMsg"
      container
      view="hHh Lpr fff"
      class="absolute-full"
    >
      <q-header>
        <q-toolbar class="bg-white text-black">
          <q-toolbar-title class="mn-toolbar-title">{{
            state.info.name
          }}</q-toolbar-title>
          <q-tabs v-model="state.mode" shrink>
            <q-tab name="raw" no-caps label="原文本" />
            <q-tab name="md" no-caps label="格式化" />
          </q-tabs>
        </q-toolbar>
        <q-separator />
      </q-header>
      <q-footer class="bg-white text-black">
        <q-separator />
        <div class="row q-pa-sm q-gutter-sm justify-end">
          <q-btn dense color="primary" @click="tryCopy">复制原文</q-btn>
          <q-btn dense color="primary" @click="trySave">下载内容文件</q-btn>
        </div>
      </q-footer>
      <q-page-container>
        <q-page padding>
          <div
            v-if="state.mode == 'raw'"
            class="q-ma-none"
            style="white-space: pre-wrap"
          >
            {{ state.decrpytMsg }}
          </div>
          <ByteMdView v-else :model-value="state.decrpytMsg" />
        </q-page>
      </q-page-container>
    </q-layout>
    <q-card v-else-if="state.info" class="absolute-full column">
      <q-card-section>
        <div class="text-h6">{{ state.info.uid }}给你发了一条密文</div>
        <div class="text-subtitle2">
          该密文将于 {{ state.info.viewTs }} 过期
        </div>
      </q-card-section>
      <q-separator />

      <q-card-section class="col column flex-center">
        <div class="text-h5">需要解密后才能查看该密文</div>
        <q-btn class="q-mt-lg" color="primary" @click="tryDecrypt">解密</q-btn>
      </q-card-section>
    </q-card>
    <div
      v-else-if="state.errInfo"
      class="text-red absolute-full column flex-center"
    >
      {{ state.errInfo.why }}
    </div>
    <div v-else class="absolute-full column flex-center">加载中</div>
  </q-page>
</template>
<script setup>
import { exportFile } from 'quasar';
import { gutil, MikCall, MyAes } from 'miknas/utils';
import { onMounted, reactive } from 'vue';
import { useExtension } from '../extMain.js';
import { openUnlockSecretDlg } from '../helper';
import { ByteMdView } from 'miknas/exts/Official/shares';

const props = defineProps({
  mid: {
    type: String,
    required: true,
  },
});

const state = reactive({
  mode: 'raw',
  info: undefined,
  decrpytMsg: undefined,
  errInfo: undefined,
});

const extsObj = useExtension();
async function tryQuery() {
  let iRet = await extsObj.mcpost('viewOne', { mid: props.mid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.errInfo = iRet;
    return;
  }
  let info = iRet.ret;
  info.viewTs = gutil.formatTs(info.ts * 1000);
  state.info = info;
}

async function tryDecrypt() {
  let [isOk, formData] = await openUnlockSecretDlg(state.info.hint);
  if (!isOk) return;
  let pwd = formData.pwd;
  if (!pwd) {
    MikCall.sendErrorTips('密码不能为空!');
    return;
  }
  let myaes = new MyAes(pwd);
  let [decrpytMsg, decryptErr] = myaes.decryptEx(state.info.txt);
  if (!decrpytMsg) {
    MikCall.sendErrorTips(decryptErr);
    return;
  }
  state.decrpytMsg = decrpytMsg;
}

function tryCopy() {
  if (!state.decrpytMsg) return MikCall.sendErrorTips('需要先解密');
  gutil
    .copyToClipboard(state.decrpytMsg)
    .then(() => {
      MikCall.sendSuccTips('复制成功!');
    })
    .catch(() => {
      MikCall.sendErrorTips('复制失败');
    });
}

function trySave() {
  if (!state.decrpytMsg) return MikCall.sendErrorTips('需要先解密');
  let fileName = props.mid;
  if (state.info.name) {
    fileName = state.info.name;
  }
  const status = exportFile(`${fileName}.txt`, state.decrpytMsg);
  if (status !== true) {
    MikCall.sendErrorTips(`保存失败: ${status}`);
  }
}

onMounted(() => {
  tryQuery();
});
</script>
