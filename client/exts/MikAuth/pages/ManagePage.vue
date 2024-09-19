<template>
  <q-page>
    <q-table
      :rows="state.tableRows"
      :columns="state.tableColumns"
      :filter="state.filterTxt"
      :loading="state.isLoading"
      :rows-per-page-options="[0]"
      class="mn-sticky-header-table mn-sticky-last-column-table absolute-full"
      row-key="pid"
      v-bind="$attrs"
    >
      <template #top>
        <div class="q-table__title">用户管理</div>
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
      </template>

      <template #body-cell-mid="cellProps">
        <q-td :props="cellProps">
          <q-btn
            flat
            size="sm"
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
              <q-menu>
                <q-list>
                  <q-item clickable>
                    <q-item-section>设置权限为</q-item-section>
                    <q-item-section side>
                      <q-icon name="keyboard_arrow_right" />
                    </q-item-section>
                    <q-menu anchor="top end" self="top start">
                      <q-list>
                        <template v-for="tmprole in state.roles">
                          <q-item
                            v-if="tmprole != cellProps.row.role"
                            :key="tmprole"
                            clickable
                            @click="trySetRole(cellProps.row, tmprole)"
                          >
                            <q-item-section>{{ tmprole }}</q-item-section>
                          </q-item>
                        </template>
                      </q-list>
                    </q-menu>
                  </q-item>
                  <q-item clickable @click="tryRemove(cellProps.row)">
                    <q-item-section>删除用户</q-item-section>
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
    name: 'uid',
    label: 'uid',
    field: 'uid',
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
    name: 'role',
    label: 'role',
    field: 'role',
    sortable: true,
    align: 'left',
  },
  {
    name: 'cts',
    label: '注册时间',
    field: 'viewCts',
    sortable: true,
    align: 'left',
    style: 'width: 140px',
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

const state = reactive({
  tableRows: [],
  filterTxt: '',
  isLoading: true,
  roles: [],
  tableColumns: tableColumns,
});

const extsObj = useExtension();
async function tryRefreshTable() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('queryAllUser', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let { users, roles } = iRet.ret;
  state.roles = roles;
  for (let info of users) {
    info.viewCts = gutil.formatTs(info.cts * 1000);
  }
  users.sort((a, b) => a.bts - b.bts);
  state.tableRows = users;
  state.isLoading = false;
}

async function trySetRole(info, newRole) {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('modifyUserRole', {
    uid: info.uid,
    role: newRole,
  });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  MikCall.sendSuccTips('设置成功');
  state.isLoading = false;
  // 保险起见，重新刷新一次，应该用户也不多
  tryRefreshTable();
}

async function tryRemove(info) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除用户(${info.uid})?`);
  if (!isOk) return;

  state.isLoading = true;

  let iRet = await extsObj.mcpost('removeUser', { uid: info.uid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  MikCall.sendSuccTips('删除成功');
  state.isLoading = false;
  await tryRefreshTable();
}

onMounted(() => {
  tryRefreshTable();
});
</script>
