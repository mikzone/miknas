<template>
  <div>
    <q-list separator bordered>
      <q-item
        v-for="jobConf in props.pluginDetail.pluginDef.Jobs"
        :key="jobConf.Id"
        clickable
        @click="execJob(jobConf.Id)"
      >
        <q-item-section avatar>
          <q-icon :name="jobConf.Icon || 'construction'" />
        </q-item-section>
        <q-item-section>
          <q-item-label lines="1">
            <q-badge color="secondary" :label="jobConf.Id" />
          </q-item-label>
          <q-item-label>{{ jobConf.Name }}</q-item-label>
          <q-item-label caption>{{ jobConf.Shell }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
    <q-inner-loading :showing="loadingMgr.isloading.value" :label="loadingMgr.loadingLabel.value" />
  </div>
</template>
<script setup>
// import { computed, onMounted, reactive } from 'vue';
import { useExtension } from '../extMain';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { MikCall } from 'miknas/utils';
import { fetchResult } from 'miknas/exts/CmdExec/exec_cmd_util';
import { coOpenFormDlg } from 'miknas/exts/Official/shares';
// import { useWorkStore } from '../stores/work';

// const workStore = useWorkStore();

const loadingMgr = useMikLoading();

const props = defineProps({
  spaceId: {
    type: String,
    required: true
  },
  fspath: {
    type: String,
    default: ''
  },
  pluginDetail: {
    type: Object,
    required: true
  }
});

// const state = reactive({
//   pluginStats: [],
//   tab: undefined
// });

const extsObj = useExtension();

async function execJob(jobId, formData, ignoreConfirm) {
  let jobConf = props.pluginDetail.pluginDef._jobMap[jobId];
  if (!jobConf) return;
  if (!ignoreConfirm && jobConf.Confirm) {
    let isOk = await MikCall.coMakeConfirm(`是否确认执行: ${jobConf.Name}`);
    if (!isOk) return;
  }
  let stateName = `正在请求执行`;
  loadingMgr.addLoadingState(stateName);
  let iRet = await extsObj.mcpost('execPluginJob', {
    spaceId: props.spaceId,
    fspath: props.fspath,
    pluginId: props.pluginDetail.pluginId,
    jobId: jobId,
    formData: formData
  });
  loadingMgr.removeLoadingState(stateName);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let result = iRet.ret;
  if (result.nextAction == 'ShowExec') {
    let jobInfo = result.jobInfo;
    let execJobId = jobInfo.jobId;
    if (!execJobId) {
      MikCall.sendErrorTips('执行错误，创建Job失败');
      return;
    }
    fetchResult({ jobId: execJobId });
  } else if (result.nextAction == 'FillForm') {
    let formProps = result.form;
    if (!formProps.title) {
      formProps.title = jobConf.Name;
    }
    let [isOk, newFormData] = await coOpenFormDlg(formProps);
    if (!isOk) return;
    return await execJob(jobId, newFormData, true);
  }
}
</script>
