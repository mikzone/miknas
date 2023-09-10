<!-- eslint-disable vue/no-v-html -->
<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    transition-show="slide-up"
    transition-hide="slide-down"
    @hide="onDialogHide"
  >
    <q-card style="width: 1200px; max-width: 90vw">
      <q-card-section class="q-pa-sm row items-center text-white bg-secondary">
        <div class="col text-h6 mn-word-break-all">{{ props.title }}</div>
        <q-btn icon="close" no-caps flat round dense @click="onCloseDlg"></q-btn>
      </q-card-section>
      <q-card-section style="max-height: calc(100vh - 115px); overflow: auto">
        <q-form class="q-gutter-md" @submit="onSubmit">
          <div v-for="confItem in myform.formConfs" :key="confItem.id">
            <div class="row item-centers q-gutter-xs">
              <label class="text-subtitle2">{{ confItem.title }}</label>
              <q-space />
              <q-btn
                v-for="helpAction in confItem.helpActions"
                :key="helpAction.label"
                :label="helpAction.label"
                size="sm"
                dense
                color="primary"
                @click="helpAction.func(myform.state.formData)"
              >
              </q-btn>
            </div>
            <span class="text-caption block" v-html="confItem.desc"></span>
            <div>
              <component
                :is="confItem.component"
                v-bind="confItem.componentProps"
                v-model="myform.state.formData[confItem.id]"
                :conf-item="confItem"
              />
            </div>
          </div>
          <q-btn type="submit" :label="props.confirmLabel" color="primary" />
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useFormView } from 'miknas/exts/Official/shares';
import { useDialogPluginComponent } from 'quasar';

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
  title: {
    type: String,
    default: '填写表单',
  },
  formConfs: {
    type: Array,
    required: true,
  },
  initData: {
    type: Object,
    default: () => {
      return {};
    },
  },
  confirmLabel: {
    type: String,
    default: '确认',
  },
});

const myform = useFormView(props.formConfs, props.initData);

function onCloseDlg() {
  onDialogCancel();
  // MikCall.makeConfirm('尚未完成当前表单，是否确认关闭？', () => {
  //   onDialogCancel();
  // });
}

function onSubmit() {
  let ret = myform.action.tryGetValidFormData();
  if (ret) {
    onDialogOK(ret);
  }
  return false;
}
</script>
<style scoped></style>
