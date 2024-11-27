<template>
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
</template>

<script setup>
import WikiMenu from '../components/WikiMenu.vue';
import { useNoteStore } from '../stores/note';
import { ResizeDrawer } from 'miknas/exts/Official/shares';
const noteStore = useNoteStore();
</script>
