<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    @hide="onDialogHide"
  >
    <q-card class="q-dialog-plugin" style="width: 1200px; max-width: 80vw">
      <q-card-section class="row items-center bg-secondary text-white q-py-sm">
        <div class="text-h6">{{ props.title }}</div>
        <q-space></q-space>
        <q-btn icon="close" flat round dense @click="onCloseDlg"></q-btn>
      </q-card-section>
      <q-card-section>
        <q-tree
          v-model:selected="selected"
          v-model:expanded="state.expandedKeys"
          :nodes="treeNodes"
          node-key="id"
          label-key="label"
          selected-color="secondary"
        />
      </q-card-section>
      <q-separator />
      <q-card-actions>
        <q-input dense standout class="col q-mr-sm" :model-value="curSelectTxt" label="当前选择:" readonly />
        <q-btn color="primary" @click="confirmSelect">确定</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useNoteStore } from 'miknas/exts/Note/stores/note';
import { useDialogPluginComponent } from 'quasar';
import { MikCall } from 'miknas/utils';
import { computed, reactive, ref } from 'vue';

const noteStore = useNoteStore();
const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  exclueId: {
    type: Number,
    default: 0,
  },
});

const state = reactive({
  expandedKeys: [0],
})

defineEmits([
  // REQUIRED; need to specify some events that your
  // component will emit through useDialogPluginComponent()
  ...useDialogPluginComponent.emits,
]);

const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
  useDialogPluginComponent();
// dialogRef      - Vue ref to be applied to QDialog
// onDialogHide   - Function to be used as handler for @hide on QDialog
// onDialogOK     - Function to call to settle dialog with "ok" outcome
//                    example: onDialogOK() - no payload
//                    example: onDialogOK({ /*...*/ }) - with payload
// onDialogCancel - Function to call to settle dialog with "cancel" outcome

function onCloseDlg() {
  onDialogCancel();
}

function dfs(folder) {
  if (!folder.children) return; // 不含chilren说明不是文件夹
  let ret = { id: folder.id, label: folder.name };
  let childList = [];
  for (let child of Object.values(folder.children)) {
    let childRet = dfs(child);
    if (childRet) childList.push(childRet);
  }
  if (childList.length > 0) ret.children = childList;
  return ret;
}

const treeNodes = computed(() => {
  let ret = [dfs(noteStore.rootFolder)];
  return ret;
});

const selected = ref(null);
const curSelectTxt = computed(()=>{
  if (!selected.value) return '';
  return noteStore.folderDict[selected.value].name;
})

function confirmSelect() {
  if (!selected.value) {
    MikCall.sendErrorTips('当前未选择任何文件夹');
    return;
  }
  onDialogOK(selected.value);
}
</script>
<style scoped></style>
