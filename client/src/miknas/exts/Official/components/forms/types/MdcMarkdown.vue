<template>
  <!-- <q-input
    dense
    hide-bottom-space
    :model-value="props.modelValue"
    :rules="confItem.rules"
    @update:model-value="updateModelValue"
  /> -->
  <q-field
    :model-value="props.modelValue"
    :rules="confItem.rules"
  >
    <template #control>
      <ByteMdInput
        class="full-width"
        :model-value="props.modelValue"
        v-bind="$props.byteMdOpts"
        @update:model-value="updateModelValue"
      />
    </template>
  </q-field>
</template>

<script setup>
import { ByteMdInput } from 'miknas/exts/Official/shares';
const props = defineProps({
  confItem: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: String,
    required: true,
  },
  byteMdOpts: {
    type: Object,
    default: () => {return {}},
  },
});

const emit = defineEmits(['update:modelValue']);

function updateModelValue(newValue) {
  emit('update:modelValue', newValue);
}
</script>
