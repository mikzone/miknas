<template>
  <LayoutHeader>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="所有工作区" :to="extsObj.routePath('')"></PageMenuItem>
        <PageMenuItem title="所有操作" :to="extsObj.routePath('alljobs')"></PageMenuItem>
        <q-btn-dropdown stretch flat label="所有插件">
          <q-list separator bordered>
            <q-item v-for="defInfo in workStore.pluginSimples" :key="defInfo.Id">
              <q-item-section>
                <q-item-label>
                  {{ defInfo.Id }}
                  <q-badge color="grey" :label="defInfo.Version" />
                </q-item-label>
                <q-item-label caption>{{ defInfo.Desc }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </q-tabs>
    </template>
  </LayoutHeader>
</template>

<script setup>
import { LayoutHeader, PageMenuItem } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { useWorkStore } from '../stores/work';
import { onMounted } from 'vue';

const workStore = useWorkStore();
const extsObj = useExtension();

onMounted(() => {
  workStore.refreshMainInfo();
});
</script>
