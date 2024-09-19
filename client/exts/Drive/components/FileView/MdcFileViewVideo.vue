<template>
  <div class="absolute-full row justify-center items-center">
    <VideoPlayer :options="state.videoOptions" />
  </div>
</template>

<script setup>
import { useExtension } from 'miknas/exts/Drive/extMain';
import { reactive } from 'vue';
import { VideoPlayer } from 'miknas/exts/Official/shares';

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
});

const state = reactive({
  videoOptions: {
    controls: true,
    fill: true,
    sources: [
      {
        src: extsObj.serverUrl(`view/${props.fsid}/${props.fspath}`),
      }
    ]
  },
})
</script>
