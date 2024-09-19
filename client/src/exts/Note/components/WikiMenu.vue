<template>
  <q-toolbar>
    <q-toolbar-title class="mn-toolbar-title">
      我的文件夹
    </q-toolbar-title>
    <q-btn
      round
      dense
      flat
      icon="create_new_folder"
      @click="noteStore.addFolder(RootFolderId, loadingMgr)"
    />
  </q-toolbar>
  <q-separator />
  <q-list>
    <FolderItem
      v-for="subfolder in noteStore.rootFolder.children"
      :key="subfolder.gid"
      :level="0"
      :data="subfolder"
      :loading-mgr="loadingMgr"
    ></FolderItem>
  </q-list>
  <q-inner-loading
    :showing="loadingMgr.isloading.value"
    :label="loadingMgr.loadingLabel.value"
  />
</template>

<script setup>
import { useMikLoading } from 'miknas/exts/Official/shares';
import { useNoteStore, RootFolderId } from '../stores/note';
import FolderItem from './FolderItem.vue';
import { onMounted } from 'vue';

const loadingMgr = useMikLoading();

const noteStore = useNoteStore();

onMounted(()=>{
  noteStore.refreshMyFolders(loadingMgr);
  noteStore.refreshMyNotes(loadingMgr);
})
</script>
