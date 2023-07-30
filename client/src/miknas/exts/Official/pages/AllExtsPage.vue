<template>
  <q-page>
    <q-list bordered>
      <q-item-label header>{{ label }}</q-item-label>
      <template v-for="(extsInfo, idx) in allExtsInfos" :key="extsInfo.id">
        <q-separator v-if="idx > 0" />
        <q-item v-ripple :to="extsInfo.to" active-class="">
          <q-item-section avatar>
            <q-icon :name="extsInfo.icon" />
          </q-item-section>
          <q-item-section>
            <q-item-label lines="1"> <q-badge color="secondary" :label="extsInfo.id" /> </q-item-label>
            <q-item-label lines="1"> {{ extsInfo.title }}</q-item-label>
            <q-item-label caption> {{ extsInfo.desc }}</q-item-label>
          </q-item-section>
          <q-item-section v-if="Boolean(extsInfo.to)" side>
            <q-icon name="arrow_forward" />
          </q-item-section>
        </q-item>
      </template>
    </q-list>
  </q-page>
</template>

<script setup>
import { reactive, computed } from 'vue';
import { getAllExtensions } from 'miknas/utils';
import { useOfficialStore } from 'miknas/exts/Official/stores/official';
import { useExtension } from '../extMain';

const officialStore = useOfficialStore();
const officialExtsObj = useExtension();

function CalcExtsInfos() {
  //
  let ret = [];
  let allExtsObjs = getAllExtensions();
  for (let extsId of officialStore.extids) {
    if (extsId == officialExtsObj.id) continue;
    let extsObj = allExtsObjs[extsId];
    if (!extsObj) continue;
    let info = {
      id: extsObj.id,
      desc: extsObj.desc,
      title: extsObj.title,
      icon: extsObj.icon || 'extension',
      to: extsObj.getIndex(),
    };
    ret.push(info);
  }
  return ret;
}

const allExtsInfos = reactive(CalcExtsInfos());

const label = computed(() => {
  if (allExtsInfos.length > 0) return '所有扩展';
  return '当前角色没有可用的扩展，请联系系统管理员添加';
})
</script>
