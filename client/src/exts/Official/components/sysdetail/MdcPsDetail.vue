<template>
  <q-table
    :rows="state.processList"
    :columns="state.processColumns"
    :filter="state.filterTxt"
    :loading="state.isLoading"
    :rows-per-page-options="[0]"
    class="mn-sticky-header-table"
    no-data-label="没有匹配的进程信息"
    row-key="pid"
    v-bind="$attrs"
  >
    <template #top>
      <div class="q-table__title">查看进程({{ state.refreshTs }})</div>
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

    <template #loading>
      <q-inner-loading showing color="primary" />
    </template>
  </q-table>
</template>
<script>
const PROCESS_COLUMNS = [
  {
    name: 'pid',
    label: 'PID',
    field: 'pid',
    sortable: true,
    align: 'left',
  },
  {
    name: 'ppid',
    label: 'PPID',
    field: 'ppid',
    sortable: true,
    align: 'left',
  },
  {
    name: 'gids',
    label: 'gids',
    field: 'gids',
    sortable: true,
    align: 'left',
  },
  {
    name: 'username',
    label: 'USER',
    field: 'username',
    sortable: true,
    align: 'left',
  },
  {
    name: 'cpu_percent',
    label: '%CPU',
    field: 'cpu_percent',
    sortable: true,
    align: 'left',
  },
  {
    name: 'memory_percent',
    label: '%MEM',
    field: 'memory_percent',
    sortable: true,
    align: 'left',
  },
  {
    name: 'pass_time',
    label: 'pass_time',
    field: 'pass_time',
    sortable: true,
    align: 'left',
  },
  {
    name: 'cmdline',
    label: 'COMMAND',
    field: 'cmdline',
    sortable: true,
    align: 'left',
  },
];
</script>

<script setup>
import { gutil, MikCall } from 'miknas/utils';
import { onMounted, reactive } from 'vue';
import { useExtension } from '../../extMain';
const props = defineProps({
  grepStr: {
    type: String,
    default: '',
  },
});

const state = reactive({
  processList: [],
  refreshTs: '',
  filterTxt: '',
  isLoading: true,
  processColumns: PROCESS_COLUMNS,
});

const extsObj = useExtension();
function tryRefreshJobsDict() {
  let filterStr = props.grepStr || '';
  extsObj
    .mpost('query_all_process', { grepStr: filterStr })
    .then((processList) => {
      state.processList = processList;
      state.refreshTs = gutil.getNowFormatDate();
    })
    .catch(MikCall.alertRespErrMsg)
    .finally(() => {
      state.isLoading = false;
    });
  state.isLoading = true;
}

onMounted(() => {
  tryRefreshJobsDict();
});
</script>
