<template>
  <q-drawer
    :side="props.side"
    :width="state.width"
  >
    <div ref="drawerElm" class="fit">
      <slot></slot>
    </div>
    <div class="resize-div" :style="computedStyles" @mousedown="onResize()"></div>
  </q-drawer>
</template>

<script setup>
import { throttle } from 'quasar';
import { computed, reactive, ref } from 'vue';
const props = defineProps({
  side: {
    type: String,
    default: 'left',
  },
  minWidth: {
    type: Number,
    default: 250,
  },
  width: {
    type: Number,
    default: 250,
  },
});

const drawerElm = ref();
const state = reactive({
  width: props.width,
})

const computedStyles = computed(()=>{
  if (props.side == 'right') {
    return {
      left: '-6px',
    }
  } else {
    return {
      right: '-12px',
    }
  }
})

function resizeX(IsRightSide, minW) {
  let startW;
  let startClientX;
  let incRatio = (IsRightSide && -1) || 1;
  let element = drawerElm.value;
  function getIntStyle(key) {
    return parseInt(window.getComputedStyle(element).getPropertyValue(key));
  }
  function dragMouseDown(e) {
    if (e && e.button !== 0) return;
    e = e || window.event;
    e.preventDefault();
    const { clientX } = e;
    startW = getIntStyle('width');
    startClientX = clientX;
    document.addEventListener('mouseup', closeDragElement);
    document.addEventListener('mousemove', throttleEvt);
    // document.addEventListener('mousemove', elementDrag);
  }

  function elementDrag(e) {
    const { clientX } = e;
    let w = startW + (clientX - startClientX) * incRatio;
    if (w < minW) w = minW;
    state.width = w;
  }

  let throttleEvt = throttle(elementDrag, 30);

  function closeDragElement() {
    document.removeEventListener('mouseup', closeDragElement);
    document.removeEventListener('mousemove', throttleEvt);
  }
  return dragMouseDown;
}

const onResize = computed(()=>{
  let IsRightSide = props.side == 'right' ? true : false;
  return resizeX(IsRightSide, props.minWidth);
});
// let IsRightSide = props.side == 'right' ? true : false;
// const onResize = resizeX(IsRightSide, props.minWidth);
// </script>

<style scoped>
.resize-div {
  width: 12px;
  height: 100%;
  background-color: transparent;
  position: absolute;
  top: 0;
  cursor: e-resize;
}
</style>
