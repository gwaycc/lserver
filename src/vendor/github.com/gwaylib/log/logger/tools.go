package logger

import (
	"os"
)

var HostName = getHost()

func getHost() string {
	hostName, _ := os.Hostname()
	return hostName
}
