# Go 微信机器人

⚠️本项目只是作者的一个玩具，大部分功能只为个人定制，并不通用。

## 功能演示

[//]: # ([查看 gif 演示]&#40;https://cdn.xiaobinqt.cn/%E6%BC%94%E7%A4%BA.gif&#41;)

<div align="center"><a href="https://cdn.xiaobinqt.cn/xiaobinqt.io/20230103/f5978e2b41cb46fd856bb988bc78cfee.jpg?imageView2/0/q/75|watermark/2/text/eGlhb2JpbnF0/font/dmlqYXlh/fontsize/1000/fill/IzVDNUI1Qg==/dissolve/52/gravity/SouthEast/dx/15/dy/15" target="_blank"><img src="https://cdn.xiaobinqt.cn/xiaobinqt.io/20230103/f5978e2b41cb46fd856bb988bc78cfee.jpg?imageView2/0/q/75|watermark/2/text/eGlhb2JpbnF0/font/dmlqYXlh/fontsize/1000/fill/IzVDNUI1Qg==/dissolve/52/gravity/SouthEast/dx/15/dy/15" width=  /></a></div>

## 部署说明

clone 项目到本地，然后进入项目目录，将 `config/dev.yaml` 文件改成 `config/prod.yaml`， yaml 配置文件需要配置下，可以去对应的网站获取 apiKey。

执行如下命令：


```shell
go mod tidy # 下载依赖

go build -v -o wxbot  # 编译

nohup ./wxbot > core.log & # 后台运行, 可以查看日志 core.log
```

`less core.log` 可以查看日志，日志里有二维码，可以扫码登录。

## 功能列表

### 定时给女朋友推消息

每天早上 9:30 给女朋友推送一条早安消息，每天晚上 23:00 给女朋友推送一条晚安消息。好吧，我要被女朋友锤了:cry:。

### 自定义事件提醒

生活中的很多事情都是通过微信提醒的，比如快递消息等。这个自定义事件消息，可以通过固定规则，让机器人定时给我们发消息提醒我们某事。

比如`+s15:32,消息内容,3,60`就会定时「今天 15:31 提醒我「消息内容」,提醒 3 次每次间隔 60s」

又如`+st20221227 15:35,,消息内容,3,60`会定时「20221227日 15:35 提醒我「消息内容」,提醒 3 次每次间隔 60s」。具体可以参看功能演示部分。

### 定时给群推送消息

**现在只能通过群名获取群信息**，每天定时给群推送上班打卡等消息，比如每天提醒吃饭。

### 根据关键字回复

基于 [天行](https://www.tianapi.com/) api 和 [和风天气](https://console.qweather.com/#/console?lang=zh) 查询接口开发。

比如在群里发送【泾县天气】机器人会回复泾县今日的天气情况。

现在支持的关键字查询如下：

```
天气查询，如：泾县天气。
菜谱查询，如: 红烧肉菜谱，红烧肉做法。
输入【打赏】打赏卫小兵。
输入【程序员鼓励师】收到程序员鼓励师的回复。
输入【毒鸡汤】关键字回复毒鸡汤。
输入【事件提醒】获取设置事件提醒的格式。
输入【圣诞帽】关键字回复简单处理后的圣诞帽头像，个别用户获取不到头像信息。
输入【英语一句话】关键字回复一句学习英语。
```

## 联系方式

![](https://cdn.xiaobinqt.cn/xiaobinqt.io/20220319/d5616bfc809a45608437f9cc94b14044.jpg?imageView2/0/interlace/1/q/50|imageslim)

## 感谢

基于 [openwechat](https://github.com/eatmoreapple/openwechat) 开发，感谢作者。


