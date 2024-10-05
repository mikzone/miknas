<template>
  <q-item v-ripple clickable>
    <q-item-section avatar>
      <q-icon name="enhanced_encryption" />
    </q-item-section>
    <q-item-section>
      <q-item-label>{{ attachInfo.jsonData.name }}</q-item-label>
      <q-item-label caption>{{ attachInfo.jsonData.hint }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-btn dense round flat icon="more_vert" @click.stop.prevent="">
        <q-menu>
          <q-list dense>
            <q-item v-close-popup clickable @click="viewSecret">
              <q-item-section>
                <q-item-label>查看加密内容</q-item-label>
              </q-item-section>
            </q-item>
            <q-item v-close-popup clickable @click="noteInfoUtil.startModifyAttach(attachInfo)">
              <q-item-section>
                <q-item-label>修改</q-item-label>
              </q-item-section>
            </q-item>
            <q-item v-close-popup clickable @click="noteInfoUtil.startDeleteAttach(attachInfo)">
              <q-item-section>
                <q-item-label>删除</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </q-item-section>
  </q-item>
</template>

<script setup>
import { attachCfgs } from 'miknas/exts/Note/helpers';
import { openTextCopyDlg } from 'miknas/exts/Official/shares';
import { computed } from 'vue';

const props = defineProps({
  attachInfo: {
    type: Object,
    required: true,
  },
  noteInfoUtil: {
    type: Object,
    required: true,
  },
});

const attachConf = computed(()=>{
  return attachCfgs[props.attachInfo.type]
})

async function viewSecret() {
  let formData = await attachConf.value.genFormData(props.attachInfo.jsonData);
  if (formData) {
    openTextCopyDlg({ title: '加密内容如下:', txt: formData.txt, showSwitch: true });
  }
}

</script>
