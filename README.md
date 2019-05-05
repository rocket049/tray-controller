## What is this?
本程序用于控制其他小程序，例如 `pydoc3`、`godoc`，我们直接使用他们的时候需要输入命令，比较低效。
使用本程序作为控制器，可以点击系统托盘控制后台小程序的运行状态。

## 控制方式
为了在一个电脑上控制多个程序，本程序使用程序名字作为配置文件目录，目录路径为：`HOME/config/prog-name`，
目录中需要1个配置文件`app.json`和2个图标`run.png`、`stop.png`。
配置文件包含如下内容：

```
	{
		"exec":"/full/path/to/prog",
		"args":"-name2 value1 -name2 value2 ...",
		"envs":"Key1=Value1;Key2=Value2;...",
		"wd":"/path/to/work/dir"
	}
	
	// godoc 示例：
	{
		"exec":"/usr/local/go/bin/godoc",
		"args":"-http :6060"
	}
	
	// pydoc3 示例：
	{
		"exec":"/usr/bin/pydoc3",
		"args":"-b"
	}
```

其中的"args"、"envs"、"wd"可以省略。

图标和配置文件在同一目录，分别是：

- run.png ：代表正在运行
- stop.png ：代表停止状态

如果没有配置，启动时会弹出提示窗口。

## linux上的图形配置工具
`linux`版中增加了一个配置工具`traycontroller-config`，在命令中运行 `traycontroller-config`，可以用图形界面生成配置文件和菜单项。

`windows 版的 `traycontroller-config ，在解压后的目录中： `bin/traycontroller-config.exe`，不能生成菜单项目，配置后在 bin 目录中寻找新建的 `ControllerName.exe` 运行。

## 界面
主窗口显示程序的输出，包含标准输出和错误输出，下方的输入框可以用来向被控制的程序输入信息。
主窗口关闭后，可以从系统托盘的弹出菜单重新打开。
