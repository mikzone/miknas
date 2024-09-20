import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import { router, boot } from 'miknas/utils'

import { Quasar, Notify, Dialog, LocalStorage } from 'quasar'
import quasarLang from 'quasar/lang/zh-CN'

// Import icon libraries
import '@quasar/extras/roboto-font/roboto-font.css'
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

const myApp = createApp(App)

myApp.use(createPinia())

myApp.use(Quasar, {
  plugins: {
    Notify,
    Dialog,
    LocalStorage
  }, // import Quasar plugins and add here
  lang: quasarLang
})

async function run() {
  await boot({
    app: myApp,
    router: router
  })

  myApp.use(router)

  // Assumes you have a <div id="app"></div> in your index.html
  myApp.mount('#app')
}

run();