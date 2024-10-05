<template>
  <q-page>
    <q-toolbar>
      <q-toolbar-title class="mn-toolbar-title">
        查看所有文件分享
      </q-toolbar-title>
      <q-btn
        class="q-ml-sm"
        color="primary"
        :disable="state.isLoading"
        label="刷新"
        @click="tryRefreshTable()"
      />
    </q-toolbar>
    <q-list>
      <q-item
        v-for="shareInfo in state.tableRows"
        :key="shareInfo.sid"
        v-ripple
        tag="a"
        target="_blank"
        :href="viewShareRoute(shareInfo)"
        clickable
      >
        <q-item-section avatar>
          <q-avatar>
            <q-icon name="link" />
          </q-avatar>
        </q-item-section>

        <q-item-section>
          <q-item-label lines="1" class="text-weight-bold">{{
            shareInfo.name
          }}</q-item-label>
          <q-item-label
            v-if="shareInfo.hint"
            caption
            lines="2"
            class="text-weight-bold"
          >
            {{ shareInfo.hint }}
          </q-item-label>
          <q-item-label v-if="shareInfo.isExpired" class="text-red" caption
            >有效期至: {{ shareInfo.viewEts }}</q-item-label
          >
          <q-item-label v-else caption class="text-green"
            >有效期至: {{ shareInfo.viewEts }}</q-item-label
          >
        </q-item-section>

        <q-item-section side top>
          <q-item-label v-if="shareInfo.isExpired" class="text-red" caption
            >已过期</q-item-label
          >
          <q-btn dense flat round icon="more_vert" @click.stop.prevent="">
            <q-menu auto-close>
              <q-list>
                <q-item clickable @click="tryShare(shareInfo)">
                  <q-item-section>查看分享链接</q-item-section>
                </q-item>
                <q-item clickable @click="tryRemove(shareInfo)">
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
import { viewShareRoute } from '../FileHelper.js';
import { openTextCopyDlg } from 'miknas/exts/Official/shares';

const state = reactive({
  tableRows: [],
  isLoading: true,
});

const extsObj = useExtension();
async function tryRefreshTable() {
  let now = Date.now();
  state.isLoading = true;

  let iRet = await extsObj.mcpost('queryShares', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let tableRows = iRet.ret;
  for (let info of tableRows) {
    if (info.intv && info.intv > 0) {
      let endTs = info.bts + info.intv;
      info.viewEts = gutil.formatTs(endTs * 1000);
      if (endTs * 1000 <= now) {
        info.isExpired = true;
      }
    }
  }
  tableRows.sort((a, b) => a.bts - b.bts);
  state.tableRows = tableRows;
  state.isLoading = false;
}

function tryShare(info) {
  openTextCopyDlg({ title: '分享链接如下:', txt: viewShareRoute(info, true) });
}

async function tryRemove(info) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除?`);
  if (!isOk) return;
  state.isLoading = true;
  let iRet = await extsObj.mcpost('removeShare', { sid: info.sid });
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
