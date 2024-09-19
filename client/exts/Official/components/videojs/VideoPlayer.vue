<template>
  <video ref="videoPlayer" class="video-js vjs-big-play-centered"></video>
</template>

<script setup>
import 'video.js/dist/video-js.css';
import videojs from 'video.js';
import 'videojs-mobile-ui/dist/videojs-mobile-ui.css';
import 'videojs-mobile-ui';
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue';
const props = defineProps({
  options: {
    type: Object,
    default() {
      return {};
    },
  },
});

const state = reactive({
  player: null,
});

const videoPlayer = ref(null);
onMounted(() => {
  state.player = videojs(videoPlayer.value, props.options, () => {
    state.player.log('onPlayerReady', state.player);
  });
  state.player.mobileUi({
    fullscreen: {
      lockToLandscapeOnEnter: true,
    },
  });
});

onBeforeUnmount(() => {
  if (state.player) {
    state.player.dispose();
  }
});
</script>
