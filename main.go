package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	//webAPI
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {

		domain := r.URL.Query().Get("domain")

		if domain == "" {
			fmt.Fprintln(w, "please provide a domain")
			return
		}

		checkDomain(w, domain)
	})

	// index.html
	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Println("Server running at http://localhost:8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain , hasMX , hasSPF , spfRecord, hasDMARC, dmarchRecord\n")

	for scanner.Scan() {
		checkDomain(nil, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v/n", err)
	}
}

func checkDomain(w http.ResponseWriter, domain string) {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarchRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}
	dmarchRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	for _, record := range dmarchRecords {
		if strings.HasPrefix(strings.ToLower(record), "v=dmarc1") {
			hasDMARC = true
			dmarchRecord = record
			break
		}
	}

	if w != nil {
		fmt.Fprintf(w,
			"Domain: %v\nMX: %v\nSPF: %v\nSPF Record: %v\nDMARC: %v\nDMARC Record: %v\n",
			domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarchRecord)
	} else {
		fmt.Fprintf(w,
			"Domain: %v\nMX: %v\nSPF: %v\nSPF Record: %v\nDMARC: %v\nDMARC Record: %v\n",
			domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarchRecord)
	}
}
