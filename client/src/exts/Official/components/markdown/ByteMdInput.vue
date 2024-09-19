<template>
  <Editor
    :value="props.modelValue"
    :plugins="plugins"
    @change="updateModelValue"
  />
</template>

<script setup>
import gfm from '@bytemd/plugin-gfm';
import highlight from '@bytemd/plugin-highlight';
import 'bytemd/dist/index.css';
import 'highlight.js/styles/default.css';
import '../../css/github.css';
import { Editor } from '@bytemd/vue-next';
import { drivePlugin } from '../../helpers/bytemdHelpers';

const plugins = [gfm(), highlight(), drivePlugin()];

const props = defineProps({
  modelValue: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue']);

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}
</script>

<style>
.fullheight-bytemd .bytemd {
  height: 100%;
}
.bytemd-fullscreen.bytemd {
  z-index: 10000;
}
</style>
