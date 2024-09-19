<template>
  <div ref="mydiv" class="absolute-full" style="overflow: hidden">
    <img
      ref="myimg"
      class="myimg"
      :src="state.src"
      @load="imgOnload"
    />
  </div>
</template>

<script setup>
import { useExtension } from 'miknas/exts/Drive/extMain';
import { reactive, ref } from 'vue';
import Hammer from 'hammerjs';

const extsObj = useExtension();
const props = defineProps({
  fsid: {
    type: String,
    required: true,
  },
  fspath: {
    type: String,
    required: true,
  },
  viewOp: {
    type: Object,
    required: true,
  }
});

const state = reactive({
  src: extsObj.serverUrl(`view/${props.fsid}/${props.fspath}`),
  debugMsg: '',
});

function hammerIt(elm, parentElm) {
  elm.style.transformOrigin = 'left top';
  let hammertime = new Hammer(elm, {});
  hammertime.get('pinch').set({
    enable: true,
  });
  hammertime.get('pan').set({ direction: Hammer.DIRECTION_ALL });
  var posX = 0,
    posY = 0,
    scale = 1,
    lastScale = 1,
    scaleCenter = {}, // pinch中心点
    lastPosX = 0,
    lastPosY = 0,
    isPinching = false,
    isPanning = false,
    isAnim = false,
    el = elm;

  function calcOrigScale(){
    let origWidth = elm.naturalWidth;
    let origHeight = elm.naturalHeight;
    let rect1 = parentElm.getBoundingClientRect();
    let containerWidth = rect1.width;
    let containerHeight = rect1.height;
    let scale = Math.min(containerHeight / origHeight, containerWidth / origWidth, 1);
    return scale;
  }

  function getRestrict(){
    let origScale = calcOrigScale();
    let cx = origScale * elm.naturalWidth;
    let cy = origScale * elm.naturalHeight;
    let rect1 = parentElm.getBoundingClientRect();
    let bx = rect1.width;
    let by = rect1.height;
    let ret = {
      minX: bx - scale * (bx + cx) / 2,
      minY: by - scale * (by + cy) / 2,
      maxX: (cx - bx) * scale / 2,
      maxY: (cy - by) * scale / 2,
    }
    if (ret.minX > ret.maxX) {
      let avg = (ret.minX + ret.maxX) / 2;
      ret.minX = avg;
      ret.maxX = avg;
    }
    if (ret.minY > ret.maxY) {
      let avg = (ret.minY + ret.maxY) / 2;
      ret.minY = avg;
      ret.maxY = avg;
    }
    return ret;
  }

  function getScale(times) {
    return times / calcOrigScale();
  }

  function getCenter(ev) {
    let rect1 = parentElm.getBoundingClientRect();
    // rect1 不一定是占满屏幕的，所以rect1的xy分别是它的坐上偏移量
    let ret = {
      x: ev.center.x - rect1.x,
      y: ev.center.y - rect1.y,
    }
    return ret;
  }

  function UpdateNeedTransition(flag) {
    if (isAnim != flag) {
      isAnim = flag;
      if (flag) {
        el.style.transition = '100ms';
      } else {
        el.style.transition = '';
      }
    }
  }

  function UpdateImg() {
    if (!isPinching && !isPanning) {
      let rest = getRestrict();
      if (posX > rest.maxX) {
        posX = rest.maxX;
      }
      if (posX < rest.minX) {
        posX = rest.minX;
      }
      if (posY > rest.maxY) {
        posY = rest.maxY;
      }
      if (posY < rest.minY) {
        posY = rest.minY;
      }
      lastPosX = posX;
      lastPosY = posY;
      UpdateNeedTransition(true)
    } else {
      UpdateNeedTransition(false)
    }
    if (!isPinching) {
      lastScale = scale;
    }
    let transform = `translate3d(${posX}px, ${posY}px, 0) scale3d(${scale}, ${scale}, 1)`;
    el.style.webkitTransform = transform;
  }

  hammertime.on('doubletap', function (ev) {
    try {
      if (scale != 1) {
        scale = 1;
        posX = 0;
        posY = 0;
        UpdateImg();
      } else {
        scale = 2;
        scaleCenter = getCenter(ev);
        posX = scaleCenter.x - (scaleCenter.x - lastPosX) * scale / lastScale;
        posY = scaleCenter.y - (scaleCenter.y - lastPosY) * scale / lastScale;
        UpdateImg();
      }
    } catch (err) {
      console.log(err);
    }
  });

  hammertime.on('pinch pinchend', function (ev) {
    //pinch
    if (ev.type == 'pinch') {
      if (!isPinching) {
        isPinching = true;
        scaleCenter = getCenter(ev);
      }
      let maxScale = getScale(2);
      scale = Math.max(1, Math.min(lastScale * ev.scale, maxScale));
      posX = scaleCenter.x - (scaleCenter.x - lastPosX) * scale / lastScale;
      posY = scaleCenter.y - (scaleCenter.y - lastPosY) * scale / lastScale;
    } else if (ev.type == 'pinchend') {
      isPinching = false;
    }
    UpdateImg();
  });

  hammertime.on('pan panend', function (ev) {
    //pan
    posX = lastPosX + ev.deltaX;
    posY = lastPosY + ev.deltaY;
    if (ev.type == 'pan') {
      if (!isPanning) isPanning = true;
    } else if (ev.type == 'panend') {
      isPanning = false;
      if (scale <= 1){
        let ratio = 0.2;
        let checkwidth = parentElm.clientWidth * ratio;
        if (posX < -checkwidth) {
          props.viewOp.chooseNext(1, ['img']);
        } else if (posX > checkwidth) {
          props.viewOp.chooseNext(-1, ['img']);
        }
      }
    }
    UpdateImg();
  });
}

const myimg = ref(null);
const mydiv = ref(null);
function imgOnload() {
  let mydivElm = mydiv.value;
  let myimgElm = myimg.value;
  hammerIt(myimgElm, mydivElm);
}
</script>

<style scoped>
.myimg{
  width: 100%;
  height: 100%;
  object-fit: scale-down;
}
</style>
