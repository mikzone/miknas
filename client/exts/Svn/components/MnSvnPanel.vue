<template>
  <q-card>
    <q-table
      ref="svnTable"
      v-model:selected="selectedRows"
      :rows="svnStEntrys"
      :columns="tableColumns"
      bordered
      dense
      style="max-height: 700px;"
      :loading="svnStLoading"
      :rows-per-page-options="[0]"
      hide-bottom
      selection="multiple"
      class="mn-sticky-header-table"
      no-data-label="当前svn目录没有任何修改"
      row-key="path"
    >
      <template #top>
        <div class="q-table__title">svn目录({{ svnFolder }})</div>
        <q-space></q-space>
        <q-btn
          class="q-ml-sm"
          color="primary"
          :disable="svnStLoading"
          label="刷新"
          @click="refreshSvnSt()"
        ></q-btn>
      </template>

      <template #loading>
        <q-inner-loading
          showing
          color="primary"
        ></q-inner-loading>
      </template>
    </q-table>
    <q-card-actions>
      <MdcCmdExec ref="svnExecCmdRef">
        <q-btn-group>
          <q-btn
            color="white"
            text-color="black"
            label="svn up"
            @click="execSvn('up', null)"
          ></q-btn>
          <q-btn
            v-show="selectedRows.length"
            color="white"
            text-color="black"
            label="查diff"
            @click="execSvn('diff', null)"
          ></q-btn>
          <q-btn
            v-show="selectedRows.length"
            color="white"
            text-color="black"
            label="svn commit"
            @click="promptSvnCommit"
          ></q-btn>
          <q-btn
            v-show="selectedRows.length"
            color="white"
            text-color="black"
            label="svn revert"
            @click="execSvn('revert', null)"
          ></q-btn>
          <q-btn-dropdown
            color="white"
            text-color="black"
            label="其它"
          >
            <q-list>
              <q-item
                v-show="selectedRows.length"
                v-close-popup
                clickable
                @click="execSvn('add', null)"
              >
                <q-item-section>
                  <q-item-label>svn add</q-item-label>
                </q-item-section>
              </q-item>
              <q-item
                v-show="selectedRows.length"
                v-close-popup
                clickable
                @click="execSvn('deleteForce', null)"
              >
                <q-item-section>
                  <q-item-label>svn delete --force</q-item-label>
                </q-item-section>
              </q-item>
              <q-item
                v-close-popup
                clickable
                @click="execSvn('cleanup', null)"
              >
                <q-item-section>
                  <q-item-label>svn cleanup --include-externals</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </q-btn-group>
      </MdcCmdExec>
    </q-card-actions>
    <q-dialog
      v-model="showCommitDialog"
      persistent
    >
      <q-card style="width: 1200px; max-width: 80vw;">
        <q-card-section class="row items-center">
          <div class="text-h6">请输入提交的message</div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <MdcRecentInput
            ref="commitMessageRef"
            v-model="commitMessage"
            placeholder="请输入提交的message"
            store-key="commitMessage"
          >
          </MdcRecentInput>
        </q-card-section>
        <q-card-actions align="right" class="text-primary">
          <q-btn v-close-popup flat label="取消" />
          <q-btn flat label="提交" @click="tryCommit()" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-card>
</template>

