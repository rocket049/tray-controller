//+build gtk_3_12

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var application *gtk.Application
var myapp *myApp

type myApp struct {
	window     *gtk.ApplicationWindow
	grid       *gtk.Grid
	scrollWin  *gtk.ScrolledWindow
	board      *gtk.Label
	input      *gtk.Entry
	popMenu    *gtk.Menu
	itemStatus *gtk.MenuItem
	tray       *gtk.StatusIcon
	iconRun    string
	iconStop   string
	status     bool
	running    bool
	myCmd      *myCommand
	srv        *exec.Cmd
	chStaus    chan int
	pipeStdin  io.Writer
}

func (s *myApp) Create() {
	exe1, err := os.Executable()
	if err != nil {
		panic(err)
	}
	name1 := filepath.Base(exe1)
	home1, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	json1 := filepath.Join(home1, "config", name1, "app.json")
	cfg, err := GetCfgFromJSON(json1)
	if err != nil {
		s.showHelp(json1)
		panic(err)
	}
	s.myCmd = new(myCommand)
	s.myCmd.Name = cfg.Exec
	args := []string{s.myCmd.Name}

	sp := strings.Split(cfg.Args, " ")
	for _, v := range sp {
		if len(v) == 0 {
			continue
		}
		args = append(args, v)
	}
	s.myCmd.Args = args
	//fmt.Println(s.myCmd.Args)

	sp = strings.Split(cfg.Envs, ";")
	envs := []string{}
	for _, v := range os.Environ() {
		envs = append(envs, v)
	}
	for _, v := range sp {
		if len(v) == 0 {
			continue
		}
		envs = append(envs, v)
	}
	s.myCmd.Envs = envs
	if len(cfg.Wd) > 0 {
		os.Chdir(cfg.Wd)
	}

	s.iconRun = path.Join(home1, "config", name1, "run.png")
	s.iconStop = path.Join(home1, "config", name1, "stop.png")

	s.window, err = gtk.ApplicationWindowNew(application)
	if err != nil {
		panic(err)
	}

	s.window.SetSizeRequest(400, 300)
	s.window.SetTitle("Controller-" + s.myCmd.Name)
	s.grid, err = gtk.GridNew()
	if err != nil {
		panic(err)
	}
	s.window.Add(s.grid)
	s.scrollWin, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.scrollWin.SetHExpand(true)
	s.scrollWin.SetVExpand(true)
	s.scrollWin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	s.grid.Attach(s.scrollWin, 0, 0, 1, 1)
	s.board, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal(err)
	}
	s.board.SetLineWrap(true)
	s.grid.SetMarginStart(10)
	s.grid.SetMarginEnd(10)
	s.scrollWin.Add(s.board)
	s.input, err = gtk.EntryNew()
	if err != nil {
		log.Fatal(err)
	}
	s.input.SetHExpand(true)
	s.grid.Attach(s.input, 0, 1, 1, 1)
	s.grid.ShowAll()
	s.window.Show()
	s.setTray()
	s.setSignals()
}

func (s *myApp) setSignals() {
	//signals
	s.input.Connect("activate", func() {
		txt1, _ := s.input.GetText()
		defer s.input.SetText("")
		s.pipeStdin.Write([]byte(txt1 + "\n"))
	})

	s.window.Connect("delete-event", func() bool {
		s.window.HideOnDelete()
		s.window.Hide()
		return true
	})
}

func (s *myApp) setTray() {
	s.status = false
	s.chStaus = make(chan int, 1)
	go s.waitClose()

	s.popMenu = s.createMenu()
	var err error
	s.tray, err = gtk.StatusIconNewFromFile(s.iconStop)
	if err != nil {
		log.Fatal("Could not create status icon:", err)
	}
	s.tray.SetVisible(true)
	s.tray.SetHasTooltip(true)

	s.changeStatus()

	s.tray.Connect("popup-menu", func(statusIcon *gtk.StatusIcon, button, activateTime uint) {
		s.popMenu.PopupAtStatusIcon(statusIcon, button, uint32(activateTime))
	})
	s.tray.Connect("activate", func() {
		s.changeStatus()
	})
}

func (s *myApp) changeStatus() {
	if s.status {
		s.srv.Process.Signal(os.Interrupt)
		s.status = false
	} else {
		go s.runCmd()
		s.status = true
	}
}

