### 1.构建基础server，监听连接
> window build 需要加exe，否则无法执行
> 测试tcp连接，可以直接用浏览器访问

### 2.用户上线功能
> 用wsl 的nc 命令时，经测试127.0.0.1，0.0.0.0这样的地址都不能正常访问到，
> 使用本机ip地址可以访问
> 也可以使用cmd命令的curl连接
新建OnlineMap，记录在线用户
监听用户发来的消息，并广播给在线用户消息

### 3.用户消息广播
### 4.业务层封装
### 5.修改用户名
### 6.用户超时下线
### 7.私聊
### 8.实现客户端client
### 9.设置命令行参数并解析参数
### 10.展示菜单
### 11.修改用户用户名
### 12.群聊模式
### 13.私聊模式
> 1. 查询在线用户
> 2. 选择在线用户