# 星际竞技场

## 快速上手

- 安装SC2游戏客户端
- 安装Go语言编译器1.18版本或以上
- 访问https://goproxy.io/，按照提示配置GOPROXY
- 在项目根目录下执行命令`go mod tidy`，下载依赖
- 编辑项目目录下的`conf/conf.toml`文件，修改`bilidanmu.room_id`为你的直播间ID，如果不需要测试弹幕，可以将`bilidanmu.enable`设置为`false`
- 将项目目录下的`sc2maps/product/StarArena.SC2Map`文件复制到SC2游戏客户端安装目录下的`Maps/Custom`路径下（如：`D:\Blizzard\StarCraft II\Maps\Custom`，如果没有此路径可以手动创建）
- 在项目目录下执行命令`go run .`启动程序测试

## 地图开发

- 编辑地图
打开SC2地图编辑器（游戏安装目录下的`StarCraft II Editor_x64.exe`)，点击左上角的打开按钮选择项目根目录下的`sc2maps/documents/StarArena.SC2Map`文件夹，即可开始地图编辑。

- 编写触发器
打开`VSCode`编辑器，在插件商店中搜索安装`talv.sc2galaxy`和`talv.sc2layouts`这两个插件，然后用`VSCode`打开`sc2maps/documents/StarArena.SC2Map`文件夹，可以在`CustomScripts`目录中找到触发器脚本。

- 导出地图文件测试
在SC2地图编辑器中将地图另存为`SC2Map`文件格式到`sc2maps/product/StarArena.SC2Map`，并复制一份到SC2游戏客户端安装目录下的`Maps/Custom`目录中，即可开始测试。

## 实现机制

服务启动了两个SC2客户端，相当于两个真人玩家，一个称为`导演(director)`，另一个为`观众(audience)`。  
`director`用于接收外部程序发送的转换后的弹幕指令，具备完整的游戏控制界面，也可用于测试和活动运营。在SC2地图中对应的玩家ID为`1`。  
`audience`用于展现直播画面，无控制界面，可以显示消息。在SC2地图中对应的玩家ID为`2`。  
SC2地图中弹幕游戏玩家所属的红蓝队阵营为`敌对(hostile)`玩家（非`电脑(computer)`玩家，SC2机器人接口只支持 一个真人玩家与多个电脑玩家 或 两个真人玩家）。