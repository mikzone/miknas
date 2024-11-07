<template>
  <div>
    <q-tabs v-model="state.tab" class="bg-green-10 text-white shadow-2" :breakpoint="0">
      <q-tab name="spaces" label="工作区" />
      <q-tab name="plugins" label="所有插件" />
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
    <q-list v-show="state.tab === 'plugins'" separator bordered>
      <q-item v-for="defInfo in state.plugins" :key="defInfo.Id">
        <q-item-section avatar>
          <q-icon name="settings_input_component" />
        </q-item-section>
        <q-item-section>
          <q-item-label>
            {{ defInfo.Id }}
            <q-badge color="grey" :label="defInfo.Version" />
          </q-item-label>
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
  plugins: {},
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
  state.plugins = result.plugins;
  loadingMgr.removeLoadingState(stateName);
}

onMounted(() => {
  queryResult();
});
</script>
