<template>
  <q-page>
    <div class="row">
      <div class="col-12">
        <MdcDriveAliveView
          :fsid="`Cw_${props.spaceId}`"
          :fsrela="props.fspath"
          :kind="props.kind"
          :extra-conf="myExtraConf"
        ></MdcDriveAliveView>
      </div>
      <div v-if="props.kind == 'list'" class="col-12">
        <FolderPluginsView :space-id="props.spaceId" :fspath="props.fspath"></FolderPluginsView>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { MdcDriveAliveView } from 'miknas/exts/Drive/shares';
import FolderPluginsView from '../components/FolderPluginsView.vue';

const props = defineProps({
  spaceId: {
    type: String,
    required: true
  },
  fspath: {
    type: String,
    default: ''
  },
  kind: {
    type: String,
    default: 'list'
  }
});

const myExtraConf = {
  // viewFn: null
};
</script>
