package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UnixGet : Perform a get request an return the body
func UnixGet(url string) (string, error) {
	machinePath := filepath.Join(os.Getenv("DOCKER_CERT_PATH"))
	machineIp := strings.Trim(os.Getenv("DOCKER_HOST"), "tcp://")
	certFile := filepath.Join(machinePath, "cert.pem")
	keyFile := filepath.Join(machinePath, "key.pem")
	caFile := filepath.Join(machinePath, "ca.pem")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get(fmt.Sprintf("https://%s/%s", machineIp, url))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unixget: request returned %d when expected 200\n%s", resp.StatusCode, data)
	}

	return string(data), err
}
