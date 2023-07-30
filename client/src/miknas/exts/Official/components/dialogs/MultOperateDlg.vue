<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    transition-show="slide-up"
    transition-hide="slide-down"
    @hide="onDialogHide"
  >
    <q-card class="q-dialog-plugin" style="min-width: 60vw">
      <q-card-section>
        <div class="text-h6">
          {{ props.title }}({{ FinishCnt }} / {{ totalCnt }})
        </div>
        <q-linear-progress
          dark
          stripe
          rounded
          size="20px"
          :animation-speed="200"
          :value="progressValue"
          color="primary"
          class="q-mt-sm q-mb-sm"
        />
        <div v-if="state.curJobKey" class="q-mb-sm">
          {{ props.jobDescFn(state.curJobKey) }}
        </div>
        <q-list
          dense
          class="rounded-borders scroll"
          style="max-height: 60vh"
        >
          <q-item-label v-if="failJobCnt > 0" class="text-red" header>
            共 {{ failJobCnt }} 项处理失败
          </q-item-label>
          <q-item
            v-for="(iRet, jobKey) in state.failJobs"
            :key="jobKey"
            v-ripple
            clickable
            class="text-red"
          >
            <q-item-section>
              <q-item-label>{{ props.jobDescFn(jobKey) }}</q-item-label>
              <q-item-label caption>{{ iRet.why || iRet }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-btn
                v-show="state.jobState == 'done'"
                flat
                color="primary"
                dense
                size="sm"
                label="重试"
                @click="tryDoJob(jobKey)"
              />
            </q-item-section>
          </q-item>
          <q-item-label v-if="SucJobCnt > 0" class="text-green" header>
            共 {{ SucJobCnt }} 项处理成功
          </q-item-label>
        </q-list>
      </q-card-section>
      <q-card-actions>
        <q-space />
        <q-btn
          v-if="state.jobState == 'done'"
          color="primary"
          label="完成"
          @click="onOKClick"
        />
        <q-btn
          v-if="state.jobState == 'running'"
          color="primary"
          label="中止剩余任务"
          @click="onCancelLeft"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { MikCall } from 'miknas/utils';
import { useDialogPluginComponent } from 'quasar';
import { computed, onMounted, reactive } from 'vue';

defineEmits([
  // REQUIRED; need to specify some events that your
  // component will emit through useDialogPluginComponent()
  ...useDialogPluginComponent.emits,
]);

const { dialogRef, onDialogHide, onDialogOK } = useDialogPluginComponent();
// dialogRef      - Vue ref to be applied to QDialog
// onDialogHide   - Function to be used as handler for @hide on QDialog
// onDialogOK     - Function to call to settle dialog with "ok" outcome
//                    example: onDialogOK() - no payload
//                    example: onDialogOK({ /*...*/ }) - with payload
// onDialogCancel - Function to call to settle dialog with "cancel" outcome

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  jobKeys: {
    type: Array,
    required: true,
  },
  jobDescFn: {
    type: Function,
    required: true,
  },
  jobProcFn: {
    type: Function,
    required: true,
  },
});

const state = reactive({
  jobState: 'waiting',
  isCancel: false,
  curJobKey: null,
  failJobs: {},
  doneJobs: {},
});

const totalCnt = computed(()=>{
  return props.jobKeys.length;
})

const failJobCnt = computed(()=>{
  return Object.keys(state.failJobs).length;
})

const SucJobCnt = computed(()=>{
  return Object.keys(state.doneJobs).length;
})

const FinishCnt = computed(()=>{
  return SucJobCnt.value + failJobCnt.value;
})

const progressValue = computed(()=>{
  if (totalCnt.value <= 0) return 0;
  return (1.0 * FinishCnt.value) / totalCnt.value;
})

async function tryDoJob(jobKey) {
  if (jobKey === undefined) return;
  state.curJobKey = jobKey;
  if (jobKey in state.failJobs) delete state.failJobs[jobKey];
  let iRet = await props.jobProcFn(jobKey);
  if (iRet === null || iRet === undefined) {
    state.failJobs[jobKey] = MikCall.failRet('miknas: 处理函数未返回任何值');
  } else if (!iRet.suc) {
    state.failJobs[jobKey] = iRet;
  } else {
    state.doneJobs[jobKey] = iRet;
  }
  state.curJobKey = null;
}

async function startDoAllJobs() {
  state.jobState = 'running'
  for (let jobKey of props.jobKeys) {
    await tryDoJob(jobKey);
    if (state.isCancel) break;
  }
  state.jobState = 'done'
}

function onCancelLeft() {
  state.isCancel = true;
}

function onOKClick() {
  onDialogOK({
    dones: state.doneJobs,
    fails: state.failJobs,
  });
}

onMounted(()=>{
  startDoAllJobs();
})

</script>
