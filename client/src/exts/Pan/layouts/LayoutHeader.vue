<template>
  <LayoutHeader>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="首页" :to="extsObj.routePath('list')"></PageMenuItem>
        <PageMenuItem title="我的分享" :to="extsObj.routePath('shares')"></PageMenuItem>
      </q-tabs>
    </template>
  </LayoutHeader>
</template>

<script setup>
import { LayoutHeader, PageMenuItem } from 'miknas/exts/Official/shares';

import { useExtension } from '../extMain';
const extsObj = useExtension();
</script>
