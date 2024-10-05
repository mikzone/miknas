<template>
  <q-page>
    <q-table
      :rows="state.tableRows"
      :columns="state.tableColumns"
      :filter="state.filterTxt"
      :loading="state.isLoading"
      :rows-per-page-options="[0]"
      class="mn-sticky-header-table absolute-full"
      row-key="role"
      v-bind="$attrs"
    >
      <template #top>
        <div class="q-table__title">所有角色权限</div>
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
          label="添加角色"
          @click="tryCreateRole()"
        />
        <q-btn
          class="q-ml-sm gt-xs"
          color="primary"
          :disable="state.isLoading"
          label="刷新"
          @click="tryRefreshTable()"
        />
      </template>

      <template #body-cell-cans="cellProps">
        <q-td :props="cellProps" class="q-gutter-xs">
          <template
            v-for="(flag, resource) in cellProps.value"
            :key="resource"
          >
            <q-chip
              v-if="flag"
              color="teal"
              text-color="white"
              icon="add"
              size="sm"
            >
              {{ resource }}
            </q-chip>
            <q-chip
              v-else
              color="red"
              text-color="white"
              icon="remove"
              size="sm"
            >
              {{ resource }}
            </q-chip>
          </template>
        </q-td>
      </template>

      <template #body-cell-op="cellProps">
        <q-td auto-width :props="cellProps">
          <div>
            <q-btn dense flat round icon="more_vert">
              <q-menu>
                <q-list>
                  <q-item
                    clickable
                    :to="extsObj.routePath(`edit_role/${cellProps.row.role}`)"
                  >
                    <q-item-section>修改权限</q-item-section>
                  </q-item>
                  <q-item clickable @click="tryRemove(cellProps.row)">
                    <q-item-section>删除角色</q-item-section>
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
    name: 'role',
    label: '角色',
    field: 'role',
    sortable: true,
    align: 'left',
  },
  {
    name: 'cans',
    label: '修改的权限列表',
    field: 'cans',
    align: 'left',
    style: 'white-space: break-spaces',
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
import { MikCall } from 'miknas/utils';

const state = reactive({
  tableRows: [],
  filterTxt: '',
  isLoading: true,
  roleids: [],
  tableColumns: tableColumns,
});

const extsObj = useExtension();
async function tryRefreshTable() {
  state.isLoading = true;
  let iRet = await extsObj.mcpost('allRoles', {});
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  let data = iRet.ret;
  state.tableRows = data.roles;
  state.isLoading = false;
}

async function tryRemove(info) {
  let isOk = await MikCall.coMakeConfirm(`是否确认删除角色(${info.role})?`);
  if (!isOk) return;
  state.isLoading = true;
  let iRet = await extsObj.mcpost('removeRole', { role: info.role });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  MikCall.sendSuccTips('删除成功');
  await tryRefreshTable();
  state.isLoading = false;
}

async function tryCreateRole() {
  let [isOk, newRoleId] = await MikCall.coMakePrompt(`请输入角色名称`, '');
  if (!isOk) return;
  if (!newRoleId) return;
  state.isLoading = true;

  let iRet = await extsObj.mcpost('addRole', { role: newRoleId });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.isLoading = false;
    return;
  }
  MikCall.sendSuccTips('添加成功');
  await tryRefreshTable();
  state.isLoading = false;
}

onMounted(() => {
  tryRefreshTable();
});
</script>
