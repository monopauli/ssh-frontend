package db

import (
	"database/sql"
	"fmt"
	"frontend/data"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {

	// Load the configuration from the config.toml file
	config, _ := data.LoadConfiguration("config/config.toml")

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func AddEntry(db *sql.DB, host data.Host) (int, error) {
	query := `
        INSERT INTO hosts (host, hostname, username, port, identityfile, systemtype, nodetype, provider, region, internalip, portbase)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `

	var id int
	err := db.QueryRow(
		query,
		host.Host,
		host.Hostname,
		host.User,
		host.Port,
		host.IdentityFile,
		host.SystemType,
		host.NodeType,
		host.Provider,
		host.Region,
		host.InternalIP,
		host.Portbase).Scan(&id)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return id, nil
}

func DeleteEntry(db *sql.DB, host data.Host) {
	query := `
		DELETE FROM hosts WHERE id = $1
		`
	_, err := db.Exec(query, host.ID)

	if err != nil {
		log.Fatal(err)
	}
}

func Update(db *sql.DB, host data.Host, id int) {

	query := `
	UPDATE hosts 
	SET host = $1, hostname = $2, username = $3, port = $4, identityfile = $5, 
		systemtype = $6, nodetype = $7, provider = $8, region = $9, 
		internalip = $10, portbase = $11 
	WHERE id = $12
`
	_, err := db.Exec(
		query,
		host.Host,
		host.Hostname,
		host.User,
		host.Port,
		host.IdentityFile,
		host.SystemType,
		host.NodeType,
		host.Provider,
		host.Region,
		host.InternalIP,
		host.Portbase,
		id)

	if err != nil {
		log.Fatal(err)
	}

}

func SelectALL(db *sql.DB) []data.Host {
	var hostList []data.Host
	rows, err := db.Query(`
    SELECT id, host, hostname, username, port, identityfile, systemtype,
           nodetype, provider, region, internalip, portbase
    FROM hosts
    ORDER BY UPPER(host) ASC, id ASC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		var hostname string
		var username string
		var port int
		var identityfile string
		var systemtype string
		var nodetype string
		var provider string
		var region string
		var internalip string
		var portbase int
		var id int
		if err := rows.Scan(
			&id,
			&host,
			&hostname,
			&username,
			&port,
			&identityfile,
			&systemtype,
			&nodetype,
			&provider,
			&region,
			&internalip,
			&portbase,
		); err != nil {
			log.Fatal(err)
		}
		newhost := data.Host{
			ID:           id,
			Host:         host,
			Hostname:     hostname,
			User:         username,
			Port:         port,
			IdentityFile: identityfile,
			SystemType:   systemtype,
			NodeType:     nodetype,
			Provider:     provider,
			Region:       region,
			InternalIP:   internalip,
			Portbase:     portbase,
		}
		hostList = append(hostList, newhost)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return hostList
}
