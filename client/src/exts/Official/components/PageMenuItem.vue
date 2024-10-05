<template>
  <q-item
    v-if="inSubMenu"
    v-close-popup
    v-bind="$attrs"
    active-class="mn-active"
    clickable
    :to="to"
    :exact="exact"
    :tag="tag"
  >
    <q-item-section>{{ title }}</q-item-section>
  </q-item>
  <q-route-tab
    v-else-if="to"
    no-caps
    :to="to"
    :exact="exact"
    :label="title"
    v-bind="$attrs"
  ></q-route-tab>
  <q-btn
    v-else
    flat
    no-wrap
    no-caps
    stretch
    :type="tag"
    :to="to"
    :exact="exact"
    :label="title"
    v-bind="$attrs"
  ></q-btn>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  to: {
    type: String,
    default: undefined,
  },
  tag: {
    type: String,
    default: 'a',
  },
  inSubMenu: {
    type: Boolean,
    default: false,
  },
  exact: {
    type: Boolean,
    default: true,
  },
});
</script>
