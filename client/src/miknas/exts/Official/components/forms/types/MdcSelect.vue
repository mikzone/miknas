<template>
  <q-select
    :model-value="props.modelValue"
    :rules="confItem.rules"
    :options="confItem.selectOptions"
    emit-value
    map-options
    dense
    hide-bottom-space
    @update:model-value="updateModelValue"
  />
</template>

<script setup>
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

const emit = defineEmits(['update:modelValue']);

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}
</script>
