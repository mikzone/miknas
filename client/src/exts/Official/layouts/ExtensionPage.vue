<template>
  <q-layout view="hHh Lpr fff" class="bg-white">
    <q-header v-if="showHeader" class="mn-page-header" height-hint="64">
      <q-toolbar class="q-pa-none">
        <q-btn
          dense
          flat
          round
          icon="space_dashboard"
          @click="toggleLeftDrawer"
        />

        <slot name="toolbar">
          <q-toolbar-title class="mn-toolbar-title">
            {{ curExtsInfo.title }}
          </q-toolbar-title>
          <template v-if="officialStore.uid">
            <slot name="login-toolbar">
              <q-tabs shrink stretch>
                <PageMenuItem
                  title="首页"
                  :to="curExtsInfo.index"
                ></PageMenuItem>
              </q-tabs>
            </slot>
          </template>
          <template v-else>
            <slot name="unlogin-toolbar"></slot>
          </template>
        </slot>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      side="left"
      elevated
      behavior="mobile"
      :width="260"
    >
      <q-scroll-area class="fit">
        <q-list>
          <q-item-label header class="bg-teal text-white"
            >欢迎使用 MikNas</q-item-label
          >
          <q-item v-if="!officialStore.uid" class="bg-teal text-white q-pb-lg">
            <q-item-section avatar>
              <q-icon name="person_off" />
            </q-item-section>
            <q-item-section> 未登录 </q-item-section>
            <q-item-section side>
              <q-btn
                flat
                round
                color="white"
                icon="login"
                :href="officialStore.loginUrl"
              ></q-btn>
            </q-item-section>
          </q-item>
          <q-item v-else class="bg-teal text-white q-pb-lg">
            <q-item-section avatar>
              <q-avatar>
                <!-- <img src="https://cdn.quasar.dev/img/avatar.png"/> -->
                <img
                  src="data:image/png;base64,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"
                />
              </q-avatar>
            </q-item-section>
            <q-item-section> {{ officialStore.uid }} </q-item-section>
            <q-item-section side>
              <q-btn flat round color="white" icon="more_vert">
                <q-menu auto-close>
                  <q-list dense>
                    <q-item clickable :href="officialStore.logoutUrl">
                      <q-item-section avatar>
                        <q-icon name="logout" />
                      </q-item-section>
                      <q-item-section>退出登录</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-btn>
            </q-item-section>
          </q-item>

          <template v-if="officialStore.uid">
            <q-item-label header>所有扩展</q-item-label>
            <q-item
              v-for="extsInfo in allExtsInfos"
              :key="extsInfo.id"
              v-ripple
              clickable
              :to="extsInfo.index"
            >
              <q-item-section avatar>
                <q-icon :name="extsInfo.icon" />
              </q-item-section>
              <q-item-section> {{ extsInfo.title }} </q-item-section>
            </q-item>
          </template>
        </q-list>
      </q-scroll-area>
    </q-drawer>
    <q-page-container>
      <slot>
        <router-view />
      </slot>
    </q-page-container>

    <slot name="layout"></slot>
  </q-layout>
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
    default: false,
  },
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
      index: index,
    };
    ret[extsObj.id] = info;
  }
  return ret;
}

const allExtsInfos = reactive(CalcExtsInfos());

const curExtsInfo = computed(()=>{
  let extsId = curExtsId.value;
  if (!extsId) return {};
  return allExtsInfos[extsId] || {};
});

const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>
