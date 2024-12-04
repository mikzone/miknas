<template>
  <q-table
    flat
    :rows="jobList"
    :columns="JobColumns"
    :filter="state.filterTxt"
    :loading="state.isLoading"
    :rows-per-page-options="[0]"
    class="mn-sticky-header-table"
    row-key="jobId"
    v-bind="$attrs"
  >
    <template #top>
      <div>
        <div class="text-subtitle2">最近作业</div>
        <div class="text-caption text-grey-8">[刷新时间: {{ state.refreshTs }}]</div>
      </div>
      <q-space />
      <q-input
        v-model="state.filterTxt"
        dense
        debounce="200"
        color="primary"
        placeholder="请输入文字进行筛选"
      >
        <template #append>
          <q-icon name="search"></q-icon>
        </template>
      </q-input>
      <q-btn
        v-if="refreshMgr.nextTs.value"
        class="q-ml-sm"
        :disable="state.isLoading"
        :loading="refreshMgr.percentage.value > 0"
        :percentage="refreshMgr.percentage.value"
        dark-percentage
        label="刷新"
      >
        <template #loading>
          <q-spinner-hourglass v-if="!state.isLoading" />
          <q-spinner v-else />
        </template>
        <q-tooltip> 每{{ refreshSec }}秒自动刷新 </q-tooltip>
      </q-btn>
      <q-btn
        v-else
        class="q-ml-sm"
        color="primary"
        :disable="state.isLoading"
        label="刷新"
        @click="refreshMgr.forceRefresh()"
      >
        <template #loading>
          <q-spinner-hourglass v-if="!state.isLoading" />
          <q-spinner v-else />
        </template>
        <q-tooltip> 每{{ refreshSec }}秒自动刷新 </q-tooltip>
      </q-btn>
      <q-btn dense flat stretch round icon="more_vert" @click.stop.prevent="">
        <q-menu auto-close>
          <q-list dense>
            <q-item clickable @click="state.autoRefresh = !state.autoRefresh">
              <q-item-section side>
                <q-checkbox v-model="state.autoRefresh" dense />
              </q-item-section>
              <q-item-section>自动刷新</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </template>
    <template #body-cell-jobId="cellProps">
      <q-td :props="cellProps" @click="showJob(cellProps.value)">
        <q-btn flat size="sm" color="primary" :label="cellProps.value">
          <q-tooltip>查看详情</q-tooltip>
        </q-btn>
      </q-td>
    </template>
    <template #body-cell-runningState="props">
      <q-td :props="props">
        <q-btn v-if="props.row.runningState == 'done'" size="sm" flat color="positive" label="done">
        </q-btn>
        <q-btn
          v-else-if="props.row.runningState == 'canceled' || props.row.runningState == 'errstop'"
          size="sm"
          flat
          color="negative"
          :label="props.row.runningState"
        >
          <q-tooltip> {{ props.row.failtxt }} </q-tooltip>
        </q-btn>
        <q-btn v-else size="sm" color="secondary" :label="props.row.runningState">
          <q-menu>
            <q-list style="min-width: 100px">
              <q-item v-close-popup clickable @click="tryCancel(props.row.jobId, 'kill')">
                <q-item-section>终止任务</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-td>
    </template>
    <template #loading>
      <q-inner-loading showing color="primary" />
    </template>
  </q-table>
</template>

<script>
const JobColumns = [
  {
    name: 'jobId',
    label: 'jobId',
    field: 'jobId',
    sortable: true,
    align: 'left'
  },
  {
    name: 'runningState',
    label: '运行状态',
    field: 'runningState',
    sortable: true,
    align: 'left'
  },
  {
    name: 'title',
    label: '作业名称',
    field: 'title',
    sortable: true,
    align: 'left'
  },
  {
    name: 'uid',
    label: 'uid',
    field: 'uid',
    sortable: true,
    align: 'left'
  },
  {
    name: 'nameSpace',
    label: 'nameSpace',
    field: 'nameSpace',
    sortable: true,
    align: 'left'
  }
];
</script>

<script setup>
import { gutil, MikCall } from 'miknas/utils';
import { computed, onMounted, reactive, watch } from 'vue';
import { coFetchResult } from '../../exec_cmd_util';
import useExtension from '../../extMain';
import { useAutoRefresh } from 'miknas/exts/Official/shares';
let extsObj = useExtension();

const props = defineProps({
  spaceId: {
    type: String,
    default: ''
  },
  pluginId: {
    type: String,
    default: ''
  },
  pluginRootPath: {
    type: String,
    default: ''
  }
});

const state = reactive({
  jobsDict: null,
  refreshTs: '',
  filterTxt: '',
  autoRefresh: false,
  isLoading: true
});

const jobList = computed(() => {
  if (!state.jobsDict) return [];
  return Object.values(state.jobsDict);
});

async function tryRefreshJobsDict() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('queryAllJobs', {
    spaceId: props.spaceId,
    pluginId: props.pluginId,
    pluginRootPath: props.pluginRootPath
  });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let newJobsDict = iRet.ret;
  state.jobsDict = newJobsDict;
  state.refreshTs = gutil.getNowFormatDate();
  state.isLoading = false;
}

const refreshSec = 5;

const refreshMgr = useAutoRefresh(refreshSec * 1000, tryRefreshJobsDict, true);

async function tryCancel(jobId, killType) {
  let isOk = await MikCall.coMakeConfirm(`是否取消执行该任务`);
  if (!isOk) return;
  let iRet = await extsObj.mcpost('reqCancelJob', { jobId: jobId, killType: killType });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let cbRet = iRet.ret;
  MikCall.sendSuccTips(cbRet);
  await refreshMgr.forceRefresh();
}

async function forceRefresh() {
  await refreshMgr.forceRefresh();
}

async function showJob(jobId) {
  await coFetchResult({ jobId: jobId });
  forceRefresh();
}

onMounted(() => {
  refreshMgr.forceRefresh();
});

watch(
  () => state.autoRefresh,
  (newVal) => {
    if (newVal) {
      refreshMgr.start();
    } else {
      refreshMgr.stop();
    }
  }
);

defineExpose({
  forceRefresh
});
</script>
