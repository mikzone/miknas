<template>
  <div>
    <q-tabs v-model="state.tab" class="bg-green-10 text-white shadow-2" :breakpoint="0">
      <q-tab name="spaces" label="工作区" />
      <q-tab name="defs" label="管理模板" />
    </q-tabs>
    <q-list v-show="state.tab === 'spaces'" separator bordered>
      <q-item
        v-for="space in state.spaces"
        :key="space.Id"
        clickable
        :to="extsObj.routePath(`list/${space.Id}`)"
      >
        <q-item-section avatar>
          <q-icon name="folder_open" />
        </q-item-section>
        <q-item-section>
          <q-item-label lines="1">
            <q-badge color="secondary" :label="space.Id" />
          </q-item-label>
          <q-item-label>{{ space.Name }}</q-item-label>
          <q-item-label caption>{{ space.Path }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
    <q-list v-show="state.tab === 'defs'" separator bordered>
      <q-item v-for="defInfo in state.defs" :key="defInfo.DefId">
        <q-item-section avatar>
          <q-icon name="folder_open" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ defInfo.DefId }}</q-item-label>
          <q-item-label caption>{{ defInfo.Desc }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
    <q-inner-loading :showing="loadingMgr.isloading.value" :label="loadingMgr.loadingLabel.value" />
  </div>
</template>
<script setup>
import { onMounted, reactive } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { MikCall } from 'miknas/utils';

const loadingMgr = useMikLoading();

const state = reactive({
  result: {},
  spaces: {},
  defs: {},
  tab: 'spaces'
});

const extsObj = useExtension();

async function queryResult() {
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);

  let iRet = await extsObj.mcpost('maininfo');
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let result = iRet.ret;
  state.result = result;
  state.spaces = result.spaces;
  state.defs = result.defs;
  loadingMgr.removeLoadingState(stateName);
}

onMounted(() => {
  queryResult();
});
</script>
