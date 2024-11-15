<template>
  <div>
    <q-tabs v-model="state.tab" align="left" class="bg-grey-3 shadow-1" :breakpoint="0">
      <q-tab
        v-for="pluginDetail in relPlugins"
        :key="pluginDetail.pluginId"
        :name="pluginDetail.pluginId"
        :label="pluginDetail.pluginDef.Title"
      />
    </q-tabs>
    <PluginView
      v-for="pluginDetail in relPlugins"
      v-show="state.tab === pluginDetail.pluginId"
      :key="pluginDetail.pluginId"
      :plugin-detail="pluginDetail"
      :space-id="props.spaceId"
      :fspath="props.fspath"
      separator
      bordered
    />
    <q-inner-loading :showing="loadingMgr.isloading.value" :label="loadingMgr.loadingLabel.value" />
  </div>
</template>
<script setup>
import { computed, onMounted, reactive } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { MikCall } from 'miknas/utils';
import { useWorkStore } from '../stores/work';
import PluginView from './PluginView.vue';

const loadingMgr = useMikLoading();
const workStore = useWorkStore();

const props = defineProps({
  spaceId: {
    type: String,
    required: true
  },
  fspath: {
    type: String,
    default: ''
  }
});

const state = reactive({
  pluginStats: [],
  tab: undefined
});

const extsObj = useExtension();

const relPlugins = computed(() => {
  let ret = [];
  for (let pluginStat of state.pluginStats) {
    let pluginId = pluginStat.id;
    let pluginDef = workStore.plugins[pluginId];
    if (pluginDef) {
      ret.push({ pluginId, pluginDef, pluginStat });
    }
  }
  return ret;
});

function RecheckTab() {
  if (state.tab) return;
  for (let pluginDetail of relPlugins.value) {
    let pluginId = pluginDetail.pluginId;
    state.tab = pluginId;
  }
}

async function refreshFolderDetail() {
  let stateName = `正在刷新`;
  loadingMgr.addLoadingState(stateName);

  let iRet = await extsObj.mcpost('queryFolderDetail', {
    spaceId: props.spaceId,
    fspath: props.fspath
  });

  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let result = iRet.ret;
  state.pluginStats = result.pluginStats;
  loadingMgr.removeLoadingState(stateName);
  for (let pluginStat of state.pluginStats) {
    await workStore.queryPluginDef(pluginStat.id);
  }
  RecheckTab();
}

onMounted(() => {
  refreshFolderDetail();
});
</script>
