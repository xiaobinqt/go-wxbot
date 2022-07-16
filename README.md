# golang 微信机器人

## ~~阵亡了，2021-12-31 微信关闭了 uos 微信登录，绕过网页版登录限制已经不再生效。~~

基于 [openwechat](https://github.com/eatmoreapple/openwechat) 开发，感谢作者。

## 部署说明

clone 项目到本地，然后进入项目目录，将 `config/dev.yaml` 文件改成 `config/prod.yaml`，
yaml 配置文件需要配置下，可以去对应的网站获取 apiKey。

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

### 定时给群推送消息

[openwechat](https://github.com/eatmoreapple/openwechat) 现在只能通过群名获取群信息，每天定时给群推送上班打卡消息。

### 根据关键字回复

基于 [天行](https://www.tianapi.com/) api 和 [高德天气](https://lbs.amap.com/api/webservice/guide/api/weatherinfo/) 查询接口开发。

比如在群里发送【泾县天气】机器人会回复泾县今日的天气情况。

现在支持的关键字查询如下：

```
天气查询，如：泾县天气。
菜谱查询，如: 红烧肉菜谱，红烧肉做法。
输入【打赏】打赏卫小兵。
输入【程序员鼓励师】收到程序员鼓励师的回复。
输入【毒鸡汤】关键字回复毒鸡汤。
输入【圣诞帽】关键字回复简单处理后的圣诞帽头像，个别用户获取不到头像信息。
输入【英语一句话】关键字回复一句学习英语。
```




