<template>
  <q-page>
    <div class="q-pa-md row q-gutter-sm">
      <q-btn
        color="primary"
        label="添加书签"
        @click="openAddBookmark(allKinds, onAddBookmark)"
      ></q-btn>
      <q-btn
        color="primary"
        label="刷新"
        @click="refreshResult"
      ></q-btn>
    </div>
    <div v-for="(kindInfos, kind) in filterResult" :key="kind">
      <q-separator />
      <div class="q-pa-md text-h5">{{ kind }}</div>
      <div class="q-pa-md row q-col-gutter-sm">
        <div
          v-for="urlInfo in kindInfos"
          :key="urlInfo.id"
          class="col-12 col-sm-6 col-md-4"
        >
          <q-card class="no-shadow mn-bordered link-card">
            <q-card-section horizontal>
              <q-item tag="a" :href="urlInfo.url" target="_blank" class="col">
                <q-item-section avatar>
                  <q-icon v-if="urlInfo.iconError" name="public" color="red" />
                  <q-avatar v-else-if="urlInfo.icon" size="md">
                    <img
                      :src="urlInfo.icon"
                      @error="imgerrorfun(urlInfo)"
                    />
                  </q-avatar>
                  <q-icon v-else name="public" />
                </q-item-section>
                <q-item-section>
                  <q-item-label :lines="1">{{ urlInfo.name }}</q-item-label>
                  <q-item-label :lines="1" caption>{{
                    urlInfo.url
                  }}</q-item-label>
                </q-item-section>
              </q-item>

              <q-separator vertical />

              <q-btn flat icon="more_vert" @click.stop.prevent="">
                <q-menu auto-close>
                  <q-list>
                    <q-item
                      v-close-popup
                      clickable
                      @click="
                        openModifyBookmark(urlInfo, allKinds, refreshResult)
                      "
                    >
                      <q-item-section>
                        <q-item-label>修改</q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item
                      v-close-popup
                      clickable
                      @click="deleteBookmark(urlInfo, refreshResult)"
                    >
                      <q-item-section>
                        <q-item-label>删除</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-btn>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
    <q-inner-loading
      :showing="loadingMgr.isloading.value"
      :label="loadingMgr.loadingLabel.value"
    />
  </q-page>
</template>

<script setup>
import { onMounted, reactive } from 'vue';
import { useMikLoading } from 'miknas/exts/Official/shares';
import { useExtension } from '../extMain';
import { computed } from 'vue';
import {
  deleteBookmark,
  openAddBookmark,
  openModifyBookmark,
} from '../cHelpers';
import { MikCall } from 'miknas/utils';

const loadingMgr = useMikLoading();
const state = reactive({
  result: [],
});

const filterResult = computed(() => {
  let ret = {};
  for (let info of state.result) {
    if (ret[info.kind] == undefined) {
      ret[info.kind] = [];
    }
    ret[info.kind].push(info);
  }
  return ret;
});

const allKinds = computed(() => {
  let ret = {};
  for (let info of state.result) {
    ret[info.kind] = 1;
  }
  return Object.keys(ret);
});

function convertBookmarkInfo(bookmarkInfo) {
  return bookmarkInfo;
}

function onAddBookmark(bookmarkInfo) {
  state.result.push(convertBookmarkInfo(bookmarkInfo));
}

const extsObj = useExtension();

async function refreshResult() {
  let stateName = `正在加载`;
  loadingMgr.addLoadingState(stateName);
  let iRet = await extsObj.mcpost('getall');
  if (!iRet.suc) {
    MikCall.alertRespErrMsg(iRet);
    loadingMgr.removeLoadingState(stateName);
    return;
  }
  let result = iRet.ret;
  state.result = result.map((x) => convertBookmarkInfo(x));
  loadingMgr.removeLoadingState(stateName);
}

onMounted(() => {
  refreshResult();
});

function imgerrorfun(urlInfo) {
  urlInfo.iconError = true;
}
</script>

<style scoped>
.link-card .text-caption {
  line-height: 1.25em !important;
}
</style>
