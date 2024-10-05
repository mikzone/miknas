import { defineExtension } from 'miknas/utils';

const EXTS_ID = 'Svn';

export const useExtension = defineExtension({
  id: EXTS_ID,
  title: 'SVN 面板',
  desc: '管理SVN',
  icon: 'source',
  index: false,
});
