<template>
  <ExtensionPage>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="首页" :to="extsObj.routePath('')"></PageMenuItem>
        <!-- <q-btn dense flat round icon="vertical_split" @click="noteStore.toggleDrawerOpen" /> -->
      </q-tabs>
    </template>
    <q-page>
      <q-layout container view="hHh Lpr fff" class="absolute-full">
        <ResizeDrawer v-model="noteStore.isDrawerOpen" bordered side="left" :width="260">
          <q-scroll-area class="fit">
            <WikiMenu />
          </q-scroll-area>
        </ResizeDrawer>
        <q-page-container class="bg-blue-grey-1">
          <router-view />
        </q-page-container>
      </q-layout>
    </q-page>
  </ExtensionPage>
</template>

<script setup>
import { useExtension } from '../extMain';
import WikiMenu from '../components/WikiMenu.vue';
import { useNoteStore } from '../stores/note';
import { ResizeDrawer, ExtensionPage, PageMenuItem } from 'miknas/exts/Official/shares';
const extsObj = useExtension();
const noteStore = useNoteStore();
</script>
