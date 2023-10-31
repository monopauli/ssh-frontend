package ui

import (
	"fmt"
	"frontend/data"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func CreateUI() {

	var selectedHostID int
	hosts := data.GetHosts(*data.OpenConfig("config/hosts"))
	for _, host := range hosts {
		fmt.Println(host)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")

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
	nodeNetworksLabel := widget.NewLabel("Node Networks")
	nodeNetworksInput := widget.NewEntry()
	providerLabel := widget.NewLabel("Provider")
	providerInput := widget.NewEntry()
	regionLabel := widget.NewLabel("Region")
	regionInput := widget.NewEntry()
	internalIPLabel := widget.NewLabel("Internal IP")
	internalIPInput := widget.NewEntry()
	portbaseLabel := widget.NewLabel("Portbase")
	portbaseInput := widget.NewEntry()

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
	hostList.OnSelected = func(id widget.ListItemID) {
		nameInput.SetText(hosts[id].Name)
		userInput.SetText(hosts[id].User)
		portInput.SetText(hosts[id].Port)
		keyInput.SetText(hosts[id].Key)
		systemTypeInput.SetText(hosts[id].SystemType)
		nodeTypeInput.SetText(hosts[id].NodeType)
		nodeNetworksInput.SetText(hosts[id].NodeNetworks)
		providerInput.SetText(hosts[id].Provider)
		regionInput.SetText(hosts[id].Region)
		internalIPInput.SetText(hosts[id].InternalIP)
		portbaseInput.SetText(hosts[id].Portbase)
		selectedHostID = id
	}

	newButton := widget.NewButton("New Host", func() {
		newHost := data.Host{
			Host:         "New Host",
			Name:         "",
			User:         "",
			Port:         "",
			Key:          "",
			SystemType:   "",
			NodeType:     "",
			NodeNetworks: "",
			Provider:     "",
			Region:       "",
			InternalIP:   "",
			Portbase:     "",
		}
		hosts = append(hosts, newHost)
		hostList.Refresh()
		hostList.Select(len(hosts) - 1)
	})
	saveButton := widget.NewButton("Save", func() {
		name := nameInput.Text
		user := userInput.Text
		port := portInput.Text
		key := keyInput.Text
		systemType := systemTypeInput.Text
		nodeType := nodeTypeInput.Text
		nodeNetworks := nodeNetworksInput.Text
		provider := providerInput.Text
		region := regionInput.Text
		internalIP := internalIPInput.Text
		portbase := portbaseInput.Text
		saveHost(
			hosts[:],
			selectedHostID,
			name,
			user,
			port,
			key,
			systemType,
			nodeType,
			nodeNetworks,
			provider,
			region,
			internalIP,
			portbase,
		)
	})

	hostBox := container.NewGridWithColumns(
		2,
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
		nodeNetworksLabel,
		nodeNetworksInput,
		providerLabel,
		providerInput,
		regionLabel,
		regionInput,
		internalIPLabel,
		internalIPInput,
		portbaseLabel,
		portbaseInput,
		saveButton,
		newButton,
	)
	myWindow.SetContent(container.NewHBox(
		hostList,
		hostBox,
	))
	myWindow.ShowAndRun()

}

func saveHost(
	hosts []data.Host,
	id int,
	name string,
	user string,
	port string,
	key string,
	systemType string,
	nodeType string,
	nodeNetworks string,
	provider string,
	region string,
	internalIP string,
	portbase string) {
	hosts[id].Name = name
	hosts[id].User = user
	hosts[id].Port = port
	hosts[id].Key = key
	hosts[id].SystemType = systemType
	hosts[id].NodeType = nodeType
	hosts[id].NodeNetworks = nodeNetworks
	hosts[id].Provider = provider
	hosts[id].Region = region
	hosts[id].InternalIP = internalIP
	hosts[id].Portbase = portbase

	fmt.Println("Saved Host: ")
	fmt.Println(hosts[id])
}
