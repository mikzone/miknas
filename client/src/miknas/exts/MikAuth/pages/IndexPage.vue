<template>
  <q-page padding class="row">
    <q-space />
    <div style="width: 600px; max-width: 90%">
      <q-card v-if="officialStore.uid" class="no-shadow mn-bordered">
        <q-list bordered>
          <q-item>
            <q-item-section>
              <q-item-label>UID</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-item-label caption>{{ state.userinfo.uid }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-item-label>角色</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-item-label caption>{{ state.userinfo.role }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-item v-if="state.userinfo.role != state.userinfo.realRole">
            <q-item-section>
              <q-item-label>当前实际角色</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-item-label caption>{{ state.userinfo.realRole }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-item-label>昵称</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-item-label caption>{{ state.userinfo.name }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
        <q-card-actions>
          <q-btn color="primary" @click="modifyNickname(state.userinfo.name, refresh)">修改昵称</q-btn>
          <q-btn color="primary" @click="modifyPassword(refresh)">修改密码</q-btn>
        </q-card-actions>
      </q-card>
      <div v-else>当前未登录</div>
    </div>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      :label="loadingMgr.loadingLabel.value"
    />
    <q-space />
  </q-page>
</template>

<script setup>
import { useMikLoading } from 'miknas/exts/Official/shares';
import { MikCall } from 'miknas/utils';
import { onMounted, reactive } from 'vue';
import { useOfficialStore } from 'miknas/exts/Official/stores/official';
import useExtension from '../extMain';
import { modifyNickname, modifyPassword } from '../helpers';
const officialStore = useOfficialStore();
const extsObj = useExtension();
const loadingMgr = useMikLoading();

const state = reactive({
  userinfo: {},
});

async function refresh() {
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);

  let iRet = await extsObj.mcpost('currentUserinfo');
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let userinfo = iRet.ret;
  state.userinfo = userinfo;
  loadingMgr.removeLoadingState(stateName);
}

onMounted(() => {
  refresh();
});
</script>
