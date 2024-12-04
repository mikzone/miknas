<template>
  <div ref="xtermRef" class="xterm"></div>
</template>

<script setup>
// 参考
// 作者：xinfei
// 链接：https://juejin.cn/post/7356234354118639668
import { onBeforeUnmount, onMounted, ref } from 'vue';

import '@xterm/xterm/css/xterm.css';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { debounce } from 'quasar';
const props = defineProps({
  xtermOption: {
    type: Object,
    default: () => {
      return {};
    }
  }
});

const xtermRef = ref(null);
var xterm;

var fitAddon;

// 初始化xterm
const initXterm = () => {
  const options = {
    fontSize: 14,
    convertEol: true, //启用时，光标将设置为下一行的开头
    // scrollback: 50, //终端中的回滚量
    disableStdin: false, //是否应禁用输入
    cursorBlink: true, //光标闪烁
    theme: {
      foreground: '#ECECEC', //字体
      background: '#000000', //背景色
      cursor: 'help' //设置光标
    }
  };

  const xtermOptions = Object.assign(options, props.xtermOption);
  xterm = new Terminal(xtermOptions);

  fitAddon = new FitAddon();
  xterm.loadAddon(fitAddon);
  xterm.open(xtermRef.value);
  fitAddon.fit();
  window.myfitAddon = fitAddon;
  xterm.focus();
};

const resize = debounce(() => {
  // 适应父容器的大小
  console.log('resize');
  fitAddon.fit();
}, 100);

function write(txt) {
  xterm.write(txt);
}

defineExpose({
  resize,
  write
});

onMounted(() => {
  initXterm();
  window.addEventListener('resize', resize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize);
});
</script>

<style lang="scss">
.xterm {
  width: 100%;
  height: 100%;
}
</style>
