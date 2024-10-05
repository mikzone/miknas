<template>
  <div class="column absolute-full">
    <div class="col-auto full-width">
      <q-toolbar class="full-width bg-white">
        <q-btn
          flat
          round
          dense
          icon="vertical_split"
          @click="noteStore.toggleDrawerOpen"
        />
        <slot name="toolbar"></slot>
      </q-toolbar>
      <slot name="after-toolbar"></slot>
      <q-separator />
    </div>
    <div class="col q-pa-md scroll full-width">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { useNoteStore } from '../stores/note';

const noteStore = useNoteStore();
</script>

<style scoped></style>
