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

var VERSION string = ""
var BUILD string = ""

var policy []byte

// Variables read in from the CLI flags
var policyFile string
var keyFile string
var certFile string
var bindAddress string
var logFile string
var numWorkers int
var queueSize int

var tlsCipherSuites = map[uint16]string{
	tls.TLS_RSA_WITH_RC4_128_SHA:                "TLS_RSA_WITH_RC4_128_SHA",
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_RSA_WITH_AES_128_CBC_SHA:            "TLS_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_RSA_WITH_AES_256_CBC_SHA:            "TLS_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256:         "TLS_RSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256:         "TLS_RSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384:         "TLS_RSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:        "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:          "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:   "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384: "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:    "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:  "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
	tls.TLS_FALLBACK_SCSV:                       "TLS_FALLBACK_SCSV",
}

var tlsProtocols = map[uint16]string{
	tls.VersionTLS10: "TLSv1.0",
	tls.VersionTLS11: "TLSv1.1",
	tls.VersionTLS12: "TLSv1.2",
}

func init() {

	// Setup CLI flags and map them to Variables
	flag.StringVar(&policyFile, "p", "crossdomain.xml", "policy file")
	flag.StringVar(&certFile, "c", "localhost.crt", "tls certificate")
	flag.StringVar(&keyFile, "k", "localhost.key", "tls private key")
	flag.StringVar(&bindAddress, "b", ":843", "bind address")
	flag.StringVar(&logFile, "l", "", "log file (default STDOUT)")
	flag.IntVar(&numWorkers, "w", 1, "number of workers")
	flag.IntVar(&queueSize, "q", 0, "size of queue for workers (0 is ok)")
	flag.Parse()

	// Get the PID of the current process
	pid := strconv.Itoa(os.Getpid())

	// Prepend the PID to all logs
	log.SetPrefix(fmt.Sprintf("%v : ", pid))
}

func main() {

	// If a log file was specified on the cli then write output to the log,
	// otherwise it will default to STDOUT
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf(": ERROR : %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// Log the release version and commit sha used to build this version
	log.Printf(": BOOT : version: %v commit: %v", VERSION, BUILD)

	// Load the policy file into memory to be served to clients
	policy, err := ioutil.ReadFile(policyFile)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}
	// Append a null byte since that is how the current policy server works
	policy = append(policy, '\x00')

	// Load a certificate and key for TLS encryption
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}

	// Create a queue and some workers. When the queue is 0 it blocks the main
	// thread from writing to the queue and the workers are blocked until there
	// is something to read. With a queue greater than zero, the main thread
	// blocks when the queue is full and the workers block when the queue is
	// empty.
	log.Printf(
		": BOOT : %v workers consuming queue of %v",
		numWorkers,
		queueSize,
	)
	conns := make(chan net.Conn, queueSize)
	for i := 0; i < numWorkers; i++ {
		go worker(conns, policy, i)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS10,
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

	// Setup TLS secured TCP listener
	log.Printf(": BOOT : starting listener on %v", bindAddress)
	ln, err := tls.Listen("tcp", bindAddress, tlsConfig)
	if err != nil {
		log.Fatalf(": ERROR : %v", err)
	}
	defer ln.Close()

	// Accept connections on the listener and pass them into the queue. Where
	// they will be picked up by the workers.
	log.Printf(": BOOT : listening on %v", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// This puts a connection onto the queue (using a go channel).
		conns <- conn
	}
}

// Read incomming connections and handle them. This also passes along
// the worker id and the number of items in the queue for logging purposes.
func worker(conns <-chan net.Conn, policy []byte, worker int) {
	for conn := range conns {
		handleConnection(conn, policy, worker, len(conns))
	}
}

func handleConnection(conn net.Conn, policy []byte, worker, waiting int) {
	defer conn.Close()

	// For nicer formatting post pone writing to the log
	logBuffer := fmt.Sprintf(
		"waiting %v : worker %v serving %v",
		waiting,
		worker,
		conn.RemoteAddr(),
	)

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
	tlsconn, ok := conn.(*tls.Conn)
	if ok {
		log.Printf(": INFO : %s respond with %v bytes %v %v",
			logBuffer,
			bytes,
			tlsProtocols[tlsconn.ConnectionState().Version],
			tlsCipherSuites[tlsconn.ConnectionState().CipherSuite],
		)
	} else {
		log.Printf(": INFO : %s respond with %v bytes", logBuffer, bytes)
	}
}
