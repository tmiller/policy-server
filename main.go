package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

var policy []byte
var policyFile string
var keyFile string
var certFile string
var bindAddress string
var logFile string

func init() {
	flag.StringVar(&policyFile, "p", "crossdomain.xml", "policy file")
	flag.StringVar(&certFile, "c", "tls.crt", "tls certificate")
	flag.StringVar(&keyFile, "k", "tls.key", "tls private key")
	flag.StringVar(&bindAddress, "b", ":843", "bind address")
	flag.StringVar(&logFile, "l", "", "log file")
	flag.Parse()

	pid := strconv.Itoa(os.Getpid())
	log.SetPrefix(fmt.Sprintf("%v : ", pid))
}

func main() {
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf(": ERROR : %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	policy, err := ioutil.ReadFile(policyFile)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}
	// Append a null byte since that is how the current policy server works
	policy = append(policy, '\x00')

	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		},
	}

	log.Printf(": INFO : starting listener on %v", bindAddress)
	ln, err := tls.Listen("tcp", bindAddress, tlsConfig)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}
	defer ln.Close()

	log.Printf(": INFO : listening on %v", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, policy)
	}
}

func handleConnection(conn net.Conn, policy []byte) {
	defer conn.Close()

	// For nicer formatting post pone writing to the log
	logBuffer := fmt.Sprintf("serving %v", conn.RemoteAddr())

	// Read what was requested for the logs
	var request []byte = make([]byte, 32, 32)
	bytes, err := conn.Read(request)
	if err != nil {
		log.Printf(": ERROR : %s %v", logBuffer, err)
		return
	} else if bytes > 0 {
		logBuffer = fmt.Sprintf(
			"%s requested %v bytes %q",
			logBuffer,
			bytes,
			request[:bytes],
		)
	}

	// Write the policy after reading
	bytes, err = conn.Write(policy)
	if err != nil {
		log.Printf(": ERROR : %s %v", logBuffer, err)
		return
	}

	// Log what happend during the request
	log.Printf(": INFO : %s respond with %v bytes", logBuffer, bytes)
}
