<template>
  <div class="absolute-full column" style="border: 3px solid #795548">
    <MdcAceEditor
      v-model="state.txt"
      class="col"
      :ace-lang="state.aceLang"
      :options="calcedAceOptions"
    >
    </MdcAceEditor>
    <div class="col-auto q-pa-xs q-gutter-xs">
      <q-btn
        v-for="btnInfo in aceBtns"
        :key="btnInfo.key"
        size="sm"
        :label="btnInfo.label"
        :color="(state.aceOptions[btnInfo.key] && 'primary') || 'grey'"
        @click="toggleAceOption(btnInfo.key)"
      />
    </div>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      color="primary"
      :label="loadingMgr.loadingLabel.value"
    />
  </div>
</template>

<script setup>
import { useExtension } from 'miknas/exts/Drive/extMain';
import { useMikLoading, MdcAceEditor } from 'miknas/exts/Official/shares';
import { onMounted, reactive, computed } from 'vue';
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
  aceOptions: {
    wrap: false,
    showGutter: false,
    showLineNumbers: false,
    readOnly: true,
  },
});

const calcedAceOptions = computed(() => {
  let newAceOptions = { ...state.aceOptions };
  newAceOptions.showGutter = newAceOptions.showLineNumbers;
  return newAceOptions;
});

const aceBtns = [
  { key: 'wrap', label: '换行' },
  { key: 'showLineNumbers', label: '显示行号' },
  { key: 'readOnly', label: '只读' },
];

function toggleAceOption(key) {
  state.aceOptions[key] = !state.aceOptions[key];
}

const loadingMgr = useMikLoading();

onMounted(async () => {
  let stateName = `加载文本中`;
  loadingMgr.addLoadingState(stateName);
  let iRet = await extsObj.mcpost('viewTxt', {
    fsid: props.fsid,
    fspath: props.fspath,
  });
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
