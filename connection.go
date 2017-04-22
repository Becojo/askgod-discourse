package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (s *syncer) getClient(server string, serverCert string) (*http.Client, error) {
	// Parse the server URL
	u, err := url.ParseRequestURI(server)
	if err != nil {
		return nil, err
	}

	var transport *http.Transport
	if u.Scheme == "http" {
		// Basic transport for clear-text HTTP
		transport = &http.Transport{
			DisableKeepAlives: true,
		}
	} else if u.Scheme == "https" {
		// Be picky on our cipher list
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
			PreferServerCipherSuites: true,
		}

		// If provided, pin the certificate
		if serverCert != "" {
			certBlock, _ := pem.Decode([]byte(serverCert))
			if certBlock == nil {
				return nil, fmt.Errorf("Failed to load pinned certificate")
			}

			cert, err := x509.ParseCertificate(certBlock.Bytes)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse pinned certificate: %v", err)
			}

			caCertPool := tlsConfig.RootCAs
			if caCertPool == nil {
				caCertPool = x509.NewCertPool()
			}

			caCertPool.AddCert(cert)
			tlsConfig.RootCAs = caCertPool
		}

		transport = &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true,
		}
	} else {
		return nil, fmt.Errorf("Unsupported server URL: %s", server)
	}

	// Create the new HTTP client
	client := http.Client{
		Transport: transport,
	}

	return &client, nil
}

func (s *syncer) queryStruct(server string, method string, path string, data interface{}, target interface{}) error {
	var req *http.Request
	var err error

	// Server-specific configuration
	var srv *http.Client
	var url string
	if server == "askgod" {
		srv = s.httpAskgod
		url = fmt.Sprintf("%s/1.0%s", s.config.AskgodURL, path)
	} else if server == "discourse" {
		srv = s.httpDiscourse
		url = fmt.Sprintf("%s%s?api_key=%s&api_username=%s", s.config.DiscourseURL, path, s.config.DiscourseAPIKey, s.config.DiscourseAPIUser)
	} else {
		return fmt.Errorf("Unknown server: %s", server)
	}

	// Get a new HTTP request setup
	if data != nil {
		// Encode the provided data
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(data)
		if err != nil {
			return err
		}

		// Some data to be sent along with the request
		req, err = http.NewRequest(method, url, &buf)
		if err != nil {
			return err
		}

		// Set the encoding accordingly
		req.Header.Set("Content-Type", "application/json")
	} else {
		// No data to be sent along with the request
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
	}

	// Send the request
	resp, err := srv.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err == nil && string(content) != "" {
			return fmt.Errorf("%s", strings.TrimSpace(string(content)))
		}

		return fmt.Errorf("%s: %s", url, resp.Status)
	}

	// Decode the response
	if target != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&target)
		if err != nil {
			return err
		}
	}

	return nil
}
