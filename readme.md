

## 为什么会有这个项目

在实际工作中，我司生产的采集终端设备（linux端），需要依据实际项目要求配置json文件内的参数。但是，现场实施的同事不具备linux基础知识（常见指令，vi等），需要我这边远程操作。大部分现场实施环境处于人迹罕至区域，不具备较好的网络环境（设备4G上网），远程操作受到极大限制。同时，出厂设置时因为生产同事也不具备linux操作能力，我又特别懒不愿意做重复劳动，所以打算做个界面让同事轮椅式操作出厂配置。一开始考虑让低级帕鲁写一个配置网页，但是我考虑到以后设备可能会使用低端产品（64M内存） ，再支持一个node环境可能内存不足，并且还要给设备加个wifi。

所以考虑用go写个命令行UI（编译出来就是可执行文件，而且没啥依赖）协助同事，帮助同事配置设备。

目前计划做一个通用的json文件读写，并且能依据配置执行部分简单linux指令的命令行UI，后续考虑ymal toml等格式的文件支持。

## 如何使用



```
git clone https://github.com/sanqiandamowang/peon.git
cd peon
go mod tidy
go build .
```

项目有个配置文件

```
config/config.json

{
    "version": "0.0.1",
    "configDir":"testValue",
    "plugins":[
        {
            "name":"plugin1",
            "cmd":"ls"
        },
        {
            "name":"plugin2",
            "cmd":"clear"
        }
    ]

}
```

使用启动选项命令或配置文件来更改默认设置：

```sh
$ ./pono -h
可视化编辑相应文件的cli工具，配置文件在config/config.json中。

Usage:
  peon [flags]
  peon [command]

Available Commands:
  completion  生成shell自动补全脚本
  help        Help about any command

Flags:
  -d, --dir string   打开文件目录下所有的符合格式的文件(config/config.json存在则读取配置),若无使用默认配置 (default "./")
  -h, --help         help for peon

Use "peon [command] --help" for more information about a command.
```

使用cobra提供的自动补全功能(当前会话生效)

```sh
$ ./peno completion bash > _peon_completion
$ source _peo'n_completion_completion
```

配置 `configDir`为 peon 要操作的文件的文件夹地址

```
./pean
```

显示这个界面

```
┌─Info─────────────────────────────────────────────────────────────────────────────────────────────┐ ┌─debug─────────────────────────────────────────┐
│version:0.0.1   json_dir:  testValue                                                              │ │system strat                                   │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘ │                                               │
┌─commands─────────────────────────────────────────────────────────────────────────────────────────┐ │                                               │
│Edit JSON                                                                                         │ │                                               │
│plugin1 :  ls                                                                                     │ │                                               │
│plugin2 :  clear                                                                                  │ │                                               │                                                      │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘ 

```



目前只支持 Edit JSON 功能  

方向键上下键选择 

回车进入下一级

方向键左键返回上一级
`file tree`界面回车进入 `file edit`界面，

`file edit`界面 ctrl+s 返回`file tree`界面

`file tree`界面 ctrl+s 保存文件 
