package data

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/kevinburke/ssh_config"
)

type Host struct {
	ID           int
	Host         string
	Hostname     string
	User         string
	Port         int
	IdentityFile string
	SystemType   string
	NodeType     string
	Provider     string
	Region       string
	InternalIP   string
	Portbase     int
	Networks     []string
}

type HostNetwork struct {
	ID      int
	Network string
}

func OpenConfig(path string) *ssh_config.Config {
	f, _ := os.Open(path)
	cfg, _ := ssh_config.Decode(f)
	return cfg
}

func GetHosts(cfg ssh_config.Config) []Host {
	hosts := []Host{}
	for _, host := range cfg.Hosts {
		for _, node := range host.Patterns {
			if node.String() == "*" {
				continue
			}
			hostName := node.String()
			name, _ := cfg.Get(node.String(), "HostName")
			user, _ := cfg.Get(node.String(), "User")
			portstring, _ := cfg.Get(node.String(), "Port")
			port, _ := strconv.Atoi(portstring)
			key, _ := cfg.Get(node.String(), "IdentityFile")
			systemtype, _ := cfg.Get(node.String(), "SystemType")
			nodetype, _ := cfg.Get(node.String(), "NodeType")
			provider, _ := cfg.Get(node.String(), "Provider")
			region, _ := cfg.Get(node.String(), "Region")
			internalIP, _ := cfg.Get(node.String(), "InternalIP")
			portbasestring, _ := cfg.Get(node.String(), "Portbase")
			portbase, _ := strconv.Atoi(portbasestring)
			newHost := Host{
				Host:         hostName,
				Hostname:     name,
				User:         user,
				Port:         port,
				IdentityFile: key,
				SystemType:   systemtype,
				NodeType:     nodetype,
				Provider:     provider,
				Region:       region,
				InternalIP:   internalIP,
				Portbase:     portbase,
			}
			hosts = append(hosts, newHost)
		}
	}
	return hosts[:]
}

func CompareStructs(a, b interface{}) [][]string {
	var changeList [][]string
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	for i := 0; i < aVal.NumField(); i++ {
		aField := aVal.Field(i)
		bField := bVal.Field(i)
		fieldName := aVal.Type().Field(i).Name

		// If it's a slice, we need to handle it differently since != doesn't work on slices
		if aField.Kind() == reflect.Slice {
			if !reflect.DeepEqual(aField.Interface(), bField.Interface()) {
				changeList = append(changeList, []string{fieldName, fmt.Sprintf("%v", aField.Interface()), fmt.Sprintf("%v", bField.Interface())})
			}
		} else {
			// Compare field values for non-slice types
			if aField.Interface() != bField.Interface() {
				changeList = append(changeList, []string{fieldName, fmt.Sprintf("%v", aField.Interface()), fmt.Sprintf("%v", bField.Interface())})
			}
		}
	}
	return changeList
}

func FindHost(id int, hosts []Host) int {
	for i, host := range hosts {
		if host.ID == id {
			return i
		}
	}
	return -1
}
