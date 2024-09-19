<template>
  <q-item v-if="inSubMenu" clickable>
    <q-item-section>{{ title }}</q-item-section>
    <q-item-section side>
      <q-icon name="keyboard_arrow_right" />
    </q-item-section>

    <q-menu anchor="top end" self="top start" :offset="[4, 0]">
      <q-list class="mn-header-menu">
        <slot></slot>
      </q-list>
    </q-menu>
  </q-item>
  <q-btn-dropdown
    v-else
    flat
    stretch
    no-caps
    :label="title"
    menu-anchor="bottom start"
    menu-self="top start"
    :menu-offset="[0, 4]"
  >
    <q-list class="mn-header-menu">
      <slot></slot>
    </q-list>
  </q-btn-dropdown>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  inSubMenu: {
    type: Boolean,
    default: false,
  },
});
</script>
