<template>
  <q-page>
    <MdcDriveAliveView
      :fsid="props.fsid"
      :fsroot="props.fsroot"
      :fsrela="props.fsrela"
      :kind="props.kind"
      :extra-conf="myExtraConf"
    ></MdcDriveAliveView>
  </q-page>
</template>

<script setup>
import { MdcDriveAliveView } from 'miknas/exts/Drive/shares';

const props = defineProps({
  fsid: {
    type: String,
    required: true
  },
  fsroot: {
    type: String,
    default: ''
  },
  fsrela: {
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
