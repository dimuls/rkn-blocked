package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/sirupsen/logrus"
	"github.com/someanon/rkn-blocked/blocked"
)

func main() {

	if len(os.Args) == 1 {
		logrus.Fatal("Expected argument domain or IP address")
	}

	ips, err := blocked.IPs()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load blocked IPs")
	}

	logrus.Infof("Loaded %d blocked IPs", ips.Size())

	domains, err := blocked.Domains()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load blocked domains")
	}

	logrus.Infof("Loaded %d blocked domains", domains.Size())

	color.Blue("Result:\n")

	for _, p := range os.Args[1:] {
		if ips.Has(p) || domains.Has(p) {
			fmt.Print(p, ": ")
			color.Red("blocked")
		} else {
			fmt.Print(p, ": ")
			color.Green("not blocked")
		}
	}

}
