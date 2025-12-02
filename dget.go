package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	bases       []string
	tlds        []string
	domainsFile string
	tldsFile    string
	help        bool
}

type DomainResult struct {
	domain    string
	available bool
	index     int
}

func main() {
	var debug = flag.Bool("debug", false, "enable debug output")

	config := parseFlags()

	if config.help {
		printHelp()
		return
	}

	domains := buildDomainList(config)

	if len(domains) == 0 {
		printUsageExamples()
		return
	}

	checkDomainsConcurrent(domains)
	if *debug {
		fmt.Println(domains)
	}
}

func parseFlags() Config {
	var config Config
	var baseFlags, tldFlags arrayFlags

	flag.Var(&baseFlags, "base", "Base domain name(s)")
	flag.Var(&baseFlags, "b", "Base domain name(s) (shorthand)")
	flag.Var(&tldFlags, "tld", "Top-level domain(s)")
	flag.Var(&tldFlags, "t", "Top-level domain(s) (shorthand)")
	flag.StringVar(&config.domainsFile, "domains-file", "", "File containing base domain names")
	flag.StringVar(&config.domainsFile, "df", "", "File containing base domain names (shorthand)")
	flag.StringVar(&config.tldsFile, "tld-file", "", "File containing TLDs")
	flag.StringVar(&config.tldsFile, "tf", "", "File containing TLDs (shorthand)")
	flag.BoolVar(&config.help, "help", false, "Show help message")
	flag.BoolVar(&config.help, "h", false, "Show help message (shorthand)")

	flag.Parse()

	// Handle positional arguments
	args := flag.Args()
	if len(args) > 0 {
		for _, arg := range args {
			if strings.Contains(arg, ".") {
				// Full domain provided
				parts := strings.SplitN(arg, ".", 2)
				config.bases = append(config.bases, parts[0])
				config.tlds = append(config.tlds, parts[1])
			} else {
				// Just base provided
				config.bases = append(config.bases, arg)
			}
		}
	}

	// Merge flag values
	config.bases = append(config.bases, baseFlags...)
	config.tlds = append(config.tlds, tldFlags...)

	return config
}

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlags) Set(value string) error {
	// Split by comma to support comma-separated values
	parts := strings.Split(value, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			*a = append(*a, trimmed)
		}
	}
	return nil
}

func buildDomainList(config Config) []string {
	var bases, tlds []string

	// Collect bases
	if config.domainsFile != "" {
		bases = append(bases, readLinesFromFile(config.domainsFile)...)
	}
	bases = append(bases, config.bases...)

	// Collect TLDs
	if config.tldsFile != "" {
		tlds = append(tlds, readLinesFromFile(config.tldsFile)...)
	}
	tlds = append(tlds, config.tlds...)

	// Remove duplicates
	bases = unique(bases)
	tlds = unique(tlds)

	// Build domain combinations
	var domains []string
	if len(bases) > 0 && len(tlds) > 0 {
		for _, base := range bases {
			for _, tld := range tlds {
				domains = append(domains, fmt.Sprintf("%s.%s", base, tld))
			}
		}
	}

	return unique(domains)
}

func readLinesFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
	}

	return lines
}

func unique(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func checkDomainsConcurrent(domains []string) {
	fmt.Printf("Checking %d domain(s)...\n\n", len(domains))

	// Create channels for work distribution and results
	results := make(chan DomainResult, len(domains))

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Limit concurrent checks to avoid overwhelming the system
	// Adjust this number based on your needs (50 is a reasonable default)
	maxConcurrent := 50
	if len(domains) < maxConcurrent {
		maxConcurrent = len(domains)
	}

	semaphore := make(chan struct{}, maxConcurrent)

	// Launch a goroutine for each domain
	for i, domain := range domains {
		wg.Add(1)
		go func(domain string, index int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			available := isDomainAvailable(domain)
			results <- DomainResult{
				domain:    domain,
				available: available,
				index:     index,
			}
		}(domain, i)
	}

	// Close results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results in order
	resultMap := make(map[int]DomainResult)
	for result := range results {
		resultMap[result.index] = result
	}

	// Print results in original order
	for i := 0; i < len(domains); i++ {
		result := resultMap[i]
		status := "REGISTERED"
		emoji := "❌"
		if result.available {
			status = "AVAILABLE"
			emoji = "✅"
		}
		fmt.Printf("%s %-40s %s\n", emoji, result.domain, status)
	}
}

func isDomainAvailable(domain string) bool {
	// Method 1: Try DNS lookup - if it fails, domain might be available
	_, err := net.LookupHost(domain)
	if err == nil {
		// Domain resolves, definitely registered
		return false
	}

	// Method 2: Try to resolve nameservers
	ns, err := net.LookupNS(domain)
	if err == nil && len(ns) > 0 {
		// Has nameservers, definitely registered
		return false
	}

	// Method 3: Check if WHOIS server responds (simplified check)
	// Try to connect to common WHOIS servers
	whoisServer := getWhoisServer(domain)
	if whoisServer != "" {
		conn, err := net.DialTimeout("tcp", whoisServer+":43", 5*time.Second)
		if err == nil {
			defer conn.Close()

			// Send WHOIS query
			fmt.Fprintf(conn, "%s\r\n", domain)

			// Read response
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			scanner := bufio.NewScanner(conn)
			var response strings.Builder

			for scanner.Scan() {
				response.WriteString(strings.ToLower(scanner.Text()))
				response.WriteString("\n")
			}

			responseText := response.String()

			// Check for common "not found" patterns
			notFoundPatterns := []string{
				"no match",
				"not found",
				"no entries found",
				"no data found",
				"status: available",
				"not registered",
				"no matching record",
				"domain not found",
			}

			for _, pattern := range notFoundPatterns {
				if strings.Contains(responseText, pattern) {
					return true
				}
			}

			// Check for common "registered" patterns
			registeredPatterns := []string{
				"domain name:",
				"registrar:",
				"creation date:",
				"registrant",
			}

			registeredCount := 0
			for _, pattern := range registeredPatterns {
				if strings.Contains(responseText, pattern) {
					registeredCount++
				}
			}

			if registeredCount >= 2 {
				return false
			}
		}
	}

	// If we can't determine definitively, assume registered for safety
	return false
}

func getWhoisServer(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}

	tld := parts[len(parts)-1]

	// Common WHOIS servers by TLD
	whoisServers := map[string]string{
		"com":    "whois.verisign-grs.com",
		"net":    "whois.verisign-grs.com",
		"org":    "whois.pir.org",
		"info":   "whois.afilias.net",
		"biz":    "whois.biz",
		"us":     "whois.nic.us",
		"co":     "whois.nic.co",
		"io":     "whois.nic.io",
		"ai":     "whois.nic.ai",
		"me":     "whois.nic.me",
		"tv":     "whois.nic.tv",
		"cc":     "whois.nic.cc",
		"ws":     "whois.website.ws",
		"name":   "whois.nic.name",
		"mobi":   "whois.dotmobiregistry.net",
		"asia":   "whois.nic.asia",
		"tel":    "whois.nic.tel",
		"travel": "whois.nic.travel",
		"pro":    "whois.registrypro.pro",
		"jobs":   "whois.nic.jobs",
		"cat":    "whois.nic.cat",
		"aero":   "whois.aero",
		"coop":   "whois.nic.coop",
		"museum": "whois.museum",
	}

	if server, ok := whoisServers[tld]; ok {
		return server
	}

	// Try generic WHOIS server
	return fmt.Sprintf("whois.nic.%s", tld)
}

