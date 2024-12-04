<template>
  <q-card v-if="state.canWalkDir" class="q-mb-md">
    <MdcDriveAliveView
      :fsid="`Cw_${props.spaceId}`"
      :fsrela="props.fspath"
      :kind="props.kind"
      :extra-conf="myExtraConf"
    ></MdcDriveAliveView>
  </q-card>
  <template v-if="props.kind == 'list'">
    <q-card v-if="!state.canWalkDir" class="q-mb-md">
      <q-card-section class="row items-center">
        <q-breadcrumbs gutter="xs" class="text-wrap">
          <q-breadcrumbs-el
            v-for="oneBC in fileBC"
            :key="oneBC.name"
            :label="oneBC.name"
            :icon="oneBC.icon"
            :class="oneBC.path === undefined ? '' : 'text-primary cursor-pointer'"
            @click="gotoPath(oneBC.path)"
          />
        </q-breadcrumbs>
      </q-card-section>
      <q-list separator bordered>
        <q-item
          v-for="subDirInfo in state.subdirs"
          :key="subDirInfo.dirName"
          clickable
          @click="gotoSubDir(subDirInfo.dirName)"
        >
          <q-item-section avatar>
            <q-icon name="folder" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ subDirInfo.dirName }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
      <q-card-section v-if="state.subdirs.length <= 0" class="text-caption text-grey">
        无相关子目录
      </q-card-section>
    </q-card>
    <q-card class="q-mb-md">
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
      />
      <q-inner-loading
        :showing="loadingMgr.isloading.value"
        :label="loadingMgr.loadingLabel.value"
      />
    </q-card>
    <q-card class="q-mb-md"> </q-card>
  </template>
</template>
<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { MdcDriveAliveView } from 'miknas/exts/Drive/shares';
import { useExtension } from '../extMain';
import { MikCall } from 'miknas/utils';
import { useWorkStore } from '../stores/work';
import PluginView from './PluginView.vue';
import { useRouter } from 'vue-router';
import { FileUtil } from 'miknas/exts/Drive/shares';

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
  },
  kind: {
    type: String,
    default: 'list'
  }
});

const state = reactive({
  pluginStats: [],
  canWalkDir: false,
  subdirs: [],
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

const curSpaceInfo = computed(() => {
  return workStore.spaces[props.spaceId] || {};
});

function RecheckTab() {
  if (state.tab) return;
  for (let pluginDetail of relPlugins.value) {
    let pluginId = pluginDetail.pluginId;
    state.tab = pluginId;
  }
}

const router = useRouter();

function gotoPath(newFspath) {
  let newloc = { params: { routeSubPath: newFspath } };
  router.push(newloc);
}

function gotoSubDir(dirName) {
  let newFspath;
  if (props.fspath === '') newFspath = dirName;
  else newFspath = props.fspath + '/' + dirName;
  return gotoPath(newFspath);
}

const fileBC = computed(() => {
  let ret = [];
  let curPath = props.fspath;
  ret.push({
    icon: 'home',
    name: curSpaceInfo.value.Path,
    path: ''
  });
  let pathSep = FileUtil.getPathSep();
  let subFodlers = curPath.split(pathSep);
  let genPath = '';
  for (let folderName of subFodlers) {
    if (folderName) {
      genPath = FileUtil.contactFolderName(genPath, folderName);
      ret.push({
        name: folderName,
        path: genPath
      });
    }
  }
  // 最后一个不能点击
  delete ret[ret.length - 1]['path'];
  return ret;
});

async function refreshFolderDetail() {
  if (props.kind !== 'list') return;
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
  state.canWalkDir = result.canWalkDir || false;
  state.subdirs = result.subdirs || [];
  loadingMgr.removeLoadingState(stateName);
  for (let pluginStat of state.pluginStats) {
    await workStore.queryPluginDef(pluginStat.id);
  }
  RecheckTab();
}

const myExtraConf = {
  // viewFn: null
};

onMounted(() => {
  refreshFolderDetail();
});
</script>
