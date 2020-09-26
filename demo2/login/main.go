package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	mathrand "math/rand"
	"sync"
	"time"
)

func showAboutBoxAction() {
	walk.MsgBox(mainWindow, "About", "v0.1.0 20200926", walk.MsgBoxIconInformation)
}

var mainWindow *walk.MainWindow

var statusBar *walk.StatusBarItem
var networkBar *walk.StatusBarItem

var upFlowBar *walk.StatusBarItem
var downFlowBar *walk.StatusBarItem

func randFloat() float32 {
	up := mathrand.Int63() % 10240
	return float32(up)/10
}

func init()  {
	var once sync.Once

	go func() {
		for {
			time.Sleep(time.Second)

			once.Do(func() {
				mainWindow.SetSize(walk.Size{500, 301})
				mainWindow.SetSize(walk.Size{500, 300})
			})

			statusBar.SetText("normal")
			networkBar.SetText("default")
			upFlowBar.SetText(fmt.Sprintf("Up:%.1fkb/s", randFloat()))
			downFlowBar.SetText(fmt.Sprintf("Up:%.1fkb/s", randFloat()))
		}
	}()
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		StatusBarItem{
			AssignTo: &statusBar,
			Text:     "status",
			Width:    80,
		},
		StatusBarItem{
			AssignTo: &networkBar,
			Text:     "network",
			Width:    80,
		},
		StatusBarItem{
			AssignTo: &upFlowBar,
			Width:    80,
			Text:    "up:0kb/s",
		},
		StatusBarItem{
			AssignTo: &downFlowBar,
			Width:    80,
			Text:    "down:0kb/s",
		},
	}
}

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "&Setting",
			Items: []MenuItem{
				Action{
					Text:        "&login",
				},
				Menu{
					Text:     "Recent",
				},
				Separator{},
				Action{
					Text:        "E&xit",
					OnTriggered: func() { mainWindow.Close() },
				},
			},
		},
		Menu{
			Text: "&Help",
			Items: []MenuItem{
				Action{
					Text:        "About",
					OnTriggered: showAboutBoxAction,
				},
			},
		},
	}
}

func main()  {
	mw := MainWindow{
		Icon: walk.IconInformation(),
		Title:   "客户端程序",
		Size: Size{500, 300},
		Layout:  HBox{},
		StatusBarItems: StatusBarInit(),
		AssignTo: &mainWindow,
		MenuItems: MenuBarInit(),
		Children: []Widget{
			TextEdit{
				ReadOnly: true,
				Text:     "Drop files here, from windows explorer...",
			},
		},
	}
	cnt, err := mw.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("cnt: %d\n", cnt)
}