package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"frontend/data"
	"frontend/db"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

// type HostDetails struct {
// 	HostInput            *widget.Entry
// 	NameInput            *widget.Entry
// 	UserInput            *widget.Entry
// 	PortInput            *widget.Entry
// 	KeyInput             *widget.Entry
// 	SystemTypeInput      *widget.Entry
// 	NodeTypeInput        *widget.Entry
// 	ProviderInput        *widget.Entry
// 	RegionInput          *widget.Entry
// 	InternalIPInput      *widget.Entry
// 	PortbaseInput        *widget.Entry
// 	NetworksList         *widget.List
// 	NetworksDropDown     *widget.Select
// 	NetworksAddButton    *widget.Button
// 	NetworksDeleteButton *widget.Button
// 	SaveButton           *widget.Button
// 	DeleteButton         *widget.Button
// }

func CreateHostInformationBox(selectedListItem int, hosts *[]data.Host, allNetworks *[]string, database *sql.DB, hostList *widget.List) *fyne.Container {
	var selectedNetworkID int
	var selectedNetworkDropDown string

	//Create Inputs
	hostInput := widget.NewEntry()
	nameInput := widget.NewEntry()
	userInput := widget.NewEntry()
	portInput := widget.NewEntry()
	keyInput := widget.NewEntry()
	systemTypeInput := widget.NewEntry()
	nodeTypeInput := widget.NewEntry()
	providerInput := widget.NewEntry()
	regionInput := widget.NewEntry()
	internalIPInput := widget.NewEntry()
	portbaseInput := widget.NewEntry()

	// hostDetails := &HostDetails{
	// 	HostInput:       widget.NewEntry(),
	// 	NameInput:       widget.NewEntry(),
	// 	UserInput:       widget.NewEntry(),
	// 	PortInput:       widget.NewEntry(),
	// 	KeyInput:        widget.NewEntry(),
	// 	SystemTypeInput: widget.NewEntry(),
	// 	NodeTypeInput:   widget.NewEntry(),
	// 	ProviderInput:   widget.NewEntry(),
	// 	RegionInput:     widget.NewEntry(),
	// 	InternalIPInput: widget.NewEntry(),
	// 	PortbaseInput:   widget.NewEntry(),
	// 	NetworksList: widget.NewList(
	// 		func() int { return 0 },
	// 		func() fyne.CanvasObject { return widget.NewLabel("template") },
	// 		func(i widget.ListItemID, o fyne.CanvasObject) {
	// 		},
	// 	),

	// }

	//Fill Inputs
	hostInput.SetText((*hosts)[selectedListItem].Host)
	nameInput.SetText((*hosts)[selectedListItem].Hostname)
	userInput.SetText((*hosts)[selectedListItem].User)
	portInput.SetText(strconv.Itoa((*hosts)[selectedListItem].Port))
	keyInput.SetText((*hosts)[selectedListItem].IdentityFile)
	systemTypeInput.SetText((*hosts)[selectedListItem].SystemType)
	nodeTypeInput.SetText((*hosts)[selectedListItem].NodeType)
	providerInput.SetText((*hosts)[selectedListItem].Provider)
	regionInput.SetText((*hosts)[selectedListItem].Region)
	internalIPInput.SetText((*hosts)[selectedListItem].InternalIP)
	portbaseInput.SetText(strconv.Itoa((*hosts)[selectedListItem].Portbase))

	networks := (*hosts)[selectedListItem].Networks

	networksList := widget.NewList(
		func() int { return len(networks) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(networks[i])
		},
	)

	networksList.OnSelected = func(id widget.ListItemID) {
		selectedNetworkID = id
	}

	// Create networks control
	networksDropDown := widget.NewSelect(*allNetworks, func(value string) {
		selectedNetworkDropDown = value
	})
	networksAddButton := widget.NewButton("+", func() {
		db.AddNetworkToHost(database, (*hosts)[selectedListItem].ID, selectedNetworkDropDown)
		// Refresh the hosts list and UI
		RefreshUI(database, hosts, hostList)
		networks = (*hosts)[selectedListItem].Networks
		networksList.Refresh()

	})
	networksDeleteButton := widget.NewButton("-", func() {
		db.DeleteNetworkFromHost(database, networks[selectedNetworkID], (*hosts)[selectedListItem].ID)
		// Refresh the hosts list and UI
		RefreshUI(database, hosts, hostList)
		networks = (*hosts)[selectedListItem].Networks
		networksList.Refresh()
	})

	//Create networks container
	networksButtons := container.NewHBox(networksAddButton, networksDeleteButton)
	networksBox := container.NewVBox(networksList, networksButtons, networksDropDown)

	saveButton := createSaveButton(
		selectedListItem,
		hosts,
		hostInput,
		nameInput,
		userInput,
		portInput,
		keyInput,
		systemTypeInput,
		nodeTypeInput,
		providerInput,
		regionInput,
		internalIPInput,
		portbaseInput,
		database,
		hostList,
		networks)
	deleteButton := createDeleteButton(selectedListItem, hosts, database, hostList)

	labelsAndInputs := []*widget.FormItem{
		widget.NewFormItem("Host", hostInput),
		widget.NewFormItem("Hostname", nameInput),
		widget.NewFormItem("Username", userInput),
		widget.NewFormItem("Port", portInput),
		widget.NewFormItem("Key", keyInput),
		widget.NewFormItem("System Type", systemTypeInput),
		widget.NewFormItem("Node Type", nodeTypeInput),
		widget.NewFormItem("Provider", providerInput),
		widget.NewFormItem("Region", regionInput),
		widget.NewFormItem("Internal IP", internalIPInput),
		widget.NewFormItem("Portbase", portbaseInput),
	}

	hostForm := widget.NewForm(labelsAndInputs...)
	hostForm.Append("Networks", networksBox)

	buttonBox := container.NewHBox(saveButton, deleteButton)

	hostBox := container.NewVBox(hostForm, buttonBox)
	return hostBox
}

