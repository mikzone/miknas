
<template>
  <div ref="editordiv"/>
</template>

<script setup>
import Editor from '@toast-ui/editor';
import '@toast-ui/editor/dist/toastui-editor.css';
import { onMounted, ref } from 'vue';

const props = defineProps({
  options: {
    type: Object,
    default: ()=>{return {}},
  },
});

const editordiv = ref();
var tuiInst;

onMounted(() => {
  let options = {
    el: editordiv.value,
    height: '50vh',
    initialEditType: 'markdown',
    // previewStyle: 'vertical',
    previewStyle: 'tab',
    usageStatistics: false,
    language: 'zh',
  }
  for (let [k, v] of Object.entries(props.options)) {
    options[k] = v;
  }
  tuiInst = new Editor(options);
});

function getEditor() {
  return tuiInst;
}

defineExpose({
  getEditor,
})
</script>
