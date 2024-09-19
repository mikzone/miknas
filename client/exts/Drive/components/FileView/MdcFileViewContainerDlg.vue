<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    maximized
    persistent
    style="background: rgba(0, 0, 0, 1)"
    @hide="onDialogHide"
  >
    <q-layout view="hHh lpR fFf">
      <q-header class="bg-secondary text-white">
        <q-toolbar class="overflow-auto">
          <div class="text-subtitle2 col overflow-auto hide-scrollbar">
            {{ viewState.curFileStat.name }}
            <q-tooltip>
              修改时间: {{ viewState.curFileStat.viewModify }}
            </q-tooltip>
          </div>
          <q-space />
          <q-btn
            class="q-px-sm col-auto"
            stretch
            flat
            :label="viewState.curFileStat.viewSize"
          >
          </q-btn>
        </q-toolbar>
      </q-header>
      <q-page-container>
        <q-page>
          <!-- <div class="absolute-full" @click="onDialogCancel"></div> -->
          <div>
            <MdcFileViewImg
              v-if="viewState.mode == 'img'"
              :key="viewState.curFilePath + '6'"
              :fsid="props.fsid"
              :fspath="viewState.curFilePath"
              :view-op="viewOp"
            />
            <MdcFileViewText
              v-else-if="viewState.mode == 'text'"
              :key="viewState.curFilePath + '1'"
              :fsid="props.fsid"
              :fspath="viewState.curFilePath"
            />
            <MdcFileViewVideo
              v-else-if="viewState.mode == 'video'"
              :key="viewState.curFilePath + '2'"
              :fsid="props.fsid"
              :fspath="viewState.curFilePath"
            />
            <div
              v-else
              class="absolute-full row justify-center items-center text-white"
            >
              <h3>暂不支持预览</h3>
            </div>
          </div>
        </q-page>
      </q-page-container>

      <q-footer reveal elevated class="bg-secondary text-white">
        <q-toolbar class="overflow-auto">
          <q-space />
          <q-btn class="q-px-sm" stretch flat :icon="viewGetter.curModeIcon.value">
            <q-menu>
              <q-list dense>
                <q-item
                  v-for="(modeInfo, modeKey) in viewState.modeDesc"
                  :key="modeKey"
                  v-close-popup
                  clickable
                  @click="viewState.mode = modeKey"
                >
                  <q-item-section avatar>
                    <q-icon :name="modeInfo.icon" />
                  </q-item-section>
                  <q-item-section>{{ modeInfo.name }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
          <q-btn
            class="q-px-sm"
            stretch
            flat
            icon="navigate_before"
            @click="viewOp.chooseNext(-1)"
          />
          <q-btn
            class="q-px-sm"
            stretch
            flat
            icon="navigate_next"
            @click="viewOp.chooseNext(1)"
          />
          <q-btn
            class="q-px-sm"
            stretch
            flat
            icon="file_download"
            @click="viewOp.downloadCurrent"
          />
          <!-- <q-btn class="q-px-sm" stretch flat icon="fullscreen" /> -->
          <q-btn
            class="q-px-sm"
            stretch
            flat
            icon="close"
            @click="onDialogCancel"
          />
          <q-space />
        </q-toolbar>
        <q-inner-loading
          class="bg-secondary"
          :showing="viewState.loadingMgr.isloading"
          color="primary"
          :label="viewState.loadingMgr.loadingLabel"
        />
      </q-footer>
    </q-layout>
  </q-dialog>
</template>

<script setup>
import { useDialogPluginComponent } from 'quasar';
import {
  usePreviewView,
  MdcFileViewText,
  MdcFileViewVideo,
  MdcFileViewImg,
} from 'miknas/exts/Drive/shares';

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
  fsid: {
    type: String,
    required: true,
  },
  initFilePath: {
    type: String,
    required: true,
  },
});

const { viewState, viewGetter, viewOp } = usePreviewView({
  fsid: props.fsid,
  initFilePath: props.initFilePath,
});
</script>
<style scoped></style>
