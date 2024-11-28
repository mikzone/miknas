<template>
  <q-page>
    <q-list bordered>
      <q-item-label header>{{ label }}</q-item-label>
      <template v-for="(extsInfo, idx) in allExtsInfoList" :key="extsInfo.id">
        <q-separator v-if="idx > 0" />
        <q-item v-ripple :to="extsInfo.indexTo" active-class="">
          <q-item-section avatar>
            <q-icon :name="extsInfo.icon" />
          </q-item-section>
          <q-item-section>
            <q-item-label lines="1">
              <q-badge color="secondary" :label="extsInfo.id" />
            </q-item-label>
            <q-item-label lines="1"> {{ extsInfo.title }}</q-item-label>
            <q-item-label caption> {{ extsInfo.desc }}</q-item-label>
          </q-item-section>
          <q-item-section v-if="Boolean(extsInfo.indexTo)" side>
            <q-icon name="arrow_forward" />
          </q-item-section>
        </q-item>
      </template>
    </q-list>
  </q-page>
</template>

<script setup>
import { reactive, computed } from 'vue';
import { useExtension } from '../extMain';
import { calcValidExtsInfos } from 'miknas/utils';

const officialExtsObj = useExtension();

const allExtsInfoList = reactive(Object.values(calcValidExtsInfos([officialExtsObj.id])));

const label = computed(() => {
  if (allExtsInfoList.length > 0) return '所有扩展';
  return '当前角色没有可用的扩展，请联系系统管理员添加';
});
</script>
