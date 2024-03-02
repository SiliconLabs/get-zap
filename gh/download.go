/*
Copyright © 2024 Silicon Labs
*/
package gh

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type DownloadOptions struct {
	skipCertCheck  bool
	proxyUrl       *url.URL
	allowHttp      bool
	showPercentage bool
}

func (dso *DownloadOptions) SetProxy(proxyS string) error {
	url, err := url.Parse(proxyS)
	if err != nil {
		return err
	}
	dso.proxyUrl = url
	return nil
}

func (dso *DownloadOptions) SetSkipCertCheck(skipCertCheck bool) {
	dso.skipCertCheck = skipCertCheck
}

func (dso *DownloadOptions) SetAllowHttp(allowHttp bool) {
	dso.allowHttp = allowHttp
}

func (dso *DownloadOptions) SetShowPercentage(showPercentage bool) {
	dso.showPercentage = showPercentage
}

// Returns the default security options.
func DefaultSecurityOptions() *DownloadOptions {
	s := DownloadOptions{
		skipCertCheck:  false,
		proxyUrl:       nil,
		allowHttp:      false,
		showPercentage: true,
	}
	return &s
}

// This function downloads a file from a given URL and puts it into the
// destination path.
func DownloadFileFromUrl(urlAsString string, destinationPath string, sec *DownloadOptions) error {

	tlsConfig := &tls.Config{}
	if sec.skipCertCheck {
		tlsConfig.InsecureSkipVerify = sec.skipCertCheck
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	if sec.proxyUrl != nil {
		tr.Proxy = http.ProxyURL(sec.proxyUrl)
	}

	client := &http.Client{Transport: tr}

	u, err := url.Parse(urlAsString)
	if err != nil {
		return err
	}

	if !sec.allowHttp && u.Scheme == "http" {
		return fmt.Errorf("only secure encrypted HTTPS protocol is allowed, downloads via HTTP are blocked: %v", urlAsString)
	}

	// Security alert: Let's do an actual get now
	response, err := client.Get(urlAsString)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %v", response.StatusCode)
	}

	defer response.Body.Close()

	len := response.ContentLength
	fmt.Printf("Downloading %v bytes to %v ...\n", len, destinationPath)

	output, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer output.Close()

	var chunk int64
	if len > 100 {
		chunk = len / 100
	} else {
		chunk = len
	}
	cnt := 0
	var totalDownloaded int64 = 0
	for {
		written, err := io.CopyN(output, response.Body, chunk)
		totalDownloaded += written
		percentage := (100 * totalDownloaded) / len

		if err == io.EOF {
			fmt.Printf("%v%%: Downloaded %v out of %v bytes. Done!\n", percentage, totalDownloaded, len)
			break
			// done.
		} else if err != nil {
			// real error.
			return err
		} else {
			if sec.showPercentage {
				fmt.Printf("%v%%: Downloaded %v out of %v bytes...\r", percentage, totalDownloaded, len)
			}
			cnt++
		}
	}
	return nil
}
