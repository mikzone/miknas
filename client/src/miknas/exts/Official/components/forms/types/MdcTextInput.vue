<template>
  <q-input
    dense
    hide-bottom-space
    :model-value="props.modelValue"
    :rules="confItem.rules"
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
    type: String,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue'])

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}

</script>
