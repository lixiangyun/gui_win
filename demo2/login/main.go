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

var StatusBar *walk.StatusBarItem
var flowBar   *walk.StatusBarItem

func randFloat() float32 {
	up := mathrand.Int63() % 10240
	return float32(up)/10
}

var mainWindowWidth = 300
var mainWindowHeight = 600

func init()  {
	go func() {
		for {
			if mainWindow != nil && mainWindow.Visible() {
				mainWindow.SetSize(walk.Size{mainWindowWidth, mainWindowHeight})
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(5 * time.Second)

		var cnt int
		for {
			if (cnt % 2) == 0 {
				StatusBar.SetIcon(network_online1_icon)
			} else {
				StatusBar.SetIcon(network_online2_icon)
			}
			cnt++

			time.Sleep(time.Second)
			flowBar.SetText(fmt.Sprintf("%.1fkb/s", randFloat()))
		}
	}()
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		StatusBarItem{
			AssignTo: &StatusBar,
			Width:    80,
			Icon:     network_offline_icon,
		},
		StatusBarItem{
			AssignTo: &flowBar,
			Icon:     network_flow_icon,
			Width:    120,
			Text:    "0kb/s",
		},
	}
}

type LoginInfo struct {
	User     string
	Password string
	Remenber bool
	Auto     bool
}

func SettingDialog()  {
	
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
		Size: Size{150, 220},
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

var selfIP *walk.LineEdit
var statusShow *walk.LineEdit

type Species struct {
	Id   int
	Name string
}

func KnownSpecies() []*Species {
	return []*Species{
		{1, "Dog"},
		{2, "Cat"},
		{3, "Bird"},
		{4, "Fish"},
		{5, "Elephant"},
	}
}

var network *walk.ComboBox
var iface *walk.ComboBox

func StatusWidget() []Widget {
	return []Widget{
		Label{
			Text: "IP:",
			MaxSize: Size{Width: 50},
		},
		LineEdit{
			AssignTo: &selfIP,
			ReadOnly: true,
			Text: "192.168.0.1",
		},
		Label{
			Text: "Network:",
		},
		ComboBox{
			AssignTo: &network,
			BindingMember: "Id",
			DisplayMember: "Name",
			Model:         KnownSpecies(),
		},
		Label{
			Text: "Interface:",
		},
		ComboBox{
			AssignTo: &iface,

			BindingMember: "Id",
			DisplayMember: "Name",
			Model:         KnownSpecies(),
		},
		Label{
			Text: "Status:",
		},
		LineEdit{
			AssignTo: &statusShow,
			ReadOnly: true,
			Text: "no login",
		},
		HSpacer{},

		Composite{
			Layout: Grid{Columns: 2, MarginsZero: true},
			Children: []Widget{
				PushButton{
					Text:     "&Join",
					OnClicked: func() {

					},
				},
				PushButton{
					Text:     "&Leave",
					OnClicked: func() {

					},
				},
			},
		},
	}
}

func ListWidget() []Widget {
	return []Widget{
		TableView{
			AlternatingRowBG: true,
			CheckBoxes:       true,
			ColumnsOrderable: true,
			Columns: []TableViewColumn{
				{Title: "Tag", Width: 60},
				{Title: "IP", Width: 100},
				{Title: "Status", Width: 50},
			},
		},
	}
}

var main_windows_icon *walk.Icon
var network_offline_icon *walk.Icon
var network_online1_icon *walk.Icon
var network_online2_icon *walk.Icon
var network_flow_icon *walk.Icon


func iconLoad() error {
	var err error

	main_windows_icon, err = walk.NewIconFromFile("./main_windows.ico")
	if err != nil {
		return err
	}

	network_offline_icon, err = walk.NewIconFromFile("./network_offline.ico")
	if err != nil {
		return err
	}

	network_online1_icon, err = walk.NewIconFromFile("./network_online1.ico")
	if err != nil {
		return err
	}

	network_online2_icon, err = walk.NewIconFromFile("./network_online2.ico")
	if err != nil {
		return err
	}

	network_flow_icon, err = walk.NewIconFromFile("./network_flow.ico")
	if err != nil {
		return err
	}
	return nil
}

func main()  {
	err := iconLoad()
	if err != nil {
		log.Println(err.Error())
		return
	}

	mw := MainWindow{
		Icon: main_windows_icon,
		Title:   "客户端程序",
		Size: Size{mainWindowWidth, mainWindowHeight-1},
		Layout:  VBox{},
		StatusBarItems: StatusBarInit(),
		AssignTo: &mainWindow,
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: StatusWidget(),
			},
			Composite{
				Layout: VBox{},
				Children: ListWidget(),
			},
		},
	}
	cnt, err := mw.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("cnt: %d\n", cnt)
}