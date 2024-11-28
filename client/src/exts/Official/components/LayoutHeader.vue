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
              <PageMenuItem title="首页" :to="curExtsInfo.indexTo"></PageMenuItem>
            </q-tabs>
          </slot>
        </template>
        <template v-else>
          <slot name="unlogin-toolbar">
            <q-tabs shrink stretch>
              <PageMenuItem title="登录" :href="officialStore.loginUrl"></PageMenuItem>
            </q-tabs>
          </slot>
        </template>
      </slot>
    </q-toolbar>
  </q-header>
</template>

<script setup>
import { reactive } from 'vue';

import { calcValidExtsInfos } from 'miknas/utils';
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useOfficialStore } from '../stores/official.js';
import PageMenuItem from '../components/PageMenuItem.vue';

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

const allExtsInfos = reactive(calcValidExtsInfos());

const curExtsInfo = computed(() => {
  let extsId = curExtsId.value;
  if (!extsId) return {};
  return allExtsInfos[extsId] || {};
});

function toggleLeftDrawer() {
  officialStore.updateLeftDrawerOpen(!officialStore.leftDrawerOpen);
}
</script>
