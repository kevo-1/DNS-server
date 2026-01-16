package data

func GetRootServers() []string {
	servers := make([]string, len(RootServers))
	for i, server := range RootServers {
		servers[i] = server.IPv4
	}
	return servers
}

func GetRootServerByName(name string) *NameServer {
	manager := GetRootServerManager()
	return manager.LookUpByName(name)
}

func GetRandomRootServer() string {
	manager := GetRootServerManager()
	root := manager.GetFirstRoot()
	if root != nil {
		return root.IPv4
	}
	return ""
}