func (s *myApp) createMenu() *gtk.Menu {
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}
	item, err := gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		log.Fatal(err)
	}
	item.Show()
	item.Connect("activate", func() {
		if s.status {
			s.srv.Process.Signal(os.Interrupt)
		}
		application.Quit()
	})
	menu.Append(item)
	s.itemStatus, err = gtk.MenuItemNewWithLabel("Run")
	if err != nil {
		log.Fatal(err)
	}
	s.itemStatus.Show()
	s.itemStatus.Connect("activate", func() {
		menu.Popdown()
		s.changeStatus()
	})
	menu.Append(s.itemStatus)

	item, err = gtk.MenuItemNewWithLabel("ShowWindow")
	if err != nil {
		log.Fatal(err)
	}
	item.Show()
	item.Connect("activate", func() {
		s.window.ShowAll()
	})
	menu.Append(item)
	return menu
}

func (s *myApp) runCmd() {
	s.srv = s.myCmd.GetCmd()
	if s.srv == nil {
		s.chStaus <- 2
		return
	}
	r, _ := s.srv.StdoutPipe()
	go s.showOutput(r)
	r, _ = s.srv.StderrPipe()
	go s.showOutput(r)
	s.chStaus <- 1
	w, _ := s.srv.StdinPipe()
	defer w.Close()
	s.pipeStdin = w
	//fmt.Println("run")
	err := s.srv.Run()
	if err != nil {
		log.Println(err)
	}
	s.chStaus <- 3
	//fmt.Println("stop")
}

func (s *myApp) waitClose() {
	for {
		ret, ok := <-s.chStaus
		if !ok {
			application.Quit()
		}
		if ret == 1 {
			glib.IdleAdd(s.updateRun)
			//fmt.Println("update run")
		} else if ret == 3 {
			glib.IdleAdd(s.updateStop)
			//fmt.Println("update stop")
		} else if ret == 2 {
			application.Quit()
		}
	}
}

func (s *myApp) updateRun() {
	s.tray.SetFromFile(s.iconRun)
	s.tray.SetTooltipText(s.myCmd.Name + " Running")
	s.tray.SetHasTooltip(true)
	s.itemStatus.SetLabel("Stop")
	s.status = true
}

func (s *myApp) updateStop() {
	s.tray.SetFromFile(s.iconStop)
	s.tray.SetTooltipText(s.myCmd.Name + " Stopped")
	s.tray.SetHasTooltip(true)
	s.itemStatus.SetLabel("Run")
	s.status = false
}

func (s *myApp) showOutput(r io.ReadCloser) {
	defer r.Close()
	rd := bufio.NewReader(r)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Println(err)
			break
		}
		glib.IdleAdd(func() bool {
			txt1, _ := s.board.GetText()
			s.board.SetText(txt1 + string(line) + "\n")

			glib.IdleAdd(func() bool {
				vadj := s.scrollWin.GetVAdjustment()
				vadj.SetValue(vadj.GetUpper())
				return false
			})

			return false
		})
	}
}

func (s *myApp) showHelp(name string) {
	msg := `需要配置文件和2个图标！
	配置文件：` + name + ` 包含内容：
	{
		"exec":"/full/path/to/prog",
		"args":"-name2 value1 -name2 value2 ...",
		"envs":"Key1=Value1;Key2=Value2;...",
		"wd":"/path/to/work/dir"
	}
	"args"、"envs"、"wd"可以省略。
	图标和配置文件在同一目录，分别是：
	run.png ：代表正在运行
	stop.png ：代表停止状态`
	dlg := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "OK")
	if dlg == nil {
		return
	}
	dlg.SetTitle("Need Config")
	dlg.SetTooltipText("need config file:" + name)
	dlg.SetMarkup(msg)
	dlg.Run()
	dlg.Destroy()
}

func main() {
	exe1, err := os.Executable()
	if err != nil {
		panic(err)
	}
	application, err = gtk.ApplicationNew("tray.ctrl.p-"+filepath.Base(exe1), glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(err)
	}
	application.Connect("activate", func() {
		if myapp != nil {
			return
		}
		myapp = new(myApp)
		myapp.Create()
	})

	//gtk.Main()
	application.Run(os.Args)
}
