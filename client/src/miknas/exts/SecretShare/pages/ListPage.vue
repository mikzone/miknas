<template>
  <q-page>
    <q-toolbar>
      <q-toolbar-title class="mn-toolbar-title"> 查看所有密文 </q-toolbar-title>
      <q-btn
        class="q-ml-sm"
        color="primary"
        :disable="state.isLoading"
        label="刷新"
        @click="tryRefreshTable()"
      />
      <q-btn
        class="q-ml-sm"
        color="primary"
        label="添加密文"
        @click="tryAddSecret"
      />
    </q-toolbar>
    <q-list>
      <q-item
        v-for="secretInfo in state.tableRows"
        :key="secretInfo.mid"
        v-ripple
        :to="viewSecretRoute(secretInfo)"
        clickable
      >
        <q-item-section avatar>
          <q-avatar>
            <q-icon name="folder" />
          </q-avatar>
        </q-item-section>

        <q-item-section>
          <q-item-label lines="1" class="text-weight-bold">{{
            secretInfo.name
          }}</q-item-label>
          <q-item-label
            v-if="secretInfo.hint"
            caption
            lines="2"
            class="text-weight-bold"
          >
            {{ secretInfo.hint }}
          </q-item-label>
          <q-item-label v-if="secretInfo.isExpired" class="text-red" caption
            >有效期至: {{ secretInfo.viewTs }}</q-item-label
          >
          <q-item-label v-else caption class="text-green"
            >有效期至: {{ secretInfo.viewTs }}</q-item-label
          >
        </q-item-section>

        <q-item-section side top>
          <q-item-label v-if="secretInfo.isExpired" class="text-red" caption
            >已过期</q-item-label
          >
          <q-btn dense flat round icon="more_vert" @click.stop.prevent="">
            <q-menu auto-close>
              <q-list>
                <q-item clickable @click="tryShare(secretInfo)">
                  <q-item-section>分享</q-item-section>
                </q-item>
                <q-item clickable @click="tryModifySecret(secretInfo)">
                  <q-item-section>编辑</q-item-section>
                </q-item>
                <q-item clickable @click="tryRemove(secretInfo)">
                  <q-item-section>删除</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-item-section>
      </q-item>
    </q-list>
    <q-inner-loading :showing="state.isLoading" />
  </q-page>
</template>
<script></script>
<script setup>
import { onMounted, reactive } from 'vue';
import { useExtension } from '../extMain.js';
import { gutil, MikCall } from 'miknas/utils';
import {
  createNewSecret,
  viewSecretRoute,
  decryptAndModifySecret,
} from '../helper.js';
import { openTextCopyDlg } from 'miknas/exts/Official/shares';

const state = reactive({
  tableRows: [],
  isLoading: true,
});

const extsObj = useExtension();
async function tryRefreshTable() {
  let now = Date.now();
  state.isLoading = true;

  let iRet = await extsObj.mcpost('querySecrets', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let tableRows = iRet.ret;
  for (let info of tableRows) {
    info.viewTs = gutil.formatTs(info.ts * 1000);
    if (info.ts * 1000 <= now) {
      info.isExpired = true;
    }
  }
  tableRows.sort((a, b) => a.bts - b.bts);
  state.tableRows = tableRows;
  state.isLoading = false;
}

async function tryAddSecret() {
  let result = await createNewSecret();
  if (result) {
    tryRefreshTable();
  }
}

async function tryModifySecret(secretInfo) {
  let result = await decryptAndModifySecret(secretInfo);
  if (result) {
    tryRefreshTable();
  }
}

function tryShare(info) {
  openTextCopyDlg({ title: '分享链接如下:', txt: viewSecretRoute(info, true) });
}

async function tryRemove(info) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除?`);
  if (!isOk) return;
  state.isLoading = true;
  let iRet = await extsObj.mcpost('removeSecret', { mid: info.mid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let result = iRet.ret;
  MikCall.sendSuccTips('删除成功');
  await tryRefreshTable();
  state.isLoading = false;
}

onMounted(() => {
  tryRefreshTable();
});
</script>
