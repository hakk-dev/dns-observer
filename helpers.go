package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func readZone(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getId(url string) string {
	return strings.Split(url, ".")[0]
}

func checkIP(ip string, ipRange []string) bool {
	ip_net := net.ParseIP(ip)

	for _, r := range ipRange {
		_, ipnetA, _ := net.ParseCIDR(r)
		if ipnetA.Contains(ip_net) {
			return true
		}
	}
	return false
}

func checkProvider(ipAddr string) string {
	for k, v := range rangeMap {
		rv := checkIP(ipAddr, v)
		if rv {
			return k
		}
	}
	return "Unknown"
}

func configSetup() {
	// read in the ranges
	f := map[string]string{"Cloudflare": "iplists/cloudflare_ipv4.txt", 
        "NextDNS": "iplists/nextdns_ipv4.txt"}

	for k, v := range f {
		lines, err := readZone(v)
		if err != nil {
			log.Fatalf("readZone: %s", err)
		}
		rangeMap[k] = lines
	}

	// set config defaults
	viper.SetDefault("dns_addr", "127.0.0.1")
	viper.SetDefault("dns_port", "8053")

	viper.SetDefault("http_addr", "127.0.0.1")
	viper.SetDefault("http_port", "8080")

	viper.SetDefault("responseIP", "127.0.0.1")

	// load a config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("Error: %s. Using defaults.\n", err)
	}

	if viper.Get("dns_addr") != nil {
		dns_addr = viper.GetString("dns_addr")
	}

	if viper.Get("dns_port") != nil {
		dns_port = viper.GetInt("dns_port")
	}

	if viper.Get("http_addr") != nil {
		http_addr = viper.GetString("http_addr")
	}

	if viper.Get("http_port") != nil {
		http_port = viper.GetInt("http_port")
	}

	if viper.Get("responseIP") != nil {
		responseIP = viper.GetString("responseIP")
	}
}
