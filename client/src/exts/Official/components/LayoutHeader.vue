<template>
  <q-header v-if="showHeader" class="mn-page-header" height-hint="64">
    <q-toolbar class="q-pa-none">
      <q-btn dense flat round icon="space_dashboard" @click="toggleLeftDrawer" />

      <slot name="toolbar">
        <q-toolbar-title class="mn-toolbar-title">
          {{ curExtsInfo.title }}
        </q-toolbar-title>
        <template v-if="officialStore.uid">
          <slot name="login-toolbar">
            <q-tabs shrink stretch>
              <PageMenuItem title="首页" :to="curExtsInfo.index"></PageMenuItem>
            </q-tabs>
          </slot>
        </template>
        <template v-else>
          <slot name="unlogin-toolbar"></slot>
        </template>
      </slot>
    </q-toolbar>
  </q-header>
</template>

<script setup>
import { reactive, ref } from 'vue';

import { getAllExtensions } from 'miknas/utils';
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useOfficialStore } from '../stores/official.js';
import PageMenuItem from '../components/PageMenuItem.vue';
let allExtsObjs = getAllExtensions();

const props = defineProps({
  toolbarNeedLogined: {
    type: Boolean,
    default: false
  }
});

const officialStore = useOfficialStore();

const showHeader = computed(() => {
  if (!props.toolbarNeedLogined) return true;
  return !!officialStore.uid;
});

const route = useRoute();
const curExtsId = computed(() => {
  return route.meta.extsId;
});

function CalcExtsInfos() {
  let ret = {};
  for (let extsId of officialStore.extids) {
    let extsObj = allExtsObjs[extsId];
    if (!extsObj) continue;
    let index = extsObj.getIndex();
    if (!index) continue;
    let info = {
      id: extsObj.id,
      desc: extsObj.desc,
      title: extsObj.title,
      icon: extsObj.icon || 'extension',
      index: index
    };
    ret[extsObj.id] = info;
  }
  return ret;
}

const allExtsInfos = reactive(CalcExtsInfos());

const curExtsInfo = computed(() => {
  let extsId = curExtsId.value;
  if (!extsId) return {};
  return allExtsInfos[extsId] || {};
});

function toggleLeftDrawer() {
  officialStore.updateLeftDrawerOpen(!officialStore.leftDrawerOpen);
}
</script>
