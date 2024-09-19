<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="q-gutter-y-md">
    <div
      v-for="confItem in myform.formConfs"
      :key="confItem.id"
    >
      <label>{{ confItem.title }}</label>
      <div>
        <component
          :is="confItem.component"
          v-model="myform.state.formData[confItem.id]"
          :conf-item="confItem"
        />
      </div>
      <span class="help-block" v-html="confItem.desc"></span>
    </div>
    <div class="row q-gutter-x-xs">
      <q-btn
        label="保存"
        color="primary"
        @click="trySaveSetting"
      />
      <q-btn v-close-popup label="关闭" color="primary" />
    </div>
  </div>
</template>
<script setup>
import { useFormView } from 'miknas/exts/Official/shares';

const props = defineProps({
  name: {
    type: String,
    default: '',
  },
  formConfs: {
    type: Object,
    required: true,
  },
  initData: {
    type: Object,
    default: () => {return {}},
  },
});

const myform = useFormView(props.formConfs, props.initData);

</script>
