<template>
  <q-btn
    v-if="fileInfo.isFile"
    dense
    flat
    round
    icon="more_vert"
    @click.stop.prevent=""
  >
    <q-menu auto-close>
      <q-list>
        <q-item clickable @click="fileOp.downloadFile(fileInfo.name)">
          <q-item-section>下载</q-item-section>
        </q-item>
        <template v-if="!fileGetter.isReadOnly.value">
          <q-item clickable @click="fileOp.removeFile(fileInfo.name)">
            <q-item-section>删除文件</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.reqCopyFile(fileInfo.name)">
            <q-item-section>复制到</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.reqMvFile(fileInfo.name)">
            <q-item-section>移动到</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.rename(fileInfo.name)">
            <q-item-section>重命名</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.shareFile(fileInfo.name)">
            <q-item-section>分享</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.genTmpDownUrl(fileInfo.name)">
            <q-item-section>生成临时下载链接</q-item-section>
          </q-item>
          <q-item clickable @click="fileOp.startSelectMode(fileInfo)">
            <q-item-section>多选</q-item-section>
          </q-item>
        </template>
      </q-list>
    </q-menu>
  </q-btn>
  <q-btn v-else-if="!fileGetter.isReadOnly.value" dense flat round icon="more_vert" @click.stop.prevent="">
    <q-menu auto-close>
      <q-list>
        <q-item clickable @click="fileOp.removeFile(fileInfo.name)">
          <q-item-section>删除</q-item-section>
        </q-item>
        <q-item clickable @click="fileOp.reqCopyFile(fileInfo.name)">
          <q-item-section>复制到</q-item-section>
        </q-item>
        <q-item clickable @click="fileOp.reqMvFile(fileInfo.name)">
          <q-item-section>移动到</q-item-section>
        </q-item>
        <q-item clickable @click="fileOp.rename(fileInfo.name)">
          <q-item-section>重命名</q-item-section>
        </q-item>
        <q-item clickable @click="fileOp.shareFile(fileInfo.name)">
          <q-item-section>分享</q-item-section>
        </q-item>
        <q-item clickable @click="fileOp.startSelectMode(fileInfo)">
          <q-item-section>多选</q-item-section>
        </q-item>
      </q-list>
    </q-menu>
  </q-btn>
</template>

<script setup>
import { toRefs } from 'vue';

const pros = defineProps({
  fileViewProxy: {
    type: Object,
    required: true,
  },
  fileInfo: {
    type: Object,
    required: true,
  },
});

const { fileViewProxy } = toRefs(pros);
const { fileOp, fileGetter } = fileViewProxy.value;
</script>
