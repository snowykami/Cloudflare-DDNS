语言Language [zh_CN](README.md) | [EN](README_EN.md)
## Cloudflare DDNS
Cloudflare-DDNS 是一个动态域名解析服务的程序，使用Cloudflare的API进行域名解析到服务器IP记录的更新，基于第4版接口。
- 使用Go编写，跨平台支持，资源占用极低，甚至可以在路由器，树莓派等嵌入式设备上运行
- 支持IPv4和IPv6
- 有轻量级的动态链接库插件系统，可以随意扩展功能
- 适用于HomeLab和家庭网络环境
- 支持插件系统，因此不局限于Cloudflare，你可自行编写插件实现其他服务商接入

### 使用前须知
1. 首先确保你要有公网IP

### 部署

1. 从[Actions](https://github.com/snowykami/Cloudflare-DDNS/actions/)下载你设备对应系统和构架的程序，运行程序
2. 从源码自行构建，需要安装Go环境，然后执行```go build -o ddns main.go```，然后运行```./ddns```

### 如何使用
1. 将域名添加到你的Cloudflare账户中
2. 第一次运行会自动生成```config.yml```，进行修改后再次运行

### 配置文件说明

```yaml
api_key: xxxxxxxx # 从(https://dash.cloudflare.com/profile/api-tokens)获取的apikey
api_email: user@example.com # 你的Cloudflare账户邮箱
zone_id: xxxxxxxx # 从Cloudflare域名主页获取的区域id
duration: 5 # 更新间隔，单位为秒
ddns: # 你的域名记录，此处可以添加多条记录，一个子域名可以同时存在A和AAAA记录
  - name: v4.example.com
    type: A
    ttl: 60
    proxied: false
    comment: "DDNS auto update"

  - name: v6.example.com
    type: AAAA
    ttl: 60
    proxied: false
    comment: "DDNS auto update"
... 其他插件配置
```

### 一些建议
1. 考虑到用DDNS的用户大部分是有公网的家庭宽带用户，v4和v6地址不一定在一个主机上，所以建议在本地配置端口转发，然后A和AAAA使用不同的记录（一个域名可以同时有A和AAAA记录）
2. 可以给设置自启动服务，方便服务器断电重启拨号后能自动更新域名解析
3. 如果你的主机是在NAT环境下，可以使用UPnP或者NAT-PMP协议进行端口映射
4. 如果你的网络设备不支持IPv6，那就不要配置AAAA记录

### 插件系统
DDNS支持一个轻量的插件系统，可以实现一些功能例如发送邮件提醒或在IP变动时重启DDNS
- 插件仅支持Linux/macOS平台
- 可参考[插件示例](./plugins/example.go)进行编写
- 编写完插件后使用build_plugins.sh进行编译，可接文件名指定插件，不指定则编译全部插件
- 插件可通过导入config包来获取全局配置

### 预制插件(部分需要单独编译)
- 预编译的插件在[prebuild](./prebuild/)目录下
- OneBot通知：在IP变动时通过OneBot发送通知
  配置项
  ```yaml
  onebot_instance: http://127.0.0.1:8080
  onebot_token: xxxxxxxx
  onebot_user: 12345678
  ```
- NginxManager：在IP变动时重启正代服的Nginx
  配置项
  ```yaml
  ssh_host: xxx.xxxx.xxxx
  ssh_port: "22"
  ssh_user: root
  ssh_key_path: /root/.ssh/id_rsa
  # 暂时不支持密码登录，有需要的自己改
  ```

### TODO

- [ ] Lua脚本插件支持(由[synodriver](https://github.com/synodriver)提出)