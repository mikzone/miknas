<template>
  <q-select
    :model-value="props.modelValue"
    :rules="confItem.rules"
    :options="state.curHintList"
    hide-selected
    fill-input
    use-input
    input-debounce="0"
    dense
    hide-bottom-space
    @filter="filterFn"
    @input-value="setModel"
  />
</template>

<script setup>
import { computed } from 'vue';
import { reactive } from 'vue';

const props = defineProps({
  confItem: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Number],
    required: true,
  },
});

const inHintList = computed(()=>{
  return props.confItem.hintList || [];
})
const emit = defineEmits(['update:modelValue']);
const state = reactive({
  curHintList: [],
})

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}

function filterFn(val, update) {
  update(() => {
    const needle = val.toLocaleLowerCase();
    state.curHintList = inHintList.value.filter(
      (v) => v.toLocaleLowerCase().indexOf(needle) > -1
    );
  });
}

function setModel(val) {
  updateModelValue(val);
}
</script>
