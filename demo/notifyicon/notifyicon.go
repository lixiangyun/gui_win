// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

var mainWindow *walk.MainWindow

func notify()  {
	// We load our icon from a file.
	icon, err := walk.Resources.Icon("../img/stop.ico")
	if err != nil {
		log.Fatal(err)
	}



	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		log.Fatal(err)
	}

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.",
			icon); err != nil {

			log.Fatal(err)
		}
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("&MainWindow"); err != nil {
		log.Fatal(err)
	}

	exitAction.Triggered().Attach(func() {
		mainWindow.SetVisible(true)

		//walk.App().Exit(0)
	})

	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		log.Fatal(err)
	}

	mainWindow.SetVisible(false)
}

func main() {
	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.

	icon1, err := walk.NewIconFromFile("../img/check.ico")
	if err != nil {
		log.Fatal(err)
	}
	icon2, err := walk.NewIconFromFile("../img/stop.ico")
	if err != nil {
		log.Fatal(err)
	}

	var sbi *walk.StatusBarItem

	mw := &MainWindow{
		AssignTo: &mainWindow,
		Title:   "Walk Statusbar Example",
		MinSize: Size{600, 200},
		Layout:  VBox{MarginsZero: true},
		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo: &sbi,
				Icon:     icon1,
				Text:     "click",
				Width:    80,
				OnClicked: func() {
					if sbi.Text() == "click" {
						sbi.SetText("again")
						sbi.SetIcon(icon2)
					} else {
						sbi.SetText("click")
						sbi.SetIcon(icon1)
					}
				},
			},
			StatusBarItem{
				Text:        "notify",
				ToolTipText: "no tooltip for me",
				OnClicked: func() {
					notify()
				},
			},
			StatusBarItem{
				Text: "\tcenter",
			},
			StatusBarItem{
				Text: "\t\tright",
			},
			StatusBarItem{
				Icon:        icon1,
				ToolTipText: "An icon with a tooltip",
			},
		},
	}

	// Run the message loop.
	mw.Run()
}
