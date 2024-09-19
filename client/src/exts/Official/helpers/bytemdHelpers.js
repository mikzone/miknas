import {
  FileUtil,
  openSelectFsidFileDlg,
} from 'miknas/exts/Drive/shares';
import { MikCall } from 'miknas/utils';

async function tryInjectFile(ctx, isImg){
  let [isOk, result] = await openSelectFsidFileDlg(
    'Priv',
    '',
    '选择图片'
  );
  if (!isOk) return;
  let filePath = result.fspath;
  if (!filePath) {
    MikCall.sendErrorTips('选择的图片为空，不做处理');
    return;
  }
  let fileViewProxy = result.fileViewProxy;
  let viewUrl = fileViewProxy.fileOp.getViewUrl(filePath);
  let fileName = FileUtil.baseName(filePath);
  let txt = `[${fileName}](${viewUrl})\n\n`;
  if (isImg) {
    txt = '!' + txt;
  }
  ctx.appendBlock(txt);

  ctx.editor.focus();
}

export function drivePlugin() {
  return {
    actions: [
      {
        // title: '插入图片',
        icon: '<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" fill="none" viewBox="0 0 48 48"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M5 10a2 2 0 0 1 2-2h34a2 2 0 0 1 2 2v28a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V10Z" clip-rule="evenodd"/><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M14.5 18a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3Z" clip-rule="evenodd"/><path stroke="currentColor" stroke-linejoin="round" stroke-width="4" d="m15 24 5 4 6-7 17 13v4a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-4l10-10Z"/></svg>', // 16x16 SVG icon
        handler: {
          type: 'dropdown',
          actions: [
            {
              title: '嵌入图片',
              handler: {
                type: 'action',
                async click(ctx) {
                  await tryInjectFile(ctx, true);
                },
              },
            },
            {
              title: '嵌入文件',
              handler: {
                type: 'action',
                async click(ctx) {
                  await tryInjectFile(ctx, false);
                },
              },
            },
          ]
        },
      },
    ],
  };
}
