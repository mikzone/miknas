<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    @hide="onDialogHide"
  >
    <q-card class="q-dialog-plugin" style="width: 1200px; max-width: 80vw">
      <q-card-section class="row items-center q-pa-sm bg-secondary text-white">
        <div class="text-h6">{{ props.title }}</div>
        <q-space></q-space>
        <q-btn icon="close" flat round dense @click="onCloseDlg"></q-btn>
      </q-card-section>
      <q-separator />
      <q-card-section class="scroll" style="min-height: 50vh; max-height: 80vh;">
        <div v-if="state.mode == 'raw'" class="q-ma-none" style="white-space: pre-wrap;">{{ props.txt }}</div>
        <!-- <q-input v-if="state.mode == 'raw'" type="textarea" :model-value="props.txt" readonly /> -->
        <ByteMdView v-else :model-value="props.txt" />
      </q-card-section>
      <q-separator />
      <q-card-actions>
        <q-tabs
          v-model="state.mode"
          align="left"
          dense
          indicator-color="purple"
          class="text-primary"
        >
          <q-tab name="raw" no-caps label="原文本" />
          <q-tab name="md" no-caps label="Markdown格式化" />
        </q-tabs>
        <q-space />
        <q-btn color="primary" @click="tryCopy">复制到粘贴板</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useDialogPluginComponent } from 'quasar';
import { gutil, MikCall } from 'miknas/utils';
import { ByteMdView } from 'miknas/exts/Official/shares';
import { reactive } from 'vue';

defineEmits([
  // REQUIRED; need to specify some events that your
  // component will emit through useDialogPluginComponent()
  ...useDialogPluginComponent.emits,
]);

const { dialogRef, onDialogHide, onDialogCancel } = useDialogPluginComponent();
// dialogRef      - Vue ref to be applied to QDialog
// onDialogHide   - Function to be used as handler for @hide on QDialog
// onDialogOK     - Function to call to settle dialog with "ok" outcome
//                    example: onDialogOK() - no payload
//                    example: onDialogOK({ /*...*/ }) - with payload
// onDialogCancel - Function to call to settle dialog with "cancel" outcome

const props = defineProps({
  txt: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    default: '',
  },
});

const state = reactive({
  isRaw: true,
  mode: 'raw',
});

function onCloseDlg() {
  onDialogCancel();
}

function tryCopy() {
  if (!props.txt) return MikCall.sendErrorTips('复制失败,内容为空');
  gutil
    .copyToClipboard(props.txt)
    .then(() => {
      MikCall.sendSuccTips('复制成功!');
    })
    .catch(() => {
      MikCall.sendErrorTips('复制失败,可能是浏览器限制了复制功能');
    });
}
</script>
<style scoped></style>
