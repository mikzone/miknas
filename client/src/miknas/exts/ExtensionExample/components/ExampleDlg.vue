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
          <q-btn
            icon="close"
            flat
            round
            dense
            @click="onCloseDlg"
          ></q-btn>
        </q-card-section>
        <q-card-section>
          <div>这是个测试的dialog</div>
        </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useDialogPluginComponent } from 'quasar';
import { MikCall } from 'miknas/utils';

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
});

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
  let suc = true;
  if (!suc) {
    MikCall.makeConfirm('尚未执行完成，是否关闭？', ()=>{
      onDialogCancel();
    });
  } else {
    onDialogOK();
  }
}

</script>
<style scoped></style>
