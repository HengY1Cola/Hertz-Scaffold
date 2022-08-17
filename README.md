## 为什么想写一个

字节搞了个开源的Hertz框架，然后公司的业务分层也是值得学习的。

日常中使用的是Gin框架，不过现在习惯了公司的思路就想着用开源的搭建个脚手架。

我按照gin_scaffold的思路用Hertz改了个脚手架，从Gin过来的成本很低。

本框架能实现Http相关的业务上的『快速开始』。

## 如何快速开始
> 前置条件：已经配置好了Git与Go相关环境
```bash
$ git clone https://github.com/HengY1Sky/Hertz-Scaffold
$ cd Hertz-Scaffold
$ chmod +x ./build.sh
$ ./build.sh && ./output/bootstrap_boe.sh # 开始体验
```

## 文件结构

```
├─biz # business
│  ├─bo # request && reponse && Object
│  ├─constant # 定义
│  ├─dal
│  ├─handler # 业务
│  ├─middleware # 中间件 
│  ├─model # 模型
│  ├─repository # 数据库
│  ├─service
│  ├─utils # 仓库
│  └─validate # 验证器
├─conf # 配置文件
└─cron_job # 定时任务
```

## 实现功能

1. 分组路由文件上的封装，按照统一格式写就好了
2. 中间件独立出来了，直接在对应的路由进行注册就好了
3. Binding的自定义校验Validator
4. Success与Error返回的统一格式
5. Logger的追踪与全局的实现
6. 配置文件的读取以及运用

## 相关链接

- Hertz官方文档：https://www.cloudwego.io/zh/docs/hertz/
- Gin脚手架：https://github.com/e421083458/gin_scaffold
