import { createApp } from 'vue'

import App from './App.vue'

import { Quasar, Notify, Dialog, LocalStorage } from 'quasar'
import quasarLang from 'quasar/lang/zh-CN'

// Import icon libraries
import '@quasar/extras/roboto-font/roboto-font.css'
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

const myApp = createApp(App)

myApp.use(Quasar, {
  plugins: {
    Notify,
    Dialog,
    LocalStorage
  }, // import Quasar plugins and add here
  lang: quasarLang
})

import { router, boot, scanAllExtension, gPinia } from 'miknas/utils';

myApp.use(gPinia)

async function run() {
  let extsObjMap = {};
  scanAllExtension(extsObjMap, [
    await import('miknas/exts/Official/extMain'),
    await import('miknas/exts/MikAuth/extMain'),
    await import('miknas/exts/Drive/extMain'),
    await import('miknas/exts/Pan/extMain'),
    await import('miknas/exts/Note/extMain'),
    await import('miknas/exts/SecretShare/extMain'),
    await import('miknas/exts/BookMarks/extMain'),
    await import('miknas/exts/CmdExec/extMain'),
    await import('miknas/exts/DevTool/extMain'),
  ]);
  await boot({
    app: myApp,
    router: router,
    extsObjMap: extsObjMap,
  })

  myApp.use(router)

  // Assumes you have a <div id="app"></div> in your index.html
  myApp.mount('#app')
}

run();
