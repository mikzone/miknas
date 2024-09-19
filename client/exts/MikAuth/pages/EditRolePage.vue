<template>
  <q-page>
    <q-table
      v-bind="$attrs"
      v-model:selected="state.selected"
      :rows="state.tableRows"
      :columns="state.tableColumns"
      :filter="state.filterTxt"
      :loading="state.isLoading"
      :rows-per-page-options="[0]"
      class="mn-sticky-header-table absolute-full"
      row-key="resource"
      selection="multiple"
    >
      <template #top>
        <div class="q-table__title">修改{{ props.roleid }}的权限</div>
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
          class="q-ml-sm gt-xs"
          color="primary"
          :disable="state.isLoading"
          label="刷新"
          @click="tryRefreshTable()"
        />
      </template>
      <template #bottom>
        <q-btn
          class="q-ml-sm"
          color="primary"
          :disable="state.isLoading"
          label="保存"
          @click="trySave()"
        />
      </template>
      <template #body-cell-resource="cellProps">
        <q-td :props="cellProps">
          <div class="text-subtitle2 text-primary">{{ cellProps.value }}</div>
          <div class="text-caption">{{ cellProps.row.desc }}</div>
        </q-td>
      </template>
      <template #body-selection="scope">
        <q-toggle
          v-model="scope.selected"
          checked-icon="check"
          :color="scope.selected ? 'green' : 'red'"
          :label="scope.selected ? '允许' : '禁止'"
          :class="scope.row.defv ? 'white-auth' : 'black-auth'"
          :keep-color="scope.row.defv"
          unchecked-icon="clear"
        />
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
    name: 'resource',
    label: '权限列表',
    field: 'resource',
    sortable: true,
    align: 'left',
  },
  {
    name: 'type',
    label: '权限类型',
    field: 'type',
    sortable: true,
    align: 'left',
  },
];
</script>
<script setup>
import { onMounted, reactive } from 'vue';
import { useExtension } from '../extMain.js';
import { gutil, MikCall } from 'miknas/utils';

const props = defineProps({
  roleid: {
    type: String,
    required: true,
  },
});

const state = reactive({
  tableRows: [],
  filterTxt: '',
  isLoading: true,
  roleids: [],
  tableColumns: tableColumns,
  selected: [],
});

const extsObj = useExtension();
async function tryRefreshTable() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('oneRole', { role: props.roleid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let data = iRet.ret;
  state.tableRows = Object.values(data.auths);
  let selected = [];
  for (let resInfo of state.tableRows) {
    let resId = resInfo.resource;
    resInfo.defv = gutil.authCheck(resId, {});
    resInfo.type = resInfo.defv ? '默认允许' : '默认禁止';
    let flag = gutil.authCheck(resId, data.roleInfo.cans);
    if (flag) {
      selected.push(data.auths[resId]);
    }
  }
  state.selected = selected;
  state.isLoading = false;
}

async function trySave() {
  state.isLoading = true;
  let cans = {};
  for (let selectedInfo of state.selected) {
    cans[selectedInfo.resource] = true;
  }
  for (let resInfo of state.tableRows) {
    if (!cans[resInfo.resource]) cans[resInfo.resource] = false;
  }

  let iRet = await extsObj.mcpost('saveRole', {
    role: props.roleid,
    cans: cans,
  });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  MikCall.sendSuccTips('修改成功');
  await tryRefreshTable();
  state.isLoading = false;
}

onMounted(() => {
  tryRefreshTable();
});
</script>

<style lang="scss">
.white-auth {
  .q-toggle__inner--truthy {
    .q-toggle__thumb .q-icon {
      opacity: 0.54;
      color: #000;
    }

    .q-toggle__thumb:after {
      background-color: #fff;
    }

    .q-toggle__thumb .q-icon {
      color: #4caf50;
      opacity: 1;
    }

    .q-toggle__track {
      color: #000;
      opacity: 0.38;
    }
    .q-toggle__thumb {
      left: 0.45em;
    }
  }

  .q-toggle__inner--falsy {
    .q-toggle__thumb:after {
      background-color: currentColor;
    }
    .q-toggle__thumb .q-icon {
      opacity: 1;
      color: #fff;
    }
  }
}

.black-auth .q-toggle__inner--falsy {
  .q-toggle__thumb {
    left: 0.45em;
  }
  .q-toggle__thumb .q-icon {
    color: #f44336;
    opacity: 1;
  }
}
</style>
