<template>
  <ExtensionPage>
    <template #toolbar>
      <template v-if="state.shareInfo">
        <div class="row no-wrap items-center justify-start q-gutter-sm mn-toolbar-title" style="font-size: 16px">
          <div>文件分享</div>
          <div>--</div>
          <q-icon name="person" />
          <div class="text-no-wrap">{{ state.shareInfo.uid }}</div>
          <div>/</div>
          <q-icon name="folder_shared" />
          <div class="text-no-wrap">{{ state.shareInfo.name }}</div>
        </div>
      </template>
    </template>
    <q-page>
      <MdcDriveAliveView
        v-if="state.shareInfo && props.kind != 'check'"
        :fsid="curFsid"
        fsroot=""
        :fsrela="props.fsrela"
        :kind="props.kind"
        :extra-conf="myExtraConf"
      ></MdcDriveAliveView>
      <div
        v-else-if="state.needPwd"
        class="absolute-full flex-center q-pa-md row q-col-gutter-sm"
      >
        <div class="col-12 col-sm-6 col-md-4">
          <q-input
            v-model="state.inputPwd"
            filled
            label="请输入提取码"
            @keyup.enter="tryVerify"
          />
          <q-btn
            size="lg"
            class="q-mt-lg full-width"
            color="primary"
            @click="tryVerify"
            >提取文件</q-btn
          >
        </div>
      </div>
      <div
        v-else-if="state.errInfo"
        class="text-red absolute-full column flex-center"
      >
        {{ state.errInfo.why }}
      </div>
      <div v-else class="absolute-full column flex-center">加载中</div>
    </q-page>
  </ExtensionPage>
</template>

<script setup>
import { ExtensionPage } from 'miknas/exts/Official/shares';
import { MdcDriveAliveView } from 'miknas/exts/Drive/shares';
import { computed, onMounted, reactive } from 'vue';
import { useExtension } from '../extMain';
import { MikCall, gutil } from 'miknas/utils';
import { useRouter } from 'vue-router';
const router = useRouter();

const props = defineProps({
  shareid: {
    type: String,
    required: true,
  },
  fsrela: {
    type: String,
    default: '',
  },
  kind: {
    type: String,
    default: 'list',
  },
});

const curFsid = `S_${props.shareid}`;

const state = reactive({
  needPwd: false,
  shareInfo: undefined,
  inputPwd: '',
});

const extsObj = useExtension();

const myExtraConf = computed(() => {
  return {
    isReadOnly: true,
  };
});

async function tryQuery() {
  let iRet = await extsObj.mcpost('viewShare', { sid: props.shareid });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    state.errInfo = iRet;
    return;
  }
  let info = iRet.ret;
  if (info.needPwd) {
    state.needPwd = true;
    return;
  }
  info.viewBts = gutil.formatTs(info.bts * 1000);
  if (info.intv > 0) {
    info.viewEts = gutil.formatTs((info.bts + info.intv) * 1000);
  }
  state.shareInfo = info;
  if (props.kind == 'check') {
    router.replace({
      name: extsObj.routeName('slist'),
      params: { shareid: props.shareid, routeSubPath: '' },
    });
  }
}

async function tryVerify() {
  if (!state.inputPwd) {
    MikCall.sendErrorTips('请输入提取码');
  }
  let iRet = await extsObj.mcpost('verifyShare', {
    sid: props.shareid,
    pwd: state.inputPwd,
  });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  tryQuery();
}

onMounted(() => {
  tryQuery();
});
</script>
