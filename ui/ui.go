package ui

import (
	"fmt"
	"frontend/data"
	"frontend/db"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func CreateUI() {
	var selectedListItem int
	var selectedHostId int
	database := db.Connect()
	hosts := db.SelectALL(database)

	for _, host := range hosts {
		fmt.Println(host)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")

	//Host List Creation
	hostList := widget.NewList(
		func() int {
			return len(hosts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(hosts[i].Host)
		})

	//Menu
	fileExportItem := fyne.NewMenuItem("Export", func() {})
	editNewHostItem := fyne.NewMenuItem("New Host", func() {
		newHost := data.Host{
			Host:         "New Host",
			Hostname:     "",
			User:         "",
			Port:         0,
			IdentityFile: "",
			SystemType:   "",
			NodeType:     "Sentry",
			Provider:     "",
			Region:       "",
			InternalIP:   "",
			Portbase:     0,
		}

		id, err := db.AddEntry(database, newHost)
		if err != nil {
			fmt.Println(err)
		}
		hosts = db.SelectALL(database)
		hostList.Refresh()
		hostList.Select(data.FindHost(id, hosts))
	})

	fileMenu := fyne.NewMenu("File", fileExportItem)
	editMenu := fyne.NewMenu("Edit", editNewHostItem)

	menuBar := fyne.NewMainMenu(
		fileMenu,
		editMenu,
	)
	myWindow.SetMainMenu(menuBar)

	//Host Information
	hostLabel := widget.NewLabel("Host")
	hostInput := widget.NewEntry()
	nameLabel := widget.NewLabel("Hostname")
	nameInput := widget.NewEntry()
	userLabel := widget.NewLabel("User")
	userInput := widget.NewEntry()
	portLabel := widget.NewLabel("Port")
	portInput := widget.NewEntry()
	keyLabel := widget.NewLabel("Key")
	keyInput := widget.NewEntry()
	systemTypeLabel := widget.NewLabel("System Type")
	systemTypeInput := widget.NewEntry()
	nodeTypeLabel := widget.NewLabel("Node Type")
	nodeTypeInput := widget.NewEntry()
	providerLabel := widget.NewLabel("Provider")
	providerInput := widget.NewEntry()
	regionLabel := widget.NewLabel("Region")
	regionInput := widget.NewEntry()
	internalIPLabel := widget.NewLabel("Internal IP")
	internalIPInput := widget.NewEntry()
	portbaseLabel := widget.NewLabel("Portbase")
	portbaseInput := widget.NewEntry()

	//Host List Selection
	hostList.OnSelected = func(id widget.ListItemID) {
		hostInput.SetText(hosts[id].Host)
		nameInput.SetText(hosts[id].Hostname)
		userInput.SetText(hosts[id].User)
		portInput.SetText(fmt.Sprint(hosts[id].Port))
		keyInput.SetText(hosts[id].IdentityFile)
		systemTypeInput.SetText(hosts[id].SystemType)
		nodeTypeInput.SetText(hosts[id].NodeType)
		providerInput.SetText(hosts[id].Provider)
		regionInput.SetText(hosts[id].Region)
		internalIPInput.SetText(hosts[id].InternalIP)
		portbaseInput.SetText(fmt.Sprint(hosts[id].Portbase))
		selectedListItem = id
		selectedHostId = hosts[id].ID
	}

	//Buttons
	saveButton := widget.NewButton("Save", func() {
		oldHost := hosts[selectedListItem]
		port, _ := strconv.Atoi(portInput.Text)
		portbase, _ := strconv.Atoi(portbaseInput.Text)
		editedHost := data.Host{
			ID:           oldHost.ID,
			Host:         hostInput.Text,
			Hostname:     nameInput.Text,
			User:         userInput.Text,
			Port:         port,
			IdentityFile: keyInput.Text,
			SystemType:   systemTypeInput.Text,
			NodeType:     nodeTypeInput.Text,
			Provider:     providerInput.Text,
			Region:       regionInput.Text,
			InternalIP:   internalIPInput.Text,
			Portbase:     portbase,
		}

		result := data.CompareStructs(oldHost, editedHost)
		var builder strings.Builder
		builder.WriteString("Do you want to make these changes?\n")

		for _, v := range result {
			builder.WriteString(fmt.Sprintf("%s: %s -> %s\n", v[0], v[1], v[2]))
		}
		changes := builder.String()
		dialog.ShowConfirm("Warning", changes, func(confirm bool) {
			if confirm {
				id := editedHost.ID
				db.Update(database, editedHost, selectedHostId)
				hosts = db.SelectALL(database)
				hostList.Refresh()
				hostList.Select(data.FindHost(id, hosts))

			}
		}, myWindow)

	})
	deleteButton := widget.NewButton("Delete", func() {
		message := fmt.Sprintf("Do you really want to delete this Host: %s?", hosts[selectedListItem].Host)
		dialog.ShowConfirm("Warning", message, func(confirm bool) {
			if confirm {
				fmt.Println("Host to delete: ")
				fmt.Println(hosts[selectedListItem])
				db.DeleteEntry(database, hosts[selectedListItem])
				hosts = db.SelectALL(database)
				hostList.Refresh()
			}
		}, myWindow)

	})

	hostBox := container.NewGridWithColumns(
		2,
		hostLabel,
		hostInput,
		nameLabel,
		nameInput,
		userLabel,
		userInput,
		portLabel,
		portInput,
		keyLabel,
		keyInput,
		systemTypeLabel,
		systemTypeInput,
		nodeTypeLabel,
		nodeTypeInput,
		providerLabel,
		providerInput,
		regionLabel,
		regionInput,
		internalIPLabel,
		internalIPInput,
		portbaseLabel,
		portbaseInput,
		saveButton,
		deleteButton,
	)
	myWindow.SetContent(container.NewHBox(
		hostList,
		hostBox,
	))
	myWindow.ShowAndRun()

}
