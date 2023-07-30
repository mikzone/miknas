<!-- eslint-disable vue/no-v-html -->
<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    @hide="onDialogHide"
  >
    <q-layout
      container
      view="hHh lpR fFf"
      style="width: 1200px; max-width: 80vw"
      class="bg-white"
    >
      <q-header class="bg-secondary text-white">
        <q-toolbar>
          <q-toolbar-title class="mn-toolbar-title"> {{ props.title }} </q-toolbar-title>
          <q-btn icon="close" flat round dense @click="onCloseDlg"></q-btn>
        </q-toolbar>
      </q-header>

      <q-page-container>
        <q-page>
          <q-scroll-area class="absolute-full q-pa-md">
          <q-form class="q-gutter-lg q-pt-md" @submit="onSubmit">
            <div v-for="confItem in myform.formConfs" :key="confItem.id">
              <label class="text-h6">{{ confItem.title }}</label>
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
          </q-scroll-area>
        </q-page>
      </q-page-container>
    </q-layout>
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
