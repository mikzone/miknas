<template>
  <q-page>
    <q-table
      :title="state.curTable && state.curTable.name"
      :rows="state.curRows"
      :columns="state.curColumns"
      class="mn-sticky-header-table absolute-full"
      row-key="_rowKey"
      :rows-per-page-options="[50, 0]"
    >
      <template #top>
        <q-select
          dense
          option-label="name"
          :options="state.tableOptions"
          label="选择表名"
          :model-value="state.curTable"
          class="col"
          @update:model-value="onSelectTable"
        />
        <q-btn v-if="state.curTable" round flat color="primary" icon="info">
          <q-popup-proxy>
            <q-banner dark>
              {{ state.curTable.sql }}
            </q-banner>
          </q-popup-proxy>
        </q-btn>
      </template>
      <template #body-cell-mddbop="cellProps">
        <q-td :props="cellProps">
          <!-- <q-btn
            color="primary"
            flat
            class="q-pa-xs"
            icon="info"
            @click="showRow(cellProps.row)"
          ></q-btn> -->
          <q-btn dense round flat icon="more_vert" @click.stop.prevent="">
            <q-menu>
              <q-list dense>
                <q-item
                  v-for="(rVal, rKey) in cellProps.row"
                  :key="rKey"
                  v-close-popup
                  clickable
                  @click="showRow(rVal)"
                >
                  <q-item-section>
                    <q-item-label>查看{{ rKey }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { onMounted, reactive } from 'vue';
import { useExtension } from '../extMain';
import { MikCall, gutil } from 'miknas/utils';
import { openTextCopyDlg } from 'miknas/exts/Official/helpers';
const extsObj = useExtension();

const state = reactive({
  tableOptions: [],
  tableDatas: {},
  curRows: [],
  curColumns: [],
  curTable: undefined,
});

async function fetchTableDatas(tableName) {
  // 根据表名获取其中的数据
  let iRet = await extsObj.mcget(`viewTable/${tableName}`);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let data = iRet.ret;
  let idx = 0;
  for (let info of data) {
    idx += 1;
    // 补充一个rowkey给显示用
    info._rowKey = `${tableName}_${idx}`;
  }
  return data;
}

async function fetchTableStructs(tableName) {
  // 根据表名获取其中的数据
  let iRet = await extsObj.mcget(`descTable/${tableName}`);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let data = iRet.ret;
  let idx = 0;
  for (let info of data) {
    idx += 1;
    // 补充一个rowkey给显示用
    info._rowKey = `struct_${tableName}_${idx}`;
  }
  return data;
}

async function tryRefreshTables() {
  let data = await fetchTableDatas('sqlite_master');
  if (data) {
    let ret = [];
    for (let info of data) {
      if (info.type == 'table') {
        ret.push({
          name: info.name,
          sql: info.sql,
        });
      }
    }
    state.tableOptions = ret;
  }
}

function showRow(row) {
  let txt;
  if (typeof(row) == 'string') txt = row;
  else txt = JSON.stringify(row, undefined, 2);
  if (!txt) return;
  openTextCopyDlg({ txt });
}

function calcTableColumn(tableStructs) {
  let ret = [];
  for (let colInfo of tableStructs) {
    let colName = colInfo.name;
    let info = {
      name: colName,
      label: colName,
      field: colName,
      sortable: true,
      classes: 'md-db-td',
      align: 'left',
    };
    if (colInfo.type == 'datetime') {
      info.format = (val) => {
        let d = new Date(val);
        return gutil.formatTs(d.getTime());
      };
    }
    ret.push(info);
  }
  ret.push({
    name: 'mddbop',
    label: '',
    sortable: false,
    align: 'left',
  });
  return ret;
}

async function onSelectTable(tableInfo) {
  if (!tableInfo) return;
  let tableName = tableInfo.name;
  if (!tableName) return;
  state.curTable = tableInfo;
  if (!state.tableDatas[tableName]) {
    let tableRows = await fetchTableDatas(tableName);
    if (!tableRows) return;
    let tableStructs = await fetchTableStructs(tableName);
    if (!tableStructs) return;
    let tableColumn = calcTableColumn(tableStructs);
    state.tableDatas[tableName] = {
      Rows: tableRows,
      Structs: tableStructs,
      Cols: tableColumn,
    };
  }
  state.curRows = state.tableDatas[tableName].Rows;
  state.curColumns = state.tableDatas[tableName].Cols;
}

onMounted(() => {
  tryRefreshTables();
});
</script>

<style>
.md-db-td {
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
