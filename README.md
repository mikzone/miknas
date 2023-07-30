## MikNas (自建Nas管理工具)

### 简介

MikNas 可以认为是一个集合若干功能的工具箱，它由一个个扩展构成，每个扩展用于解决某一类功能。

目前已有的功能有:
	1. 个人网盘(文件管理)
	2. 用户管理
	3. 书签
	4. 密文分享
	5. 笔记

Ps: 近来，入手了个小主机，装了个linux系统，当Nas服务器来使用。网上其实有很多工具，我也装了不少，但或多或少，有些点不太满足于我，于是萌生了自己做一个的想法，当然也只限于简单的需求。

### 使用方法

#### 推荐使用docker来部署

**Step 1. 先编辑 docker-compose.yaml**，可以根据下面内容修改

```yaml
version: "3.9"

services:
  miknas:
    image: mikzone/miknas
    ports:
      - "2020:2020"
    volumes:
      - ./web:/web
    restart: always
    environment:
      MIKNAS_SECRET_KEY: MnszoodXuiCcPhTL
      GIN_MODE: release
```

通过该docker-compose.yaml文件，你可以修改：
1. 修改 volumes，把你本地的目录(需要确保目录已经预先生成,并且里面有workspace文件夹)映射到容器的/web去，运行会生成miknas.sqlite数据库文件。
2. 修改MIKNAS_SECRET_KEY环境变量，改成随机的其它值，以防止session被破解

**Step 2. 用ip和端口打开miknas页面，注册一个admin账号**

MikNas默认使用admin作为管理账号(你可以修改环境变量来改变)，因此需要你新建一个admin账号才能开始使用。
