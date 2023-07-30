<template>
  <q-btn-toggle
    :model-value="props.modelValue"
    no-caps
    :options="props.confItem.selectOptions"
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
    type: Number,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue'])

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}

</script>
