<template>
  <q-page>
    <ContainerLayout @keydown.ctrl.s.prevent.stop="saveNote">
      <template #toolbar>
        <q-space />
        <q-btn
          v-if="isView"
          flat
          round
          dense
          icon="edit"
          :to="extsObj.routePath(`edit/${noteid}`)"
          replace
        />
        <q-btn
          v-else
          flat
          round
          dense
          icon="preview"
          :to="extsObj.routePath(`view/${noteid}`)"
          replace
        />
        <q-btn v-if="isModify" flat round dense icon="save" @click="saveNote" />
      </template>
      <template v-if="isView">
        <q-card class="no-shadow mn-bordered">
          <q-item>
            <q-item-section avatar>
              <q-icon size="lg" name="article" />
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-h6 text-link text-weight-bold">{{
                state.newTitle
              }}</q-item-label>
              <q-item-label caption>{{ state.note.modify }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-card-section>
            <ByteMdView v-model="state.newContent" />
          </q-card-section>
        </q-card>
      </template>
      <template v-else>
        <q-input
          v-model="state.newTitle"
          label="标题"
          outlined
          class="bg-white q-mb-sm"
          hide-bottom-space
          :rules="[DataRule.isNotEmptyString]"
        >
          <template #prepend>
            <q-icon name="title" />
          </template>
        </q-input>
        <ByteMdInput
          v-model="state.newContent"
          class="fullheight-bytemd"
          style="height: calc(100vh - 250px)"
        />
      </template>
      <q-card class="no-shadow mn-bordered q-mt-md">
        <NoteAttachList :note-info="state.note" />
      </q-card>
    </ContainerLayout>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      :label="loadingMgr.loadingLabel.value"
    />
  </q-page>
</template>

<script setup>
import { computed, onMounted, reactive, watch } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { DataRule } from 'miknas/exts/Official/shares';
import { MikCall } from 'miknas/utils';
import { useNoteStore } from '../stores/note';
import { ByteMdInput } from 'miknas/exts/Official/shares';
import ContainerLayout from '../layouts/ContainerLayout.vue';
import { ByteMdView } from 'miknas/exts/Official/shares';
import NoteAttachList from '../components/NoteAttachList.vue';
import useExtension from '../extMain';

const extsObj = useExtension();

const props = defineProps({
  noteStrid: {
    type: String,
    required: true,
  },
  initView: {
    type: Boolean,
    default: true,
  },
});

const noteid = computed(() => {
  return parseInt(props.noteStrid);
});

const isView = computed(() => {
  return props.initView;
});

const isModify = computed(() => {
  return (
    state.note.title != state.newTitle || state.note.content != state.newContent
  );
});

const loadingMgr = useMikLoading();
const state = reactive({
  newTitle: '',
  newContent: '',
  note: {},
  // isView: props.initView,
});

const noteStore = useNoteStore();

async function saveNote() {
  let ret = await noteStore.saveNote({
    id: state.note.id,
    title: state.newTitle,
    content: state.newContent,
  });
  if (!ret) return;
  MikCall.sendSuccTips('保存成功');
  state.note = ret;
}

async function refreshNote() {
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);
  let ret = await noteStore.getNote(noteid.value);
  if (!ret) return;
  state.note = ret;
  state.newTitle = ret.title;
  state.newContent = ret.content;
  loadingMgr.removeLoadingState(stateName);
}

watch(noteid, () => {
  refreshNote();
});

// watch(isView, () => {
//   state.isView = props.initView;
// });

onMounted(() => {
  refreshNote();
});
</script>

<style scoped></style>
