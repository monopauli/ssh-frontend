package data

import (
	"os"

	"github.com/kevinburke/ssh_config"
)

type Host struct {
	Host         string
	Name         string
	User         string
	Port         string
	Key          string
	SystemType   string
	NodeType     string
	NodeNetworks string
	Provider     string
	Region       string
	InternalIP   string
	Portbase     string
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
			port, _ := cfg.Get(node.String(), "Port")
			key, _ := cfg.Get(node.String(), "IdentityFile")
			systemtype, _ := cfg.Get(node.String(), "SystemType")
			nodetype, _ := cfg.Get(node.String(), "NodeType")
			nodenetworks, _ := cfg.Get(node.String(), "NodeNetworks")
			provider, _ := cfg.Get(node.String(), "Provider")
			region, _ := cfg.Get(node.String(), "Region")
			internalIP, _ := cfg.Get(node.String(), "InternalIP")
			portbase, _ := cfg.Get(node.String(), "Portbase")

			newHost := Host{
				Host:         hostName,
				Name:         name,
				User:         user,
				Port:         port,
				Key:          key,
				SystemType:   systemtype,
				NodeType:     nodetype,
				NodeNetworks: nodenetworks,
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
