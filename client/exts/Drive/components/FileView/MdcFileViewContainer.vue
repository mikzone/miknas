<template>
  <div class="column absolute-full">
    <div class="col-auto full-width">
      <q-toolbar class="bg-primary text-white">
        <q-toolbar-title class="mn-toolbar-title"
          >预览: {{ viewState.curFilePath }}</q-toolbar-title
        >
        <q-btn
          class="q-px-sm"
          stretch
          flat
          :icon="viewGetter.curModeIcon.value"
        >
          <q-menu>
            <q-list dense>
              <q-item-label header class="q-pa-sm">预览类型</q-item-label>
              <q-item
                v-for="(modeInfo, modeKey) in viewState.modeDesc"
                :key="modeKey"
                v-close-popup
                clickable
                @click="viewState.mode = modeKey"
              >
                <q-item-section avatar class="mn-item-icon">
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
          icon="file_download"
          @click="viewOp.downloadCurrent"
        />
      </q-toolbar>
    </div>
    <div class="col full-width relative-position">
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
      <div v-else class="absolute-full row justify-center items-center">
        <h3>暂不支持预览</h3>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  usePreviewView,
  MdcFileViewText,
  MdcFileViewVideo,
  MdcFileViewImg,
} from 'miknas/exts/Drive/shares';

const props = defineProps({
  fsid: {
    type: String,
    required: true,
  },
  fspath: {
    type: String,
    required: true,
  },
});

const { viewState, viewGetter, viewOp } = usePreviewView({
  fsid: props.fsid,
  initFilePath: props.fspath,
});
</script>
