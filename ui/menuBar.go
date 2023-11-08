package ui

import (
	"database/sql"
	"fmt"
	"frontend/data"
	"frontend/db"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func createMenuBar(window fyne.Window, hosts *[]data.Host, database *sql.DB, hostList *widget.List) *fyne.MainMenu {
	// Define the menu items and their respective actions
	fileExportItem := fyne.NewMenuItem("Export", func() {
		// Implement export functionality here
	})
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
		*hosts = db.SelectALL(database) // Note: updating the hosts slice
		hostList.Refresh()
		hostList.Select(data.FindHost(id, *hosts))
	})

	// Create the menus
	fileMenu := fyne.NewMenu("File", fileExportItem)
	editMenu := fyne.NewMenu("Edit", editNewHostItem)

	// Return the constructed menu bar
	return fyne.NewMainMenu(
		fileMenu,
		editMenu,
	)
}
