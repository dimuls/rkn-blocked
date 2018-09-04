package blocked

import (
	"encoding/json"
	"errors"
	"github.com/fatih/set"
	"net/http"
	"time"
)

func IPs() (set.Interface, error) {
	req, err := http.NewRequest("GET", "https://reestr.rublacklist.net/api/v2/ips/json", nil)
	if err != nil {
		return nil, err
	}

	c := http.Client{
		Timeout: time.Minute,
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not OK response status: " + res.Status)
	}

		blockedIPs := set.New(set.NonThreadSafe)

	var ips []string

	err = json.NewDecoder(res.Body).Decode(&ips)
	if err != nil {
		return nil, errors.New("failed to decode IPs' JSON: " + err.Error())
	}

	for _, ip := range ips {
		blockedIPs.Add(ip)
	}

	return blockedIPs, nil
}


func Domains() (set.Interface, error) {
	req, err := http.NewRequest("GET", "https://reestr.rublacklist.net/api/v2/domains/json", nil)
	if err != nil {
		return nil, err
	}

	c := http.Client{
		Timeout: time.Minute,
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not OK response status: " + res.Status)
	}

	blockedDomains := set.New(set.NonThreadSafe)

	var domains []string

	err = json.NewDecoder(res.Body).Decode(&domains)
	if err != nil {
		return nil, errors.New("failed to decode domains' JSON: " + err.Error())
	}

	for _, ip := range domains {
		blockedDomains.Add(ip)
	}

	return blockedDomains, nil
}