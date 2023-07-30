# MiknasClient (miknas-client)

client for miknas

## Install the dependencies
```bash
yarn
```

### Start the app in development mode (hot-code reloading, error reporting, etc.)
```bash
quasar dev
```


### Lint the files
```bash
yarn lint
```


### Format the files
```bash
yarn format
```



### Build the app for production
```bash
quasar build
```

### Customize the configuration
See [Configuring quasar.config.js](https://v2.quasar.dev/quasar-cli-vite/quasar-config-js).


### 备注
* 开发过程中需要配置 quasar.conf.js 的 devServer.proxy 的内容指向对应的后端地址

### 插件开发规范
1. 插件统一放在 miknas/exts 下，以插件名为插件目录名，必须以大写字母开头，驼峰命名法
2. 每个插件目录下必须含有一个 extMain.js 用于配置插件的相关信息。
3. 插件目录下：
  - shares.js 一般用来提供给外部公用的代码块
  - components 组件。
  - pages 存放该插件自定义的展示页
  - extMain.js 要使用 defineExtension 定义好插件的
    - id: 插件id，和目录名保持一致。
    - name: 插件名称
    - desc: 插件的介绍
    - route(可选)：路由定义
    - boot(可选): boot函数
4. 每个插件的路由首页都是 自己路由下面的根路径
