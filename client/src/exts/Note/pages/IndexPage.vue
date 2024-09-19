<template>
  <q-page>
    <ContainerLayout>
      <template #toolbar>
        <q-toolbar-title class="mn-toolbar-title">所有笔记</q-toolbar-title>
        <q-input v-model="state.searchText" standout="bg-secondary text-white" dense label="搜索笔记" @keyup.enter="doSearch">
          <template #append>
            <q-icon class="cursor-pointer" name="search" @click="doSearch" />
          </template>
        </q-input>
      </template>
      <note-list ref="noteListInst" class="q-gutter-y-md" />
    </ContainerLayout>
  </q-page>
</template>

<script setup>
import { reactive, ref } from 'vue';
import NoteList from '../components/NoteList.vue';
import ContainerLayout from '../layouts/ContainerLayout.vue';

const state = reactive({
  searchText: '',
});

const noteListInst = ref();
function doSearch() {
  noteListInst.value.doSearch(state.searchText)
}
</script>
