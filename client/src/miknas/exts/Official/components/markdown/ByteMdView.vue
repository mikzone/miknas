<template>
  <Viewer :value="props.modelValue" :plugins="plugins" />
</template>

<script setup>
import gfm from '@bytemd/plugin-gfm'
import highlight from '@bytemd/plugin-highlight'
import 'bytemd/dist/index.css'
import 'highlight.js/styles/default.css'
// import '../css/juejin.scss'
import '../../css/github.css'
import { Viewer } from '@bytemd/vue-next'

const plugins = [
  gfm(),
  highlight(),
]

const props = defineProps({
  modelValue: {
    type: String,
    required: true,
  },
});

// const emit = defineEmits(['update:modelValue'])

</script>
