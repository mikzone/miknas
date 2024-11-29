<template>
  <q-page padding>
    <FolderPluginsView
      :key="`${props.spaceId}_${props.fspath}`"
      :space-id="props.spaceId"
      :fspath="props.fspath"
      :kind="props.kind"
    ></FolderPluginsView>
  </q-page>
</template>

<script setup>
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
</script>
