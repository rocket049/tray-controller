package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocket049/gettext-go/gettext"

	"github.com/gotk3/gotk3/gtk"
	go_locale "github.com/jmshal/go-locale"
)

type myCfg struct {
	Exec string
	Args string
	Envs string
	Wd   string
}

type myApp struct {
	window      *gtk.Window
	button      *gtk.Button
	cfgName     *gtk.Entry
	cfgExec     *gtk.Entry
	cfgArgs     *gtk.Entry
	cfgEnvs     *gtk.Entry
	cfgWd       *gtk.Entry
	cfgRunIcon  *gtk.Entry
	cfgStopIcon *gtk.Entry
}

var app *myApp
var osID int

func getCfgDir(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var dir1 string
	if len(name) > 0 {
		dir1 = filepath.Join(home, "config", name)
	} else {
		return "", err
	}
	if osID == 1 {
		dir1 = dir1 + ".exe"
	}
	os.MkdirAll(dir1, os.ModePerm)
	return dir1, nil
}

func getBinPath(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("Must give me a name.")
	}
	switch osID {
	case 0:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir1 := filepath.Join(home, "bin")
		os.MkdirAll(dir1, os.ModePerm)
		return filepath.Join(dir1, name), nil
	case 1:
		exe, err := os.Executable()
		if err != nil {
			return "", err
		}
		dir1 := filepath.Dir(exe)
		return filepath.Join(dir1, name+".exe"), nil
	}
	return "", nil
}

func getControllerPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir1 := filepath.Dir(exe)
	var res string
	switch osID {
	case 0:
		res = filepath.Join(dir1, "traycontroller")
	case 1:
		res = filepath.Join(dir1, "controller.exe")
	}
	return res, nil
}

func getIconDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir1 := filepath.Dir(exe)
	return filepath.Join(dir1, "..", "share", "traycontroller", "icons")
}

func getLocaleDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir1 := filepath.Dir(exe)
	return filepath.Join(dir1, "..", "share", "traycontroller", "locale")
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *myApp) Create() {
	var err error
	s.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	errPanic(err)
	grid, err := gtk.GridNew()
	errPanic(err)
	s.window.Add(grid)

	label, err := gtk.LabelNew(gettext.T("Create New Controller Config File."))
	errPanic(err)
	label.SetSizeRequest(600, 30)
	grid.Attach(label, 0, 0, 2, 1)
	label, err = gtk.LabelNew(gettext.T("ControllerName:"))
	errPanic(err)
	grid.Attach(label, 0, 1, 1, 1)
	s.cfgName, err = gtk.EntryNew()
	s.cfgName.SetHExpand(true)
	errPanic(err)
	grid.Attach(s.cfgName, 1, 1, 1, 1)
	switch osID {
	case 0:
		s.cfgName.SetTooltipText(gettext.T("I will copy the program to HOME/bin/ControllerName"))
	case 1:
		s.cfgName.SetTooltipText(gettext.T("I will copy the program to 'ControllerName'"))
	}

	label, err = gtk.LabelNew(gettext.T("Exec:"))
	errPanic(err)
	grid.Attach(label, 0, 2, 1, 1)
	s.cfgExec, err = gtk.EntryNew()
	errPanic(err)
	s.cfgExec.SetIconFromIconName(gtk.ENTRY_ICON_SECONDARY, "edit-find")
	grid.Attach(s.cfgExec, 1, 2, 1, 1)
	s.cfgExec.SetTooltipText(gettext.T("Click ICON to show Choose Dialog"))

	label, err = gtk.LabelNew(gettext.T("Args:"))
	errPanic(err)
	grid.Attach(label, 0, 3, 1, 1)
	s.cfgArgs, err = gtk.EntryNew()
	errPanic(err)
	grid.Attach(s.cfgArgs, 1, 3, 1, 1)
	s.cfgArgs.SetTooltipText(gettext.T("Split arguments with SPACE, example:") + "args1 arg2 ...")

	label, err = gtk.LabelNew(gettext.T("Envs:"))
	errPanic(err)
	grid.Attach(label, 0, 4, 1, 1)
	s.cfgEnvs, err = gtk.EntryNew()
	errPanic(err)
	grid.Attach(s.cfgEnvs, 1, 4, 1, 1)
	s.cfgEnvs.SetTooltipText(gettext.T("Split enviroment with ';', example:") + "Key1=Value1;Key2=Value2;...")

	label, err = gtk.LabelNew(gettext.T("Wd:"))
	errPanic(err)
	grid.Attach(label, 0, 5, 1, 1)
	s.cfgWd, err = gtk.EntryNew()
	errPanic(err)
	s.cfgWd.SetIconFromIconName(gtk.ENTRY_ICON_SECONDARY, "edit-find")
	grid.Attach(s.cfgWd, 1, 5, 1, 1)
	s.cfgWd.SetTooltipText(gettext.T("Click ICON to show Choose Dialog"))

	label, err = gtk.LabelNew(gettext.T("RunIcon:"))
	errPanic(err)
	grid.Attach(label, 0, 6, 1, 1)
	s.cfgRunIcon, err = gtk.EntryNew()
	errPanic(err)
	s.cfgRunIcon.SetIconFromIconName(gtk.ENTRY_ICON_SECONDARY, "edit-find")
	grid.Attach(s.cfgRunIcon, 1, 6, 1, 1)
	s.cfgRunIcon.SetTooltipText(gettext.T("Click ICON to show Choose Dialog"))

	label, err = gtk.LabelNew(gettext.T("StopIcon:"))
	errPanic(err)
	grid.Attach(label, 0, 7, 1, 1)
	s.cfgStopIcon, err = gtk.EntryNew()
	errPanic(err)
	s.cfgStopIcon.SetIconFromIconName(gtk.ENTRY_ICON_SECONDARY, "edit-find")
	grid.Attach(s.cfgStopIcon, 1, 7, 1, 1)
	s.cfgStopIcon.SetTooltipText(gettext.T("Click ICON to show Choose Dialog"))

	s.button, err = gtk.ButtonNewWithLabel(gettext.T("OK"))
	errPanic(err)
	grid.Attach(s.button, 0, 8, 2, 1)

	grid.ShowAll()
	s.window.Show()
	s.setSignals()
}

