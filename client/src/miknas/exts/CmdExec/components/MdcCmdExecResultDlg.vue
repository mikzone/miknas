<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    transition-show="slide-up"
    transition-hide="slide-down"
    @hide="onDialogHide"
  >
    <q-layout
      container
      view="hHh lpR fFf"
      style="width: 1200px; max-width: 80vw"
      class="bg-primary text-white"
    >
      <q-header class="bg-secondary text-white">
        <q-toolbar>
          <q-toolbar-title> 运行情况 </q-toolbar-title>
          <q-btn icon="close" flat round dense @click="onCloseDlg"></q-btn>
        </q-toolbar>
      </q-header>

      <q-page-container>
        <q-page>
          <MdcCmdExecResult
            v-if="jobId"
            ref="execResultRef"
            :job-id="props.jobId"
            :ace-lang="props.aceLang"
            @finish-exec="onFinishExec"
          ></MdcCmdExecResult>
        </q-page>
      </q-page-container>
    </q-layout>
  </q-dialog>
</template>

<script setup>
import { useDialogPluginComponent } from 'quasar';
import { MikCall } from 'miknas/utils';
import { reactive } from 'vue';
import MdcCmdExecResult from './MdcCmdExecResult.vue';

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

const props = defineProps({
  jobId: {
    type: String,
    required: true,
  },
  aceLang: {
    type: String,
    default: 'ace/mode/text',
  },
});

const state = reactive({
  jobItem: undefined,
});

function onFinishExec(jobItem) {
  state.jobItem = jobItem;
}

function onCloseDlg() {
  if (!state.jobItem) {
    MikCall.makeConfirm('尚未执行完成，是否关闭？关闭并不会取消任务。', () => {
      onDialogCancel();
    });
  } else {
    onDialogOK(state.jobItem);
  }
}
</script>
<style scoped></style>