<script>
import { gutil, MikCall } from 'miknas/utils';
import MdcCmdExec from 'miknas/exts/CmdExec/MdcCmdExec';
import MdcRecentInput from 'miknas/exts/Official/components/MdcRecentInput';
import { useExtension } from '../extMain';
let extsObj = useExtension();
export default {
  name: 'MnSvnStPanel',
  components: {
    MdcCmdExec,
    MdcRecentInput,
  },
  props: {
    svnFolder: {
      type: String,
      required: true,
    },
  },
  data: function () {
    return {
      svnStJson: '',
      selectedRows: [],
      svnStLoading: false,
      showCommitDialog: false,
      commitMessage: '',
      tableColumns: [
        {
          name: 'path',
          label: '目录',
          field: 'path',
          sortable: true,
          align: 'left',
        },
        {
          name: 'item',
          label: 'item',
          field: 'item',
          sortable: true,
          align: 'left',
          style: 'width: 100px',
        },
        {
          name: 'props',
          label: 'props',
          field: 'props',
          sortable: true,
          align: 'left',
          style: 'width: 75px',
        },
        {
          name: 'tree-conflicted',
          label: 'tree-conflicted',
          field: 'tree-conflicted',
          sortable: true,
          align: 'left',
          style: 'width: 120px',
        },
        {
          name: 'wc-locked',
          label: 'wc-locked',
          field: 'wc-locked',
          sortable: true,
          align: 'left',
          style: 'width: 120px',
        },
      ],
    };
  },
  computed: {
    svnStEntrys: function () {
      let ret = [];
      if (!this.svnStJson) return ret;
      let curEntrys = this.svnStJson.status.target.entry;
      if (!curEntrys) return ret;
      if (!Array.isArray(curEntrys)) {
        // 如果是一个元素，这个就是单纯的object而已，坑
        curEntrys = [curEntrys];
      }
      for (let obj of curEntrys) {
        let addItem = {
          path: obj['@path'],
          props: obj['wc-status']['@props'],
          item: obj['wc-status']['@item'],
          'tree-conflicted': obj['wc-status']['@tree-conflicted'],
          'wc-locked': obj['wc-status']['@wc-locked'],
        };

        if (addItem.props == 'none') {
          addItem.props = '';
        }
        if (addItem.item == 'normal') {
          addItem.item = '';
        }
        if (addItem['wc-locked'] == 'false') {
          addItem['wc-locked'] = '';
        }

        if (!addItem.item && !addItem.props && !addItem['wc-locked'])
          continue;
        if (addItem.item == 'external' && !addItem.props) continue;
        ret.push(addItem);
      }
      return ret;
    },
  },
  mounted: function () {
    this.refreshSvnSt();
  },
  methods: {
    refreshSvnSt: function () {
      extsObj
        .mpost('querySvnSt', { svnFolder: this.svnFolder })
        .then((cbDict) => {
          // console.log('cbDict', cbDict);
          this.svnStJson = cbDict;
          this.svnStLoading = false;
          this.selectedRows = gutil.refreshTableSelectedRows(this.$refs.svnTable, this.svnStEntrys);
        });
      this.svnStLoading = true;
    },
    execSvn: function (opType, opExt) {
      if (!opType) {
        MikCall.sendErrorTips('opType不能为空');
        return;
      }
      opExt = opExt || {};
      let reqArgs = {
        tp: opType,
      };

      let needFileOps = ['commit', 'diff', 'revert', 'add', 'deleteForce'];
      if (needFileOps.includes(opType)) {
        let needFiles = [];
        for (let obj of this.selectedRows) {
          let fileName = obj.path;
          if (fileName) needFiles.push(fileName);
        }

        if (needFiles.length <= 0) {
          MikCall.sendErrorTips(opType + '操作需要选择具体的文件');
          return;
        }
        reqArgs['files'] = needFiles;
      }

      if (opType == 'commit') {
        let msg = opExt.msg;
        if (!msg) {
          MikCall.sendErrorTips('commit的message不能为空');
          return;
        }
        reqArgs['msg'] = msg;
      }

      this.$refs.svnExecCmdRef.tryExec({
        url: extsObj.serverUrl('execSvnMethod'),
        data: { opDict: JSON.stringify(reqArgs), svnPath: this.svnFolder },
        sucCb: (cbDict) => {
          console.log('execSvnMethod cbDict', cbDict);
          // 非 diff 模式下需要更新状态
          if (!['diff'].includes(reqArgs.tp)) {
            this.refreshSvnSt();
          }
          if (reqArgs.tp == 'commit'){
            this.$refs.commitMessageRef.addRecentInputValue(reqArgs['msg']);
            this.showCommitDialog = false;
          }
          if (reqArgs.tp == 'diff'){
            this.$refs.svnExecCmdRef.setAceLang('diff');
          }
        },
      });
    },
    promptSvnCommit: function () {
      this.showCommitDialog = true;
      // this.$q
      //   .dialog({
      //     title: "请输入提交的message",
      //     prompt: {
      //       model: "",
      //       type: "text", // optional
      //     },
      //     cancel: true,
      //     persistent: true,
      //   })
      //   .onOk((data) => {
      //     this.execSvn("commit", { msg: data });
      //   });
    },
    tryCommit: function(){
      this.execSvn('commit', { msg: this.commitMessage });
      // 不能在下面关了，关了ref就没了，会导致更新不了消息记录
      // this.showCommitDialog = false;
    },
  },
};
</script>
