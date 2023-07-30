<template>
  <ExtensionPage>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="首页" :to="extsObj.routePath('')"></PageMenuItem>
      </q-tabs>
    </template>
    <router-view />
  </ExtensionPage>
</template>

<script setup>
import { ExtensionPage, PageMenuItem } from 'miknas/exts/Official/shares';

import { useExtension } from '../extMain';
const extsObj = useExtension();
</script>