func printUsageExamples() {
	fmt.Println("                                          ")
	fmt.Println("        /$$                       /$$     ")
	fmt.Println("       | $$                      | $$     ")
	fmt.Println("   /$$$$$$$  /$$$$$$   /$$$$$$  /$$$$$$   ")
	fmt.Println("  /$$__  $$ /$$__  $$ /$$__  $$|_  $$_/   ")
	fmt.Println(" | $$  | $$| $$  \\ $$| $$$$$$$$  | $$     ")
	fmt.Println(" | $$  | $$| $$  | $$| $$_____/  | $$ /$$ ")
	fmt.Println(" |  $$$$$$$|  $$$$$$$|  $$$$$$$  |  $$$$/ ")
	fmt.Println("  \\_______/ \\____  $$ \\_______/   \\___/   ")
	fmt.Println("            /$$  \\ $$                     ")
	fmt.Println("           |  $$$$$$/                     ")
	fmt.Println("            \\______/                      ")
	fmt.Println("                                          ")

	fmt.Println("dget — Domain Availability Checker")
	fmt.Println("\nNo domains specified. Here are some usage examples:")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dget example.com")
	fmt.Println("  dget --tld com example")
	fmt.Println("  dget --base example --tld com")
	fmt.Println("  dget --base example,test --tld com,net")
	fmt.Println("  dget --base example --tld com --tld net")
	fmt.Println("  dget --domains-file domains-file.txt --tld com")
	fmt.Println("  dget --domains-file domains-file.txt --tld-file tlds-file.txt")
	fmt.Println()
	fmt.Println("For more information, run: dget --help")
}

func printHelp() {
	fmt.Println("dget — Domain Availability Checker")
	fmt.Println("\nUsage:")
	fmt.Println("  dget [options] [domain]")
	fmt.Println()
	fmt.Println("Options (with shorthand):")
	fmt.Println("  -b, --base <name>          Base domain name (can be comma-separated or used multiple times)")
	fmt.Println("  -t, --tld <tld>            Top-level domain (can be comma-separated or used multiple times)")
	fmt.Println("  -df, --domains-file <file> File containing base domain names (one per line)")
	fmt.Println("  -tf, --tld-file <file>     File containing TLDs (one per line)")
	fmt.Println("  -h, --help                 Show this help message")
	fmt.Println()
	fmt.Println("Examples (and you can use shorter flags):")
	fmt.Println("  dget example.com")
	fmt.Println("    Check if example.com is available")
	fmt.Println()
	fmt.Println("  dget --tld com example")
	fmt.Println("    Check if example.com is available")
	fmt.Println()
	fmt.Println("  dget --base example --tld com,net,org")
	fmt.Println("    Check example.com, example.net, and example.org")
	fmt.Println()
	fmt.Println("  dget --base example,test --tld com --tld net")
	fmt.Println("    Check all combinations: example.com, example.net, test.com, test.net")
	fmt.Println()
	fmt.Println("  dget --domains-file domains-file.txt --tld com")
	fmt.Println("    Check all domains from file with .com TLD")
	fmt.Println()
	fmt.Println("  dget --domains-file domains-file.txt --tld-file tlds-file.txt")
	fmt.Println("    Check all combinations from both files")
	fmt.Println()
	fmt.Println("File Formats:")
	fmt.Println("  domains-file.txt:")
	fmt.Println("    example")
	fmt.Println("    test-domain")
	fmt.Println("    mysite")
	fmt.Println()
	fmt.Println("  tlds-file.txt:")
	fmt.Println("    com")
	fmt.Println("    net")
	fmt.Println("    org")
}
