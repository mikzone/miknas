<template>
  <q-table
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
      <div class="q-table__title">查看进行中的任务({{ state.refreshTs }})</div>
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
        class="q-ml-sm"
        color="primary"
        :disable="state.isLoading"
        label="刷新"
        @click="tryRefreshJobsDict()"
      />
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
        <q-btn
          v-if="props.row.runningState == 'done'"
          size="sm"
          flat
          color="positive"
          label="done"
        >
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
        <q-btn
          v-else
          size="sm"
          color="secondary"
          :label="props.row.runningState"
        >
          <q-menu>
            <q-list style="min-width: 100px">
              <q-item
                v-close-popup
                clickable
                @click="tryCancel(props.row.jobId, 'kill')"
              >
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
    align: 'left',
  },
  {
    name: 'uid',
    label: 'uid',
    field: 'uid',
    sortable: true,
    align: 'left',
  },
  {
    name: 'runningState',
    label: '运行状态',
    field: 'runningState',
    sortable: true,
    align: 'left',
  },
  {
    name: 'cmd',
    label: 'CMD',
    field: 'cmd',
    sortable: true,
    align: 'left',
  },
  {
    name: 'nameSpace',
    label: 'nameSpace',
    field: 'nameSpace',
    sortable: true,
    align: 'left',
  },
];
</script>

<script setup>
import { gutil, MikCall } from 'miknas/utils';
import { computed, onMounted, reactive } from 'vue';
import { fetchResult } from '../exec_cmd_util';
import { useExtension } from '../extMain';
let extsObj = useExtension();

const state = reactive({
  jobsDict: null,
  refreshTs: '',
  filterTxt: '',
  isLoading: true,
});

const jobList = computed(() => {
  if (!state.jobsDict) return [];
  return Object.values(state.jobsDict);
});

async function tryRefreshJobsDict() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('queryAllJobs', {});
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
  await tryRefreshJobsDict();
}

function showJob(jobId) {
  fetchResult({ jobId: jobId })
}

onMounted(() => {
  tryRefreshJobsDict();
});
</script>
