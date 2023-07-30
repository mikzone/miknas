<template>
  <div>
    <slot></slot>
    <q-dialog v-model="showDialog" persistent>
      <q-card style="width: 1200px; max-width: 80vw">
        <q-card-section class="row items-center">
          <div class="text-h6">运行情况</div>
          <q-space></q-space>
          <q-btn
            v-show="!innerIsLoading"
            v-close-popup
            icon="close"
            flat
            round
            dense
          ></q-btn>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <MdcCmdExecResult
            v-if="jobId"
            ref="execResultRef"
            :job-id="jobId"
            @finish-exec="onFinishExec"
          ></MdcCmdExecResult>
        </q-card-section>
      </q-card>
    </q-dialog>
    <q-inner-loading :showing="innerIsLoading"></q-inner-loading>
  </div>
</template>

<script>
import { MikCall } from 'miknas/utils';
import MdcCmdExecResult from 'miknas/exts/CmdExec/MdcCmdExecResult.vue';
export default {
  name: 'MdcCmdExec',
  components: { MdcCmdExecResult },
  props: {
    isLoading: Boolean,
  },
  data: function () {
    return {
      innerIsLoading: false,
      showDialog: false,
      jobId: '',
      jobItem: null,
      sucCb: null,
    };
  },
  methods: {
    updatePropVal: function (propName, val) {
      this.$emit(`update:${propName}`, val);
      this[`inner_${propName}`] = val;
    },
    updateIsLoading: function (val) {
      this.updatePropVal('isLoading', val);
    },
    tryExec: async function ({ url, data, sucCb }) {
      // # 执行指令
      if (this.sucCb) {
        MikCall.sendErrorTips('当前已经有指令在执行中');
        return;
      }
      this.jobId = '';
      this.sucCb = sucCb;
      this.updateIsLoading(true);
      this.showDialog = true;

      let iRet = await MikCall.mcpost(url, data);
      if (!iRet.suc) {
        MikCall.alertRespErrMsg(iRet);
        this.updateIsLoading(false);
        this.showDialog = false;
        this.sucCb = null;
        return;
      }
      let cbRet = iRet.ret;
      this.jobId = cbRet.jobId;
    },
    onFinishExec: function (jobItem) {
      this.updateIsLoading(false);
      let sucCb = this.sucCb;
      this.sucCb = null;
      if (sucCb) sucCb(jobItem);
    },
    closeRunningDialog: function () {
      this.showDialog = false;
    },
    setAceLang: function (lang) {
      this.$refs.execResultRef.setAceLang(lang);
    },
  },
};
</script>
