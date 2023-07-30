<template>
  <q-table
    :rows="state.threadList"
    :columns="state.threadColumns"
    :filter="state.filterTxt"
    :loading="state.isLoading"
    :rows-per-page-options="[0]"
    class="mn-sticky-header-table"
    row-key="pid"
    v-bind="$attrs"
  >
    <template #top>
      <div class="q-table__title">查看当前所有线程({{ state.refreshTs }})</div>
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
const THREAD_COLUMNS = [
  {
    name: 'ident',
    label: 'ident',
    field: 'ident',
    sortable: true,
    align: 'left',
  },
  {
    name: 'native_id',
    label: 'native_id',
    field: 'native_id',
    sortable: true,
    align: 'left',
  },
  {
    name: 'name',
    label: 'name',
    field: 'name',
    sortable: true,
    align: 'left',
  },
  {
    name: 'daemon',
    label: 'daemon',
    field: 'daemon',
    sortable: true,
    align: 'left',
  },
  {
    name: 'is_alive',
    label: 'is_alive',
    field: 'is_alive',
    sortable: true,
    align: 'left',
  },
];
</script>
<script setup>
import { gutil, MikCall } from 'miknas/utils';
import { onMounted, reactive } from 'vue';
import { useExtension } from '../../extMain';

const state = reactive({
  threadList: [],
  refreshTs: '',
  filterTxt: '',
  isLoading: true,
  threadColumns: THREAD_COLUMNS,
});

const extsObj = useExtension();
function tryRefreshJobsDict() {
  extsObj
    .mpost('query_all_threads', {})
    .then((threadList) => {
      state.threadList = threadList;
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
