<template>
  <div>
    <div v-if="state.jobItem" class="column absolute-full">
      <q-markup-table square class="col-auto" flat dense separator="none">
        <tbody>
          <tr>
            <td width="120px">cmd</td>
            <td>
              <pre class="td-cmd">{{ state.jobItem.cmd }}</pre>
            </td>
          </tr>
          <tr>
            <td>cwd</td>
            <td>{{ state.jobItem.cwd }}</td>
          </tr>
          <tr>
            <td>nameSpace</td>
            <td>{{ state.jobItem.nameSpace }}</td>
          </tr>
        </tbody>
      </q-markup-table>
      <MdcAceEditor
        v-model="state.jobItem.out"
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
                    <q-item-section>终止任务</q-item-section>
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
    required: true,
  },
  aceLang: {
    type: String,
    default: 'ace/mode/text',
  },
});

const state = reactive({
  jobItem: null,
  isUnMount: false,
});

const emit = defineEmits(['finishExec']);

async function tryRefreshExecResult() {
  if (
    state.jobItem &&
    ['done', 'canceled', 'errstop'].includes(state.jobItem.runningState)
  ) {
    // 如果是已经完成了的话，不用处理
    return;
  }
  if (!props.jobId) return;
  let reqArgs = { jobId: props.jobId };
  let iRet = await extsObj.mcpost('queryJobResult', reqArgs);
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    return;
  }
  let newJobItem = iRet.ret;
  state.jobItem = newJobItem;
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
