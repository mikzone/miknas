<template>
  <div>
    <div class="row items-center justify-end q-pa-md">
      <q-breadcrumbs gutter="xs" class="text-wrap">
        <q-breadcrumbs-el
          v-for="oneBC in fileGetter.fileBC.value"
          :key="oneBC.name"
          :label="oneBC.name"
          :icon="oneBC.icon"
          :class="oneBC.path === undefined ? '' : 'text-primary cursor-pointer'"
          @click="fileOp.gotoPath(oneBC.path)"
        />
      </q-breadcrumbs>
      <q-space />
      <div>
        <q-btn
          dense
          flat
          round
          color="primary"
          :icon="fileState.showModeConfs[fileGetter.showMode.value].icon"
        >
          <q-menu auto-close>
            <q-list dense>
              <q-item-label header class="q-pa-sm">显示模式</q-item-label>
              <q-item
                v-for="(modeConf, mode) in fileState.showModeConfs"
                :key="mode"
                dense
                clickable
                :active="fileGetter.showMode.value == mode"
                @click="fileOp.updateShowMode(mode)"
              >
                <q-item-section avatar class="mn-item-icon">
                  <q-icon :name="modeConf.icon" />
                </q-item-section>
                <q-item-section>{{ modeConf.desc }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>

        <q-btn dense flat round color="primary" icon="more_vert">
          <q-menu>
            <q-list dense>
              <q-item-label header class="q-pa-sm">排序</q-item-label>
              <q-item
                v-for="(sortConf, sortMethod) in fileGetter.sortConfs.value"
                :key="sortMethod"
                v-ripple
                clickable
                :active="fileGetter.curSortSave.value.method == sortMethod"
                @click="fileOp.updateSort(sortMethod)"
              >
                <q-item-section>
                  {{ sortConf.label }}
                </q-item-section>
                <q-item-section side>
                  <q-icon
                    v-show="fileGetter.curSortSave.value.method == sortMethod"
                    :name="
                      (fileGetter.curSortSave.value.desc && 'arrow_downward') || 'arrow_upward'
                    "
                  />
                </q-item-section>
              </q-item>
              <q-separator spaced />
              <q-item v-ripple clickable>
                <q-item-section> 查看文件夹大小 </q-item-section>
                <q-item-section side>
                  <q-checkbox
                    :model-value="!!fileGetter.curRecData.value.hasFolderSize"
                    @update:model-value="updateNeedFolderSize"
                  />
                </q-item-section>
              </q-item>
              <q-separator spaced />
              <q-item-label header class="q-pa-sm">其它</q-item-label>
              <q-item
                v-show="!fileGetter.isloading.value"
                v-close-popup
                clickable
                @click="fileOp.tryRefreshFiles()"
              >
                <q-item-section avatar class="mn-item-icon">
                  <q-icon name="refresh" />
                </q-item-section>
                <q-item-section>刷新当前目录</q-item-section>
              </q-item>
              <template v-if="!fileGetter.isReadOnly.value">
                <q-item v-close-popup clickable @click="fileOp.newFolder">
                  <q-item-section avatar class="mn-item-icon">
                    <q-icon name="create_new_folder" />
                  </q-item-section>
                  <q-item-section>新建目录</q-item-section>
                </q-item>
                <q-item v-close-popup clickable @click="fileOp.openUploadDlg">
                  <q-item-section avatar class="mn-item-icon">
                    <q-icon name="cloud" />
                  </q-item-section>
                  <q-item-section>上传文件到此处</q-item-section>
                </q-item>
              </template>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </div>
    <q-list v-if="!fileGetter.grid.value" bordered separator>
      <q-item
        v-for="fileInfo in fileState.curFiles"
        :key="fileInfo.name"
        v-ripple
        clickable
        @click="fileOp.clickOpen(fileInfo)"
      >
        <q-item-section avatar>
          <q-icon :name="fileInfo.icon" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="mn-word-break-all" :lines="2">{{ fileInfo.name }}</q-item-label>
          <q-item-label caption>
            <template v-if="fileInfo.viewSize"> {{ fileInfo.viewSize }} | </template>
            {{ fileInfo.viewModify }}
          </q-item-label>
        </q-item-section>
        <q-item-section side>
          <mdc-file-more-action
            v-if="!fileState.isSelectMode"
            :file-view-proxy="fileViewProxy"
            :file-info="fileInfo"
          />
          <q-checkbox
            v-else
            v-model="fileState.curSelected"
            :val="fileInfo.name"
            color="secondary"
          />
        </q-item-section>
      </q-item>
    </q-list>
    <div v-else class="q-pa-sm">
      <div class="q-col-gutter-sm row items-start">
        <div
          v-for="fileInfo in fileState.curFiles"
          :key="fileInfo.name"
          :class="fileGetter.gridClass.value"
        >
          <q-card class="no-shadow mn-bordered">
            <q-img
              v-if="fileInfo.fileType == 'img'"
              :src="calcThumbSrc(fileInfo.name)"
              fit="cover"
              loading="lazy"
              :ratio="fileGetter.imgRatio.value"
              class="cursor-pointer"
              @click="fileOp.clickOpen(fileInfo)"
            />
            <q-img
              v-else-if="fileInfo.fileType == 'video2'"
              :src="calcThumbSrc(fileInfo.name)"
              fit="cover"
              loading="lazy"
              :ratio="fileGetter.imgRatio.value"
              class="cursor-pointer"
              @click="fileOp.clickOpen(fileInfo)"
            />
            <q-img
              v-else-if="fileGetter.showMode.value == 'grid'"
              :src="fileOp.getTypeSvgData(fileInfo.fileType)"
              fit="cover"
              :ratio="fileGetter.imgRatio.value"
              class="cursor-pointer"
              @click="fileOp.clickOpen(fileInfo)"
            />
            <q-card-section v-if="fileGetter.showMode.value == 'grid'" class="q-pa-xs">
              <q-item-label class="mn-word-break-all text-center" lines="2">{{
                fileInfo.name
              }}</q-item-label>
              <div class="row items-center no-wrap">
                <div class="col">
                  <div class="text-caption">{{ fileInfo.viewSize }}</div>
                </div>

                <div class="col-auto">
                  <mdc-file-more-action
                    v-if="!fileState.isSelectMode"
                    :file-view-proxy="fileViewProxy"
                    :file-info="fileInfo"
                  />
                  <q-checkbox
                    v-else
                    v-model="fileState.curSelected"
                    :val="fileInfo.name"
                    color="secondary"
                  />
                </div>
              </div>
            </q-card-section>
            <q-item v-else v-ripple clickable @click="fileOp.clickOpen(fileInfo)">
              <q-item-section avatar>
                <q-icon :name="fileInfo.icon" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="mn-word-break-all" :lines="2">{{
                  fileInfo.name
                }}</q-item-label>
                <q-item-label caption>
                  <template v-if="fileInfo.viewSize"> {{ fileInfo.viewSize }} | </template>
                  {{ fileInfo.viewModify }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <mdc-file-more-action
                  v-if="!fileState.isSelectMode"
                  :file-view-proxy="fileViewProxy"
                  :file-info="fileInfo"
                />
                <q-checkbox
                  v-else
                  v-model="fileState.curSelected"
                  :val="fileInfo.name"
                  color="secondary"
                />
              </q-item-section>
            </q-item>
          </q-card>
        </div>
      </div>
    </div>
    <div v-if="fileState.curFiles.length <= 0" class="q-pa-md">文件列表为空</div>
    <div class="fixed-bottom row flex-center" style="pointer-events: none">
      <div v-if="fileState.isSelectMode" style="pointer-events: auto">
        <q-btn-group class="bg-secondary">
          <q-btn
            dense
            :icon="fileGetter.selectAllStatusIcon.value"
            text-color="white"
            @click.stop.prevent="fileOp.toggleSelectAll"
          />
          <q-btn-dropdown auto-close color="secondary" label="操作">
            <!-- dropdown content goes here -->
            <q-list dense class="bg-secondary text-white">
              <q-item clickable @click.stop.prevent="fileOp.multThumb">
                <q-item-section> 生成缩略图 </q-item-section>
              </q-item>
              <q-item clickable @click.stop.prevent="fileOp.multMove">
                <q-item-section> 移动 </q-item-section>
              </q-item>
              <q-item clickable @click.stop.prevent="fileOp.multCopy">
                <q-item-section> 复制 </q-item-section>
              </q-item>
              <q-item clickable @click.stop.prevent="fileOp.multRemove">
                <q-item-section> 删除 </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
          <q-btn dense color="secondary" icon="close" @click.stop.prevent="fileOp.quitSelectMode" />
        </q-btn-group>
      </div>
    </div>
    <q-inner-loading
      :showing="fileGetter.isloading.value"
      color="primary"
      :label="fileGetter.loadingLabel.value"
    />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useExtension } from '../extMain';
import { useFileView } from 'miknas/exts/Drive/shares';
import MdcFileMoreAction from './subcoms/MdcFileMoreAction.vue';

const extsObj = useExtension();

const props = defineProps({
  fsid: {
    type: String,
    required: true
  },
  rootPath: {
    type: String,
    default: ''
  },
  extraConf: {
    type: Object,
    default: () => {
      return {};
    }
  }
});

const fileViewProxy = useFileView(props.fsid, props.rootPath, props.extraConf);

function calcThumbSrc(fileName) {
  let fpath = fileOp.abs(fileName);
  if (fileGetter.showMode.value == 'blog') {
    return fileOp.getViewUrl(fpath);
  }
  return extsObj.serverUrl(`thumb/${props.fsid}/${fpath}`);
}

function updateNeedFolderSize(val) {
  fileOp.updateNeedFolderSize(val);
}

const { fileState, fileGetter, fileOp } = fileViewProxy;

onMounted(() => {
  // console.log('onmounted fileviewmgr');
});
</script>
