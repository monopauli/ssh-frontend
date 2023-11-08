package ui

import (
	"database/sql"
	"frontend/data"
	"frontend/db"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

var myWindow fyne.Window

func CreateUI() {
	selectedListItem := 0
	database := db.Connect()
	hosts := db.SelectALL(database)
	networks := db.SelectAllNetworks(database)

	// for _, host := range hosts {
	// 	fmt.Println(host)
	// }

	myApp := app.New()
	myWindow = myApp.NewWindow("Hello")

	hostList := createHostList(&hosts, myWindow, database)
	hostList.Resize(fyne.NewSize(100, 100))
	hostBox := CreateHostInformationBox(selectedListItem, &hosts, &networks, database, hostList)
	hostDetailsArea := container.NewHBox() // Create a container to hold the host information box

	//Host List Creation
	hostList.OnSelected = func(id widget.ListItemID) {
		selectedListItem = id
		hostDetailsArea.Remove(hostBox)
		hostBox = CreateHostInformationBox(selectedListItem, &hosts, &networks, database, hostList)
		hostDetailsArea.Add(hostBox)
		hostDetailsArea.Refresh()
	}

	//Menu Creation
	menuBar := createMenuBar(myWindow, &hosts, database, hostList)
	myWindow.SetMainMenu(menuBar)

	myWindow.SetContent(
		container.NewHSplit(
			hostList,
			hostDetailsArea, // This will be updated when an item is selected
		),
	)
	myWindow.ShowAndRun()

}

func createHostList(hosts *[]data.Host, window fyne.Window, database *sql.DB) *widget.List {
	hostList := widget.NewList(
		func() int {
			return len(*hosts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText((*hosts)[i].Host)
		},
	)
	return hostList
}

func RefreshUI(database *sql.DB, hosts *[]data.Host, list *widget.List) {
	(*hosts) = db.SelectALL(database)
	list.Refresh()
}
