<template>
  <div>
    <div v-if="state.jobItem" class="column absolute-full">
      <div square class="col-auto bg-teal text-white" flat dense separator="none">
        <q-list dense>
          <q-item>
            <q-item-section>
              <q-item-label :lines="2">{{ state.jobItem.title }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <div>
                <q-btn dense flat round icon="info" color="white">
                  <q-menu>
                    <q-list>
                      <q-item>
                        <q-item-section>
                          <q-item-label caption>nameSpace</q-item-label>
                          <q-item-label>{{ state.jobItem.nameSpace }}</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item>
                        <q-item-section>
                          <q-item-label caption>cmd</q-item-label>
                          <q-item-label>{{ state.jobItem.cmd }}</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item>
                        <q-item-section>
                          <q-item-label caption>cwd</q-item-label>
                          <q-item-label>{{ state.jobItem.cwd }}</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
                <slot name="job-header-side-btns" />
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </div>
      <MdcAceEditor
        v-model="state.stdoutInfo.txt"
        class="col"
        auto-scroll-to-end
        :ace-lang="props.aceLang"
      ></MdcAceEditor>
      <div class="col-auto">
        <q-banner
          v-if="state.jobItem.runningState == 'done'"
          dense
          inline-actions
          class="text-white bg-positive"
          >当前任务执行完毕! {{ state.jobItem.failtxt }}</q-banner
        >
        <q-banner
          v-else-if="state.jobItem.runningState == 'canceled'"
          dense
          inline-actions
          class="text-white bg-negative"
          >当前任务已被取消：{{ state.jobItem.failtxt }}
        </q-banner>
        <q-banner
          v-else-if="state.jobItem.runningState == 'errstop'"
          dense
          inline-actions
          class="text-white bg-negative"
          >当前任务已停止：{{ state.jobItem.failtxt }}
        </q-banner>
        <q-banner v-else dense inline-actions class="bg-grey-3 text-black"
          >当前任务状态: {{ state.jobItem.runningState }}
          {{ state.jobItem.failtxt }}
          <template #action>
            <q-btn color="primary" label="取消">
              <q-menu>
                <q-list style="min-width: 100px">
                  <q-item
                    v-close-popup
                    clickable
                    @click="tryCancel(state.jobItem.jobId, 'terminate')"
                  >
                    <q-item-section>终止任务(terminate)</q-item-section>
                  </q-item>
                  <q-item v-close-popup clickable @click="tryCancel(state.jobItem.jobId, 'kill')">
                    <q-item-section>终止任务(kill)</q-item-section>
                  </q-item>
                  <q-item v-close-popup clickable @click="tryCancel(state.jobItem.jobId, 'none')">
                    <q-item-section>终止任务(none)</q-item-section>
                  </q-item>
                  <q-item v-close-popup clickable @click="tryCancel(state.jobItem.jobId, 'SIGINT')">
                    <q-item-section>终止任务(SIGINT)</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </template>
        </q-banner>
      </div>
    </div>
  </div>
</template>

<script setup>
import { MikCall } from 'miknas/utils';
import { MdcAceEditor } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { onMounted, onBeforeUnmount, reactive } from 'vue';
let extsObj = useExtension();

const props = defineProps({
  jobId: {
    type: String,
    required: true
  },
  initReadStdoutStart: {
    // 从那个位置开始读取输出
    type: Number,
    default: -1
  },
  aceLang: {
    type: String,
    default: 'ace/mode/text'
  }
});

const state = reactive({
  jobItem: null,
  stdoutInfo: {
    // 聚合的stdout
    start: undefined,
    end: undefined,
    txt: ''
  },
  isUnMount: false
});

const emit = defineEmits(['finishExec']);

async function tryRefreshExecResult() {
  if (state.jobItem && ['done', 'canceled', 'errstop'].includes(state.jobItem.runningState)) {
    // 如果是已经完成了的话，不用处理
    return;
  }
  if (!props.jobId) return;
  let readStdoutStart;
  if (state.stdoutInfo.end === undefined) {
    readStdoutStart = props.initReadStdoutStart;
  } else {
    readStdoutStart = state.stdoutInfo.end + 1;
  }
  let reqArgs = {
    jobId: props.jobId,
    readStdoutStart: readStdoutStart
  };
  let iRet = await extsObj.mcpost('queryJobResult', reqArgs);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let newJobItem = iRet.ret;
  state.jobItem = newJobItem;
  if (state.stdoutInfo.end === undefined) {
    // 首次初始化，使用服务端传来的值
    state.stdoutInfo.start = newJobItem.stdoutInfo.start;
    state.stdoutInfo.end = newJobItem.stdoutInfo.end;
    state.stdoutInfo.txt = newJobItem.stdoutInfo.txt;
    // console.log('newJobItem1', newJobItem.stdoutInfo);
  } else if (state.stdoutInfo.end + 1 == newJobItem.stdoutInfo.start) {
    // 将文本拼接起来
    if (newJobItem.stdoutInfo.txt.length > 0) {
      state.stdoutInfo.end = newJobItem.stdoutInfo.end;
      state.stdoutInfo.txt += newJobItem.stdoutInfo.txt;
      // console.log('newJobItem2', newJobItem.stdoutInfo);
    }
  }
  if (!['done', 'canceled', 'errstop'].includes(newJobItem.runningState)) {
    if (!state.isUnMount) setTimeout(tryRefreshExecResult, 200);
  } else {
    // 处理完成了，直接回调事件
    emit('finishExec', newJobItem);
  }
}

async function tryCancel(jobId, killType) {
  let isOk = await MikCall.coMakeConfirm(`是否取消执行该任务`);
  if (!isOk) return;
  let iRet = await extsObj.mcpost('reqCancelJob', { jobId: jobId, killType: killType });
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    await tryRefreshExecResult();
    return;
  }
  let cbRet = iRet.ret;
  MikCall.sendSuccTips(cbRet);
  await tryRefreshExecResult();
}

onMounted(() => {
  tryRefreshExecResult();
});

onBeforeUnmount(() => {
  state.isUnMount = true;
});
</script>
<style lang="sass">
.td-cmd
  max-height: 80px
  overflow: auto
  margin: 0
  color: white
  background: #333
  padding: 10px
</style>
