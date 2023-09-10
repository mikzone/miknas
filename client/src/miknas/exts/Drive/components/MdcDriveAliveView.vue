<template>
  <KeepAlive :max="15">
    <MdcDriveListView
      v-if="props.kind == 'list'"
      :key="route.path"
      :fsid="props.fsid"
      :root-path="props.fsroot"
      :extra-conf="myExtraConf"
    />
  </KeepAlive>
  <MdcFileViewContainer
    v-if="props.kind == 'view'"
    class="bg-black text-white"
    :fsid="props.fsid"
    :fspath="fspath"
  />
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router';
import MdcDriveListView from './MdcDriveListView.vue';
import MdcFileViewContainer from './FileView/MdcFileViewContainer.vue';
import { computed } from 'vue';
import { FileUtil } from 'miknas/exts/Drive/shares';
import { useExtension } from '../extMain';
import { gutil } from 'miknas/utils';

const router = useRouter();
const route = useRoute();

const props = defineProps({
  fsid: {
    type: String,
    required: true,
  },
  fsroot: {
    type: String,
    default: '',
  },
  fsrela: {
    type: String,
    default: '',
  },
  kind: {
    type: String,
    default: 'list',
  },
  extraConf: {
    type: Object,
    default: () => {
      return {};
    },
  },
});

const fspath = computed(()=>{
  let ret = FileUtil.contactFolderName(props.fsroot, props.fsrela);
  return ret
})

const myExtraConf = computed(()=>{
  let ret = {
    openDirFn: (fsrela) => {
      let newloc = { params: { routeSubPath: fsrela } };
      router.push(newloc);
    },
    viewFn: (fsrela) => {
      if (route.meta.fsViewRouteName) {
          let newloc = { name: route.meta.fsViewRouteName, params: { routeSubPath: fsrela } };
          router.push(newloc);
        }
        else {
          let extsObj = useExtension();
          let newFspath = FileUtil.contactFolderName(props.fsroot, fsrela);
          router.push(extsObj.routePath(`view/${props.fsid}/${newFspath}`));
        }
    },
    initFsrela: route.params.routeSubPath,
  }
  if (props.extraConf) {
    gutil.mergeDict(ret, props.extraConf)
  }
  return ret
})

// 该组件提供给那些需要跳转目录，有需要能回到之前位置的
</script>
