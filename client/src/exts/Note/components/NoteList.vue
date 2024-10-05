<template>
  <div>
    <div ref="scrollElm"></div>
    <q-card
      v-for="noteInfo in state.noteList"
      :key="noteInfo.id"
      class="no-shadow mn-bordered"
    >
      <q-item>
        <q-item-section avatar>
          <q-icon size="md" name="article" />
        </q-item-section>

        <q-item-section>
          <q-item-label>
            <RouterLink
              class="text-h6 text-link text-weight-bold text-decoration-none"
              :to="extsObj.routePath(`edit/${noteInfo.id}`)"
            >
              {{ noteInfo.title }}
            </RouterLink>
          </q-item-label>
          <q-item-label caption>{{ noteInfo.modify }}</q-item-label>
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
                        v-for="(attachConf, attachType) in attachCfgs"
                        :key="attachType"
                        v-close-popup
                        clickable
                        @click="openAddAttach(noteInfo, attachType)"
                      >
                        <q-item-section>
                          <q-item-label>{{ attachConf.name }}</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-item>
                <q-item v-close-popup clickable @click="tryMoveNote(noteInfo)">
                  <q-item-section>
                    <q-item-label>移动到</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item v-close-popup clickable @click="deleteNote(noteInfo)">
                  <q-item-section>
                    <q-item-label>删除</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-item-section>
      </q-item>

      <q-card-section>
        <ByteMdView v-model="noteInfo.content" />
      </q-card-section>

      <template v-if="noteInfo.noteAttachs.length">
        <q-separator />
        <NoteAttachList :note-info="noteInfo" />
      </template>
    </q-card>
    <q-banner v-if="!state.noteList.length" class="text-center">
      <template v-if="state.search"
        >找不到与
        <span class="text-h6 text-primary text-weight-bolder">{{
          state.search
        }}</span>
        相关的笔记</template
      >
      <template v-else>暂无笔记</template>
    </q-banner>
    <div class="flex flex-center">
      <q-pagination :model-value="state.curPage" :max="pageCnt" @update:model-value="updateCurPage" />
    </div>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      :label="loadingMgr.loadingLabel.value"
    />
  </div>
</template>

<script setup>
import { useMikLoading } from 'miknas/exts/Official/shares';
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useNoteStore } from '../stores/note';
import useExtension from '../extMain';
import { attachCfgs, moveNote, openAddAttach } from '../helpers';
import { RouterLink } from 'vue-router';
import NoteAttachList from './NoteAttachList.vue';
import { scroll } from 'quasar'
import { ByteMdView } from 'miknas/exts/Official/shares';

const props = defineProps({
  folderid: {
    // 如果带有该id，表示加载该文件夹下的文章
    type: Number,
    default: 0,
  },
});

const state = reactive({
  cntPerPage: 10,
  cnt: 0,
  search: '',
  curPage: 1,
  noteList: [],
});

const pageCnt = computed(() => {
  return Math.ceil(state.cnt / state.cntPerPage);
});

const loadingMgr = useMikLoading();
const noteStore = useNoteStore();
const extsObj = useExtension();

async function refreshNoteList(pageNum, keyword) {
  pageNum = pageNum || 1;
  keyword = keyword || '';
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);
  let result = await noteStore.listNotes({
    folder: props.folderid,
    search: keyword,
    pageNum: pageNum,
    pageSize: state.cntPerPage,
  });
  loadingMgr.removeLoadingState(stateName);
  if (!result) return;
  state.search = keyword;
  state.noteList = result.notes;
  state.cnt = result.total;
  scrollToTop();
}

async function deleteNote(noteInfo) {
  let noteid = noteInfo.id;
  let stateName = `正在删除`;
  loadingMgr.addLoadingState(stateName);
  let result = await noteStore.deleteNote(noteInfo);
  loadingMgr.removeLoadingState(stateName);
  if (!result) return;
  state.noteList = state.noteList.filter((obj) => obj.id !== noteid);
}

async function tryMoveNote(noteInfo) {
  let newNoteInfo = await moveNote(noteInfo);
  if (!newNoteInfo) return;
  if (newNoteInfo.folder == noteInfo.folder) return;
  let noteid = newNoteInfo.id
  state.noteList = state.noteList.filter((obj) => obj.id !== noteid);
}

const scrollElm = ref();

const { getScrollTarget, setVerticalScrollPosition } = scroll
function scrollToTop () {
  const target = getScrollTarget(scrollElm.value)
  const offset = 0
  const duration = 200
  setVerticalScrollPosition(target, offset, duration)
}

watch(
  () => props.folderid,
  () => {
    // 调用updateCurPage防止 pagination 变更时 refreshNoteList两次
    state.search = '';
    updateCurPage(1);
  }
);

function updateCurPage(newPage) {
  state.curPage = newPage;
  refreshNoteList(state.curPage, state.search);
}

onMounted(() => {
  refreshNoteList();
});

function doSearch(keyword) {
  refreshNoteList(null, keyword);
}

defineExpose({
  doSearch,
});
</script>
