<template>
  <q-page>
    <ContainerLayout>
      <template #toolbar>
        <q-space />
        <q-input
          v-if="isDetail"
          v-model="state.searchText"
          dense
          standout="bg-secondary text-white"
          label="在当前文件夹下搜索"
          class="q-mr-md"
          @keyup.enter="doSearch"
        >
          <template #append>
            <q-icon class="cursor-pointer" name="search" @click="doSearch" />
          </template>
        </q-input>
        <q-tabs v-model="state.viewMode" shrink>
          <q-tab name="list" label="目录" />
          <q-tab name="preview" label="笔记" />
        </q-tabs>
      </template>
      <FolderBreadcrumbs :folderid="folderid" class="q-mb-md" active-color="link" />
      <q-card v-show="!isDetail" class="no-shadow mn-bordered">
        <q-list>
          <q-item-label v-if="childNum" header>子目录和笔记</q-item-label>
          <q-item-label v-else header>暂无目录或笔记</q-item-label>
          <template v-for="data in children" :key="data.gid">
            <q-item
              v-if="data.children"
              v-ripple
              clickable
              :to="extsObj.routePath(`folder/${data.id}`)"
            >
              <q-item-section avatar class="wiki-expand-icon">
                <q-icon name="folder" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-link">{{
                  data.name
                }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item
              v-else
              v-ripple
              clickable
              :to="extsObj.routePath(`view/${data.id}`)"
            >
              <q-item-section avatar class="wiki-expand-icon">
                <q-icon name="article" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-link">{{
                  data.title
                }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-btn dense round flat icon="more_vert" @click.stop.prevent="">
                  <q-menu>
                    <q-list dense>
                      <q-item v-close-popup clickable @click="moveNote(data)">
                        <q-item-section>
                          <q-item-label>移动到</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item
                        v-close-popup
                        clickable
                        @click="noteStore.deleteNote(data)"
                      >
                        <q-item-section>
                          <q-item-label>删除</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </q-item-section>
            </q-item>
          </template>
        </q-list>
      </q-card>

      <NoteList
        v-if="isDetail"
        ref="noteListInst"
        :folderid="folderid"
        class="q-mt-md q-gutter-y-md"
      />
    </ContainerLayout>
  </q-page>
</template>

<script setup>
import { computed, reactive, ref } from 'vue';
import ContainerLayout from '../layouts/ContainerLayout.vue';
import FolderBreadcrumbs from '../components/FolderBreadcrumbs.vue';
import { useNoteStore } from '../stores/note';
import useExtension from '../extMain';
import NoteList from '../components/NoteList.vue';
import { moveNote } from '../helpers';

const props = defineProps({
  folderStrid: {
    type: String,
    required: true,
  },
});

const state = reactive({
  viewMode: 'list',
});

const folderid = computed(() => {
  return parseInt(props.folderStrid);
});

const isDetail = computed(() => {
  return state.viewMode == 'preview';
});

const extsObj = useExtension();
const noteStore = useNoteStore();

const folderInfo = computed(() => {
  return noteStore.folderDict[folderid.value] || {};
});

const children = computed(() => {
  let info = folderInfo.value;
  if (!info || !info.children) return {};
  return info.children;
});

const childNum = computed(() => {
  return Object.keys(children.value).length;
});

const noteListInst = ref();

function doSearch() {
  state.viewMode = 'preview';
  noteListInst.value.doSearch(state.searchText);
}
</script>
