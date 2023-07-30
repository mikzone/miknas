<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    @hide="onDialogHide"
  >
    <q-card class="q-dialog-plugin" style="width: 700px; max-width: 80vw">
      <q-card-section class="q-pa-none">
        <q-table
          v-bind="$attrs"
          :rows="fileState.curFiles"
          :columns="fileState.processColumns"
          :filter="fileState.filterTxt"
          :loading="fileGetter.isloading.value"
          :rows-per-page-options="[0]"
          class="mn-sticky-header-table mn-sticky-last-column-table"
          no-data-label="文件列表为空"
          row-key="name"
          style="height: 500px; max-height: 80vh"
        >
          <template #top>
            <q-breadcrumbs class="q-table__title" gutter="xs">
              <q-breadcrumbs-el
                v-for="oneBC in fileGetter.fileBC.value"
                :key="oneBC.name"
                :label="oneBC.name"
                :icon="oneBC.icon"
                :class="
                  oneBC.path === undefined ? '' : 'text-primary cursor-pointer'
                "
                @click="fileOp.gotoPath(oneBC.path)"
              />
            </q-breadcrumbs>
            <q-space />
            <q-btn
              dense
              flat
              round
              color="primary"
              icon="create_new_folder"
              @click="fileOp.newFolder"
            >
              <q-tooltip>新建文件夹</q-tooltip>
            </q-btn>
            <q-btn
              dense
              flat
              round
              color="primary"
              :disable="fileGetter.isloading.value"
              icon="refresh"
              @click="fileOp.tryRefreshFiles()"
            />
            <q-btn
              dense
              flat
              round
              color="primary"
              :disable="fileGetter.isloading.value"
              icon="close"
              @click="onDialogCancel"
            />
          </template>

          <template #body-cell-name="cellProps">
            <q-td :props="cellProps">
              <q-btn
                color="primary"
                flat
                no-caps
                no-wrap
                class="q-pa-xs"
                :icon="cellProps.row.icon"
                :label="cellProps.value"
                @click="fileOp.clickOpen(cellProps.row)"
              ></q-btn>
            </q-td>
          </template>

          <template #body-cell-size="cellProps">
            <q-td auto-width :props="cellProps">
              <div v-if="cellProps.row.viewSize">
                {{ cellProps.row.viewSize }}
              </div>
            </q-td>
          </template>
          <template #bottom>
            <q-space />
            <q-btn :label="confirmLabel" color="primary" @click="onConfirm">
            </q-btn>
          </template>
          <template #no-data>
            <q-space />
            <q-btn :label="confirmLabel" color="primary" @click="onConfirm">
            </q-btn>
          </template>
          <template #loading>
            <q-inner-loading
              showing
              color="primary"
              :label="fileGetter.loadingLabel.value"
            />
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useDialogPluginComponent } from 'quasar';
import { useFileView } from 'miknas/exts/Drive/shares';

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
  fsid: {
    type: String,
    required: true,
  },
  rootPath: {
    type: String,
    default: '',
  },
  confirmLabel: {
    type: String,
    default: '确认',
  },
});

const { fileState, fileGetter, fileOp } = useFileView(
  props.fsid,
  props.rootPath
);

function onConfirm() {
  onDialogOK(fileState.curPath);
}
</script>
<style scoped></style>
