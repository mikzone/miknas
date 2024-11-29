<template>
  <q-list separator bordered>
    <q-item
      v-for="space in workStore.spaces"
      :key="space.Id"
      clickable
      :to="extsObj.routePath(`list/${space.Id}`)"
    >
      <q-item-section avatar>
        <q-icon name="folder_open" />
      </q-item-section>
      <q-item-section>
        <q-item-label lines="1">
          <q-badge color="secondary" :label="space.Id" />
        </q-item-label>
        <q-item-label>{{ space.Name }}</q-item-label>
        <q-item-label caption>{{ space.Path }}</q-item-label>
      </q-item-section>
    </q-item>
  </q-list>
  <q-list separator bordered>
    <q-item v-for="defInfo in workStore.pluginSimples" :key="defInfo.Id">
      <q-item-section avatar>
        <q-icon name="settings_input_component" />
      </q-item-section>
      <q-item-section>
        <q-item-label>
          {{ defInfo.Id }}
          <q-badge color="grey" :label="defInfo.Version" />
        </q-item-label>
        <q-item-label caption>{{ defInfo.Desc }}</q-item-label>
      </q-item-section>
    </q-item>
  </q-list>
</template>
<script setup>
import { onMounted } from 'vue';
import { useExtension } from '../extMain';
import { useWorkStore } from '../stores/work';

const workStore = useWorkStore();

const extsObj = useExtension();

onMounted(() => {
  workStore.refreshMainInfo();
});
</script>
