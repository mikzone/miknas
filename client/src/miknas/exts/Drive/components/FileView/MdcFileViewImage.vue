<template>
  <q-img :src="state.src" fit="scale-down" class="absolute-full">
    <template #error>
      <div class="absolute-full flex flex-center text-white">
        <h3>非有效图片</h3>
      </div>
    </template>
  </q-img>
</template>

<script setup>
import { useExtension } from 'miknas/exts/Drive/extMain';
import { reactive } from 'vue';

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
  src: extsObj.serverUrl(`view/${props.fsid}/${props.fspath}`),
});
</script>
