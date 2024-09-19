<template>
  <q-page>
    <q-table
      :rows="state.tableRows"
      :columns="state.tableColumns"
      :filter="state.filterTxt"
      :loading="state.isLoading"
      :rows-per-page-options="[0]"
      class="mn-sticky-header-table mn-sticky-last-column-table absolute-full"
      row-key="mid"
      v-bind="$attrs"
    >
      <template #top>
        <div class="q-table__title">查看所有密文</div>
        <q-space />
        <q-input
          v-model="state.filterTxt"
          dense
          debounce="200"
          color="primary"
          placeholder="请输入文字进行筛选"
          class="gt-xs"
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
          @click="tryRefreshTable()"
        />
        <q-btn
          class="q-ml-sm"
          color="primary"
          label="添加密文"
          @click="tryAddSecret"
        />
      </template>

      <template #body-cell-mid="cellProps">
        <q-td :props="cellProps">
          <q-btn
            flat
            size="sm"
            no-caps
            color="primary"
            :label="cellProps.value"
            :to="viewSecretRoute(cellProps.row)"
          >
            <q-tooltip>查看详情</q-tooltip>
          </q-btn>
        </q-td>
      </template>

      <template #body-cell-op="cellProps">
        <q-td auto-width :props="cellProps">
          <div>
            <q-btn dense flat round icon="more_vert">
              <q-menu auto-close>
                <q-list>
                  <q-item clickable @click="tryShare(cellProps.row)">
                    <q-item-section>分享</q-item-section>
                  </q-item>
                  <q-item clickable @click="tryRemove(cellProps.row)">
                    <q-item-section>删除</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
        </q-td>
      </template>

      <template #loading>
        <q-inner-loading showing color="primary" />
      </template>
    </q-table>
  </q-page>
</template>
<script>
const tableColumns = [
  {
    name: 'mid',
    label: '密文ID',
    field: 'mid',
    sortable: true,
    align: 'left',
    style: 'width: 180px',
  },
  {
    name: 'bts',
    label: '添加时间',
    field: 'viewBts',
    sortable: true,
    align: 'left',
    style: 'width: 140px',
  },
  {
    name: 'ts',
    label: '过期时间',
    field: 'viewTs',
    sortable: true,
    align: 'left',
    style: 'width: 140px',
  },
  {
    name: 'hint',
    label: '提示信息',
    field: 'hint',
    sortable: true,
    align: 'left',
  },
  {
    name: 'op',
    align: 'left',
  },
];
</script>
<script setup>
import { onMounted, reactive } from 'vue';
import { useExtension } from '../extMain.js';
import { gutil, MikCall } from 'miknas/utils';
import { createNewSecret, viewSecretRoute } from '../helper.js';
import { openTextCopyDlg } from 'miknas/exts/Official/shares';

const state = reactive({
  tableRows: [],
  filterTxt: '',
  isLoading: true,
  tableColumns: tableColumns,
});

const extsObj = useExtension();
async function tryRefreshTable() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('querySecrets', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let tableRows = iRet.ret;
  for (let info of tableRows) {
    info.viewTs = gutil.formatTs(info.ts * 1000);
    info.viewBts = gutil.formatTs(info.bts * 1000);
  }
  tableRows.sort((a, b) => a.bts - b.bts);
  state.tableRows = tableRows;
  state.isLoading = false;
}

async function tryAddSecret() {
  let result = await createNewSecret();
  if (result) {
    tryRefreshTable();
  }
}

function tryShare(info) {
  openTextCopyDlg({ title: '分享链接如下:', txt: viewSecretRoute(info, true) });
}

async function tryRemove(info) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除?`);
  if (!isOk) return;
  state.isLoading = true;

  let iRet = await extsObj.mcpost('removeSecret', { mid: info.mid });
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
