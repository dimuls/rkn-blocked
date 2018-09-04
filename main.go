package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/someanon/rkn-blocked/blocked"
	"golang.org/x/net/proxy"
	"gopkg.in/tucnak/telebot.v2"
)

// RKN bypasser proxy address
const proxyAddr = "127.0.0.1:8000"

func main() {

	tokenBytes, err := ioutil.ReadFile("token")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load telegram bot's token")
	}

	token := string(tokenBytes)
	strings.TrimSpace(token)

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

	go func() {
		select {

		case <-time.After(24 * time.Hour):

			newIPs, err := blocked.IPs()
			if err != nil {
				logrus.WithError(err).Error("Failed to load blocked IPs")
			} else {
				ips = newIPs
				logrus.Infof("Loaded %d blocked IPs", ips.Size())
			}

			newDomains, err := blocked.Domains()
			if err != nil {
				logrus.WithError(err).Error("Failed to load blocked domains")
			} else {
				domains = newDomains
				logrus.Infof("Loaded %d blocked domains", domains.Size())
			}
		}
	}()

	// We are using rkn-bypasser proxy server to bypass telegram bot server
	// block by RKN.
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to dial RKN bypasser SOCKS5 proxy")
		os.Exit(1)
	}

	// Setup proxy dialer.
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial

	b, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Second},
		Client: httpClient,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/blocked", func(m *telebot.Message) {
		p := strings.TrimSpace(m.Payload)

		logrus.WithFields(logrus.Fields{
			"payload": p,
		}).Info("/blocked request")

		if ips.Has(p) || domains.Has(p) {
			b.Send(m.Sender, "blocked")
		} else {
			b.Send(m.Sender, "not blocked")
		}
	})

	logrus.Info("Starting telegram bot server")

	b.Start()

}