func createSaveButton(
	selectedListItem int,
	hosts *[]data.Host,
	hostInput,
	nameInput,
	userInput,
	portInput,
	keyInput,
	systemTypeInput,
	nodeTypeInput,
	providerInput,
	regionInput,
	internalIPInput,
	portbaseInput *widget.Entry,
	database *sql.DB,
	hostList *widget.List,
	networks []string,
) *widget.Button {
	return widget.NewButton("Save", func() {
		// Check if any host is actually selected
		if selectedListItem < 0 || selectedListItem >= len(*hosts) {
			dialog.ShowError(errors.New("no host selected or index out of range"), myWindow)
			return
		}

		oldHost := (*hosts)[selectedListItem]
		port, _ := strconv.Atoi(portInput.Text)         // Handle error appropriately
		portbase, _ := strconv.Atoi(portbaseInput.Text) // Handle error appropriately

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
			Networks:     networks,
		}

		// This function compares the old and new hosts and returns a slice of differences
		differences := data.CompareStructs(oldHost, editedHost)

		if len(differences) == 0 {
			dialog.ShowInformation("No changes", "There are no changes to save.", myWindow)
			return
		}

		var builder strings.Builder
		builder.WriteString("Do you want to make these changes?\n")
		for _, change := range differences {
			builder.WriteString(fmt.Sprintf("%s: %s -> %s\n", change[0], change[1], change[2]))
		}

		dialog.ShowConfirm("Confirm changes", builder.String(), func(confirm bool) {
			if confirm {
				db.Update(database, editedHost, oldHost.ID)
				// Refresh the hosts list and UI
				RefreshUI(database, hosts, hostList)
				newSelectedIndex := data.FindHost(editedHost.ID, *hosts)
				selectedListItem = newSelectedIndex
				hostList.Select(newSelectedIndex)
			}
		}, myWindow)
	})
}

func createDeleteButton(selectedListItem int, hosts *[]data.Host, database *sql.DB, hostList *widget.List) *widget.Button {
	return widget.NewButton("Delete", func() {
		// Check if any host is actually selected
		if selectedListItem < 0 || selectedListItem >= len(*hosts) {
			dialog.ShowError(errors.New("no host selected or index out of range"), myWindow)
			return
		}

		selectedHost := (*hosts)[selectedListItem]
		confirmationMessage := fmt.Sprintf("Do you really want to delete the host: %s?", selectedHost.Host)

		dialog.ShowConfirm("Warning", confirmationMessage, func(confirm bool) {
			if confirm {
				db.DeleteEntry(database, selectedHost)
				// Refresh the hosts list after deletion
				(*hosts) = db.SelectALL(database)
				hostList.Refresh()
				// Reset the selected list item index
				selectedListItem = -1
			}
		}, myWindow)
	})
}