func (s *myApp) setSignals() {
	s.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	s.cfgExec.Connect("icon-press", func() {
		dlg, err := gtk.FileChooserNativeDialogNew(gettext.T("Choose a Executable"), s.window, gtk.FILE_CHOOSER_ACTION_OPEN, gettext.T("OK"), gettext.T("Cancel"))
		errPanic(err)
		dlg.SetSelectMultiple(false)
		ret := dlg.Run()
		if ret == int(gtk.RESPONSE_ACCEPT) {
			s.cfgExec.SetText(dlg.GetFilename())
			s.cfgExec.SetEditable(false)
		}
	})

	s.cfgWd.Connect("icon-press", func() {
		dlg, err := gtk.FileChooserNativeDialogNew(gettext.T("Choose a Dir"), s.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER, gettext.T("OK"), gettext.T("Cancel"))
		errPanic(err)
		dlg.SetSelectMultiple(false)
		ret := dlg.Run()
		if ret == int(gtk.RESPONSE_ACCEPT) {
			s.cfgWd.SetText(dlg.GetFilename())
		}
	})

	s.cfgRunIcon.Connect("icon-press", func() {
		dlg, err := gtk.FileChooserNativeDialogNew(gettext.T("Choose a PNG Image"), s.window, gtk.FILE_CHOOSER_ACTION_OPEN, gettext.T("OK"), gettext.T("Cancel"))
		errPanic(err)
		dlg.SetSelectMultiple(false)
		dlg.SetCurrentFolder(getIconDir())
		ret := dlg.Run()
		if ret == int(gtk.RESPONSE_ACCEPT) {
			s.cfgRunIcon.SetText(dlg.GetFilename())
			s.cfgRunIcon.SetEditable(false)
		}
	})

	s.cfgStopIcon.Connect("icon-press", func() {
		dlg, err := gtk.FileChooserNativeDialogNew(gettext.T("Choose a PNG Image"), s.window, gtk.FILE_CHOOSER_ACTION_OPEN, gettext.T("OK"), gettext.T("Cancel"))
		errPanic(err)
		dlg.SetSelectMultiple(false)
		dlg.SetCurrentFolder(getIconDir())
		ret := dlg.Run()
		if ret == int(gtk.RESPONSE_ACCEPT) {
			s.cfgStopIcon.SetText(dlg.GetFilename())
			s.cfgStopIcon.SetEditable(false)
		}
	})

	s.button.Connect("clicked", func() {
		s.makeConfig()
		gtk.MainQuit()
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	fp1, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer fp1.Close()
	fp2, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fp2.Close()
	_, err = io.Copy(fp1, fp2)
	return err
}

func makeLauncher(binPath, iconPath string) error {
	if len(binPath) == 0 {
		return errors.New("Must provide program pathname")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	name1 := filepath.Base(binPath)
	path1 := filepath.Join(home, ".local", "share", "applications", name1+".desktop")
	tmpl := `[Desktop Entry]
Encoding=UTF-8
Version=1.0
Type=Application
Terminal=false
Exec=` + binPath + `
Name=` + name1 + `
Icon=` + iconPath + `
Categories=GTK;Utility;
Comment=Tray Controller`
	return ioutil.WriteFile(path1, []byte(tmpl), 0755)
}

func zeroPanic(s string) {
	name1 := strings.Trim(s, " /")
	if len(name1) == 0 {
		panic("get zero value.")
	}
}

func (s *myApp) makeConfig() {
	name1, err := s.cfgName.GetText()
	errPanic(err)
	name1 = strings.Trim(name1, " ")
	zeroPanic(name1)

	//copy program
	binPath, err := getBinPath(name1)
	errPanic(err)
	progPath, err := getControllerPath()
	errPanic(err)
	copyFile(progPath, binPath, 0755)

	cfgDir, err := getCfgDir(name1)
	//copy icons
	runIconPath := filepath.Join(cfgDir, "run.png")
	src, err := s.cfgRunIcon.GetText()
	errPanic(err)
	zeroPanic(src)
	copyFile(src, runIconPath, 0644)

	stopIconPath := filepath.Join(cfgDir, "stop.png")
	src, err = s.cfgStopIcon.GetText()
	errPanic(err)
	zeroPanic(src)
	copyFile(src, stopIconPath, 0644)

	//create config files
	exec, err := s.cfgExec.GetText()
	errPanic(err)
	zeroPanic(exec)

	args, err := s.cfgArgs.GetText()
	errPanic(err)
	envs, err := s.cfgEnvs.GetText()
	errPanic(err)
	wd, err := s.cfgWd.GetText()

	cfg := &myCfg{Exec: exec, Args: args, Envs: envs, Wd: wd}
	buf, err := json.Marshal(cfg)
	errPanic(err)
	fp, err := os.Create(filepath.Join(cfgDir, "app.json"))
	errPanic(err)
	defer fp.Close()
	fp.Write(buf)

	//create launcher
	makeLauncher(binPath, runIconPath)
}

func main() {
	gtk.Init(&os.Args)
	loc, err := go_locale.DetectLocale()
	if err != nil {
		loc = ""
		log.Println(loc)
	}
	gettext.SetLocale(loc)
	gettext.BindTextdomain("config", getLocaleDir(), nil)
	gettext.Textdomain("config")
	app = new(myApp)
	app.Create()
	gtk.Main()
	//gettext.SaveLog()
}
