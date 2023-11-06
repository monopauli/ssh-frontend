package data

import (
	"bufio"
	"os"
	"strings"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfiguration(file string) (DBConfig, error) {
	config := DBConfig{}
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") { // skip blank lines and comments
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // skip lines that don't have a key and value
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "host":
			config.Host = value
		case "port":
			config.Port = value
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "dbname":
			config.DBName = value
		case "sslmode":
			config.SSLMode = value
		}
	}

	if err := scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}

// Then use the LoadConfiguration function as before in your Connect function.
