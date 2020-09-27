package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	mathrand "math/rand"
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
	go func() {
		for {
			if mainWindow != nil && mainWindow.Visible() {
				mainWindow.SetSize(walk.Size{500, 301})
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		for {
			time.Sleep(time.Second)
			statusBar.SetText("normal")
			networkBar.SetText("default")
			upFlowBar.SetText(fmt.Sprintf("Up:%.1fkb/s", randFloat()))
			downFlowBar.SetText(fmt.Sprintf("Down:%.1fkb/s", randFloat()))
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

type LoginInfo struct {
	User     string
	Password string
	Remenber bool
	Auto     bool
}

func LoginDialog(loginOld *LoginInfo)  *LoginInfo {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	var remenber, auto *walk.RadioButton
	var result *walk.Label

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "Login",
		Icon:          walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "login",
			DataSource:     loginOld,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{300, 220},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "User:",
					},
					LineEdit{
						Text: Bind("User"),
					},
					Label{
						Text: "Password:",
					},
					LineEdit{
						Text: Bind("Password"),
						PasswordMode: true,
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							RadioButton{
								AssignTo: &remenber,
								Text: "记住密码",
								OnClicked: func() {
									if loginOld.Remenber == true {
										loginOld.Remenber = false
										remenber.SetChecked(false)
									} else {
										loginOld.Remenber = true
										remenber.SetChecked(true)
									}
									log.Printf("value : %v", remenber.Value())
								},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							RadioButton{
								AssignTo: &auto,
								Text: "自动登陆",
								OnClicked: func() {
									if loginOld.Auto == true {
										loginOld.Auto = false
										auto.SetChecked(false)
									} else {
										loginOld.Remenber = true
										remenber.SetChecked(true)

										loginOld.Auto = true
										auto.SetChecked(true)
									}
								},
							},
						},
					},
					Label{
						AssignTo: &result,
						Text: "",
						MinSize: Size{Height: 15},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								result.SetText(err.Error())
								return
							}
							if loginOld.User == "" {
								result.SetText("请输入用户名")
								return
							}
							if loginOld.Password == "" {
								result.SetText("请输入密码")
								return
							}
							acceptPB.SetEnabled(false)
							cancelPB.SetEnabled(false)
							result.SetText("登陆中..")

							go func() {
								time.Sleep(time.Second)
								result.SetText("登陆失败")

								//dlg.Accept()
								acceptPB.SetEnabled(true)
								cancelPB.SetEnabled(true)
							}()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	log.Printf("%v\n", loginOld)

	return loginOld
}

var loginInfo LoginInfo

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Action{
					Text:     "&Login",
					OnTriggered: func() {
						LoginDialog(&loginInfo)
					},
				},
				Action{
					Text:     "&Setting",
				},
				Action{
					Text:     "&Export",
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
					Text:        "&Register",
				},
				Separator{},
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