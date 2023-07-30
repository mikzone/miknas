<template>
  <q-breadcrumbs>
    <q-breadcrumbs-el icon="home" />
    <q-breadcrumbs-el
      v-for="folderInfo in bclist"
      :key="folderInfo.id"
      :label="folderInfo.name"
      icon="folder"
      :to="extsObj.routePath(`folder/${folderInfo.id}`)"
    />
  </q-breadcrumbs>
</template>

<script setup>
import { computed } from 'vue';
import { useNoteStore } from '../stores/note';
import useExtension from '../extMain';

const props = defineProps({
  folderid: {
    type: Number,
    required: true,
  },
});

const noteStore = useNoteStore();
const extsObj = useExtension();

function dfs(folderid, result) {
  if (!folderid) return;
  let folderInfo = noteStore.folderDict[folderid];
  if (!folderInfo) return;
  result.unshift({
    id: folderInfo.id,
    name: folderInfo.name,
  });
  if (folderInfo.parent) dfs(folderInfo.parent, result);
}

const bclist = computed(() => {
  let ret = [];
  dfs(props.folderid, ret);
  return ret;
});
</script>
