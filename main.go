package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	config, err := getConfig("solarflare.toml") // Placeholder
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		func() {
			iface, err := net.InterfaceByName(config.Interface)
			if err != nil {
				log.Println(err)
				return
			}

			addresses, err := iface.Addrs()
			if err != nil {
				log.Println(err)
				return
			}

			if len(addresses) == 0 {
				log.Println("Interface has no unicast addresses")
			}

			ipAddr, ok := addresses[0].(*net.IPNet)
			if ok == false {
				log.Println("Wrong addr type")
				return
			}
			log.Printf("IP address: %v\n", ipAddr.String())

			api, err := cloudflare.NewWithAPIToken(config.ApiKey)
			if err != nil {
				log.Println(err)
				return
			}

			zoneID, err := api.ZoneIDByName(config.Zone)
			if err != nil {
				log.Println(err)
				return
			}

			for _, subdomain := range config.Subdomains {
				log.Printf("Updating %v\n", subdomain)

				record := cloudflare.DNSRecord{Name: subdomain}
				records, err := api.DNSRecords(zoneID, record)
				if err != nil {
					log.Println(err)
					continue
				}

				record = cloudflare.DNSRecord{
					Type:    "A",
					Name:    subdomain,
					Content: ipAddr.IP.String(),
					TTL:     1,
				}

				if len(records) == 0 { // Zone doesn't exist, so create it
					_, err = api.CreateDNSRecord(zoneID, record)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println("Success")
					return
				}

				err = api.UpdateDNSRecord(zoneID, records[0].ID, record)

				log.Println("Success")
			}

			log.Printf("Sleeping for %v\n", config.UpdateInterval.Duration)
		}()

		time.Sleep(config.UpdateInterval.Duration)
	}
}
