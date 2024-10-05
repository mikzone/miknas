<template>
  <ExtensionPage>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="Workspace" :to="extsObj.routePath('ws')"></PageMenuItem>
        <PageMenuItem title="数据库" :to="extsObj.routePath('db')"></PageMenuItem>
        <PageSubMenu title="其它">
          <PageMenuItem in-sub-menu title="下载数据库" target="_blank" :href="extsObj.serverUrl('miknas.db')"></PageMenuItem>
        </PageSubMenu>
      </q-tabs>
    </template>
    <router-view />
  </ExtensionPage>
</template>

<script setup>
import { ExtensionPage, PageMenuItem, PageSubMenu } from 'miknas/exts/Official/shares';

import { useExtension } from '../extMain';
const extsObj = useExtension();
</script>
