<template>
  <q-dialog
    ref="dialogRef"
    no-backdrop-dismiss
    no-esc-dismiss
    @hide="onDialogHide"
  >
    <q-card style="min-width: 300px" class="q-dialog-plugin">
      <q-card-section class="q-pa-none">
        <q-uploader
          style="width: 100%"
          class="no-uploader-overlay"
          :factory="factoryFn"
          multiple
          auto-upload
          field-name="files"
          @failed="onUploadFailed"
          @removed="onRemoveFiles"
          @factory-failed="onFactoryFailed"
          @uploading="onStartUploading"
          @uploaded="onUploaded"
        >
          <template #header="scope">
            <div class="row no-wrap items-center q-pa-sm q-gutter-xs">
              <q-btn
                v-if="scope.queuedFiles.length > 0"
                icon="clear_all"
                round
                dense
                flat
                @click="
                  MikCall.makeConfirm(
                    '是否清空队列文件(包含未上传的 和 上传失败的)',
                    scope.removeQueuedFiles
                  )
                "
              >
                <q-tooltip>清空</q-tooltip>
              </q-btn>
              <q-btn
                v-if="scope.uploadedFiles.length > 0"
                icon="done_all"
                round
                dense
                flat
                @click="scope.removeUploadedFiles"
              >
                <q-tooltip>移除已上传文件</q-tooltip>
              </q-btn>
              <q-spinner v-if="scope.isUploading" class="q-uploader__spinner" />
              <div class="col">
                <div class="q-uploader__title">文件上传</div>
                <div class="q-uploader__subtitle">
                  {{ scope.uploadSizeLabel }} / {{ scope.uploadProgressLabel }}
                </div>
              </div>
              <q-btn
                v-if="scope.canAddFiles"
                type="a"
                icon="add_box"
                round
                dense
                flat
                @click="scope.pickFiles"
              >
                <q-uploader-add-trigger />
                <q-tooltip>添加文件</q-tooltip>
              </q-btn>
              <q-btn
                v-if="scope.canUpload"
                icon="cloud_upload"
                round
                dense
                flat
                @click="scope.upload"
              >
                <q-tooltip>上传</q-tooltip>
              </q-btn>

              <q-btn
                v-if="scope.isUploading"
                icon="clear"
                round
                dense
                flat
                @click="scope.abort"
              >
                <q-tooltip>取消上传</q-tooltip>
              </q-btn>
              <q-btn
                v-if="!scope.isUploading"
                icon="close"
                round
                dense
                flat
                @click="onCloseClick(scope.queuedFiles.length)"
              >
                <q-tooltip>关闭</q-tooltip>
              </q-btn>
            </div>
          </template>
          <template #list="scope">
            <q-list separator>
              <q-item v-for="file in scope.files" :key="file.__key">
                <q-item-section v-if="file.__img" avatar>
                  <q-avatar rounded>
                    <img :src="file.__img.src" />
                  </q-avatar>
                </q-item-section>
                <q-item-section v-else-if="file.__status == 'failed'" avatar>
                  <q-icon name="warning" color="red" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="full-width ellipsis">
                    {{ file.name }}
                  </q-item-label>

                  <!-- <q-item-label caption>
                    Status: {{ file.__status }}
                  </q-item-label> -->

                  <q-item-label caption>
                    {{ file.__sizeLabel }} / {{ file.__progressLabel }}
                  </q-item-label>

                  <q-item-label v-if="file.__mderr" class="text-red" caption>
                    <q-icon name="priority_high" /> {{ file.__mderr }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-btn
                    v-if="file.__status == 'failed'"
                    class="gt-xs"
                    flat
                    dense
                    round
                    color="red"
                    icon="delete"
                    @click="scope.removeFile(file)"
                  />
                  <q-btn
                    v-else-if="file.__status == 'idle'"
                    class="gt-xs"
                    flat
                    dense
                    round
                    icon="delete"
                    @click="scope.removeFile(file)"
                  />
                  <q-btn
                    v-else-if="file.__status == 'uploaded'"
                    class="gt-xs text-green"
                    flat
                    dense
                    round
                    icon="done"
                  />
                  <q-spinner
                    v-else-if="file.__status == 'uploading'"
                    class="gt-xs"
                    size="2em"
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </template>
        </q-uploader>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { Dialog, useDialogPluginComponent } from 'quasar';
import { MaxCntLocker, MikCall } from 'miknas/utils';
import { ref } from 'vue';

const props = defineProps({
  factory: {
    type: Function,
    required: true,
  },
  maxRunningCnt: {
    type: Number,
    default: 2,
  },
});

const mylocker = new MaxCntLocker(props.maxRunningCnt);

function releaseFilesLock(files) {
  let curFile = files[0];
  let lockId = curFile.__mdlockid;
  mylocker.release(lockId);
}

defineEmits([
  // REQUIRED; need to specify some events that your
  // component will emit through useDialogPluginComponent()
  ...useDialogPluginComponent.emits,
]);

const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
  useDialogPluginComponent();
// dialogRef      - Vue ref to be applied to QDialog
// onDialogHide   - Function to be used as handler for @hide on QDialog
// onDialogOK     - Function to call to settle dialog with "ok" outcome
//                    example: onDialogOK() - no payload
//                    example: onDialogOK({ /*...*/ }) - with payload
// onDialogCancel - Function to call to settle dialog with "cancel" outcome

const sucCount = ref(0);


function onCloseHandle(){
  if (sucCount.value > 0) onDialogOK();
  else onDialogCancel();
}

function onCloseClick(queuedFilesNum) {
  if (queuedFilesNum > 0) {
    Dialog.create({
      message: '是否退出上传文件，并忽视当前 未上传 或 上传失败 的文件?',
      cancel: true,
      persistent: true,
    }).onOk(() => {
      onCloseHandle();
    });
  } else{
    onCloseHandle();
  }
}

async function factoryFn(files) {
  if (files.length != 1) {
    throw '只能上传单个文件'
  }
  let curFile = files[0];
  let lockId = await mylocker.acquire();
  curFile.__mdlockid = lockId;
  return await props.factory(files);
}

function onUploadFailed(info) {
  let { files, xhr } = info;
  try {
    let ret = xhr.responseText;
    let errMsg = ret;
    try {
      ret = JSON.parse(ret);
      errMsg = ret.why;
    } catch (error) {
      if (!xhr.getAllResponseHeaders()) {
        errMsg = '请求被中断'
      }
      else {
        errMsg = '请求发生异常'
        console.warn(files, xhr, error);
      }
    }

    MikCall.sendErrorTips(errMsg);
    releaseFilesLock(files);
    for (let file of files) file.__mderr = errMsg;
  } catch (error) {
    console.warn(files, xhr, error);
  }
}

function onFactoryFailed(err, files) {
  let errMsg = `${err}`;
  for (let file of files) file.__mderr = errMsg;
  releaseFilesLock(files);
}

function onRemoveFiles(files) {
  releaseFilesLock(files);
}

function onUploaded(info) {
  let { files } = info;
  sucCount.value += files.length;
  releaseFilesLock(files);
}

function onStartUploading(info) {
  let { files } = info;
  for (let file of files) file.__mderr = '';
}

</script>
<style>
.no-uploader-overlay .q-uploader__overlay {
  display: none;
}
</style>
