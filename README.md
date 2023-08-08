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

**Step 1. 编辑生成 docker-compose.yaml**

```yaml
version: "3.9"

services:
  miknas:
    image: mikzone/miknas
    user: "1000:1000"
    ports:
      - "2020:2020"
    volumes:
      - ./workspace:/web/workspace
      - ./config:/web/config
    restart: always
    environment:
      MIKNAS_ADMIN_UID: admin
```

通过该docker-compose.yaml文件，你可以修改：
1. 修改 volumes 下的 **./workspace**，**./config** 把你本地的工作目录、配置目录映射到容器中去。
  （config顾名思义就是配置目录，存放配置相关的，默认也会把数据库文件放在这里。workspace是MikNas的工作目录，虽然MikNas管理的内容很多，为了安全考虑，它只会处理Workspace下的相关文件）
2. 修改 **MIKNAS_ADMIN_UID** 环境变量，改成你希望的管理员账号
3. 修改 user 来指定运行的uid、gid


**Step 2. 用ip和端口打开miknas页面，注册一个管理员账号**

MikNas默认使用admin作为管理账号，因此需要你在登录之前先新建一个admin账号才能开始使用。(你可以修改环境变量MIKNAS_ADMIN_UID来改变管理员账号)
