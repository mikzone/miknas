<template>
  <q-expansion-item
    v-if="data.children"
    default-opened
    switch-toggle-side
    expand-icon="keyboard_arrow_right"
    expanded-icon="keyboard_arrow_down"
    expand-icon-class="mn-item-icon"
    :header-inset-level="0.5 * level"
    :to="extsObj.routePath(`folder/${data.id}`)"
  >
    <template #header>
      <q-item-section>
        <q-item-label :lines="2">{{ data.name }}</q-item-label>
      </q-item-section>

      <q-item-section side>
        <q-btn dense round flat icon="more_vert" @click.stop.prevent="">
          <q-menu>
            <q-list dense>
              <q-item clickable>
                <q-item-section>新建</q-item-section>
                <q-item-section side>
                  <q-icon name="keyboard_arrow_right" />
                </q-item-section>
                <q-menu anchor="top end" self="top start">
                  <q-list dense>
                    <q-item
                      v-close-popup
                      clickable
                      @click="noteStore.addFolder(data.id, loadingMgr)"
                    >
                      <q-item-section>
                        <q-item-label>文件夹</q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item
                      v-close-popup
                      clickable
                      @click="noteStore.addNote(data.id, loadingMgr)"
                    >
                      <q-item-section>
                        <q-item-label>Markdown</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-item>
              <q-item
                v-close-popup
                clickable
                @click="noteStore.renameFolder(data.id, loadingMgr)"
              >
                <q-item-section>
                  <q-item-label>重命名</q-item-label>
                </q-item-section>
              </q-item>
              <q-item
                v-close-popup
                clickable
                @click="noteStore.moveFolder(data.id)"
              >
                <q-item-section>
                  <q-item-label>移动到</q-item-label>
                </q-item-section>
              </q-item>
              <q-item
                v-close-popup
                clickable
                @click="noteStore.deleteFolder(data.id, loadingMgr)"
              >
                <q-item-section>
                  <q-item-label>删除</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-item-section>
    </template>
    <FolderItem
      v-for="subfolder in data.children"
      :key="subfolder.gid"
      :level="level + 1"
      :data="subfolder"
      :loading-mgr="loadingMgr"
    ></FolderItem>
  </q-expansion-item>
  <q-item
    v-else
    v-ripple
    clickable
    :inset-level="0.5 * level"
    :to="extsObj.routePath(`view/${data.id}`)"
  >
    <q-item-section avatar class="mn-item-icon">
      <q-icon name="article" />
    </q-item-section>
    <q-item-section>
      <q-item-label :lines="2">{{ data.title }}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<script setup>
import { useExtension } from '../extMain';
import { useNoteStore } from '../stores/note';

defineProps({
  data: {
    type: Object,
    required: true,
  },
  loadingMgr: {
    type: Object,
    required: true,
  },
  level: {
    type: Number,
    required: true,
  },
});

const extsObj = useExtension();
const noteStore = useNoteStore();
</script>
