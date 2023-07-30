<template>
  <!-- <div ref="editorRef" /> -->
  <div id="vditor" />
</template>

<script setup>
import { ref, onMounted, nextTick, onBeforeUnmount } from 'vue';
import Vditor from 'vditor';
import 'vditor/dist/index.css';

// const vditor = ref<Vditor | null>(null);
let vditorInst;

const props = defineProps({
  content: {
    type: String,
    default: '',
  },
});

const editorRef = ref(null);

function init() {
  vditorInst = new Vditor(editorRef.value, {
  // vditorInst = new window.Vditor('vditor', {
    height: 720,
    mode: 'sv',
    // toolbarConfig: {
    //   pin: true,
    // },
    cache: {
      enable: false,
    },
    after: () => {
      vditorInst.setValue(props.content);
    },
    // 这里写上传
    // upload: {},
  });
}

// watch(
//   () => props.content,
//   (content) => {
//     if (vditorInst) {
//       vditorInst.setValue(content);
//     }
//   },
//   {
//     immediate: true,
//   }
// );

// 初始化
onMounted(() => {
  nextTick(() => {
    init();
  });
});

// 销毁
onBeforeUnmount(() => {
  if (vditorInst) {
    vditorInst.destroy();
    vditorInst = null;
  }
});

// 获取内容
function getEditValue() {
  return vditorInst.getValue();
}

defineExpose({
  getEditValue,
})
</script>

<style></style>
