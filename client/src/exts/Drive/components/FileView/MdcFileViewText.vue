<template>
  <div class="absolute-full">
    <MdcAceEditor
      v-model="state.txt"
      class="absolute-full"
      style="border: 3px solid #795548;"
      :ace-lang="state.aceLang"
      :options="{readOnly: true}"
    >
    </MdcAceEditor>
    <q-inner-loading :showing="loadingMgr.isloading.value" color="primary" :label="loadingMgr.loadingLabel.value" />
  </div>
</template>

<script setup>
import { useExtension } from 'miknas/exts/Drive/extMain';
import { useMikLoading, MdcAceEditor } from 'miknas/exts/Official/shares';
import { onMounted, reactive } from 'vue';
import { MikCall } from 'miknas/utils';
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

const extsObj = useExtension();

let modelist = window.ace.require('ace/ext/modelist');
const state = reactive({
  txt: '',
  aceLang: modelist.getModeForPath(props.fspath).mode,
});

const loadingMgr = useMikLoading();

onMounted(async () => {
  let stateName = `加载文本中`;
  loadingMgr.addLoadingState(stateName);
  let iRet = await extsObj.mcpost('viewTxt', { fsid: props.fsid, fspath: props.fspath });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let result = iRet.ret;
  state.txt = result;
  loadingMgr.removeLoadingState(stateName);
});
</script>
