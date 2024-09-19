<template>
  <div>
    <div class="q-pa-md text-h5">{{ state.result }}</div>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      :label="loadingMgr.loadingLabel.value"
    />
  </div>
</template>
<script setup>
import { onMounted, reactive } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { MikCall } from 'miknas/utils';

const loadingMgr = useMikLoading();
const props = defineProps({
  name: {
    type: String,
    default: '',
  },
});

const state = reactive({
  result: '',
});

const extsObj = useExtension();

async function queryResult() {
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);

  let iRet = await extsObj.mcpost('hello', { name: props.name });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let result = iRet.ret;
  state.result = result;
  loadingMgr.removeLoadingState(stateName);
}

onMounted(() => {
  queryResult();
});
</script>
