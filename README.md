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
```

其中的"args"、"envs"、"wd"可以省略。

图标和配置文件在同一目录，分别是：

- run.png ：代表正在运行
- stop.png ：代表停止状态

如果没有配置，启动时会弹出提示窗口。