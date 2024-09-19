<template>
  <ExtensionPage>
    <template #login-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="首页" :to="extsObj.routePath('')"></PageMenuItem>
        <PageMenuItem v-if="extsObj.hasAuth('manager')" title="用户管理" :to="extsObj.routePath('manage')"></PageMenuItem>
        <PageMenuItem v-if="extsObj.hasAuth('manager')" title="权限管理" :to="extsObj.routePath('roles')"></PageMenuItem>
      </q-tabs>
    </template>
    <template #unlogin-toolbar>
      <q-tabs shrink stretch>
        <PageMenuItem title="登录" :to="extsObj.routePath('login')"></PageMenuItem>
        <PageMenuItem title="注册" :to="extsObj.routePath('register')"></PageMenuItem>
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
