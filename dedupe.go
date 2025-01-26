package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	helpFlag          bool
	versionFlag       bool
	updateFlag        bool
	urlsFlag          string
	regexParseFlag    bool
	similarFlag       bool
	queryStringsOnly  bool
	noExtensionsFlag  string
	modeFlag          string
)

// ANSI escape codes for colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[97m"
	colorBold   = "\033[1m"
)

// Current version of the tool
const currentVersion = "v0.1"

// URL to fetch the latest version
const versionURL = "https://raw.githubusercontent.com/0xpugal/dedupe/master/VERSION"

func showBanner() {
	boldWhite := colorBold + colorWhite
	boldRed := colorBold + colorRed
	boldCyan := colorBold + colorCyan

	fmt.Println(boldWhite + "______  _______ ______  _     _  _____  _______" + colorReset)
	fmt.Println(boldWhite + "|     \\ |______ |     \\ |     | |_____] |______" + colorReset)
	fmt.Println(boldWhite + "|_____/ |______ |_____/ |_____| |       |______ v0.1" + colorReset)
	fmt.Println()
	fmt.Print(boldWhite + "	 Made with " + colorReset)
	fmt.Print(boldRed + "<3" + colorReset)
	fmt.Print(boldWhite + " and " + colorReset)
	fmt.Print(boldRed + "AI" + colorReset)
	fmt.Print(boldWhite + " by" + colorReset)
	fmt.Print(boldCyan + " @0xpugal" + colorReset)
	fmt.Println()
}

func showHelp() {
	showBanner() // Display banner with help
	//fmt.Println("Usage: dedupe [options] [<input_file>]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help                     Usage/help info for dedupe")
	fmt.Println("  -u, --urls <filename>          Filename containing URLs (use this if you don't pipe URLs via stdin)")
	fmt.Println("  -V, --version                  Get current version for dedupe")
	fmt.Println("  -U, --update                   Check for updates")
	fmt.Println("  -r, --regex-parse              Use regex parsing (slower but more thorough)")
	fmt.Println("  -s, --similar                  Remove similar URLs (based on integers and image/font files)")
	fmt.Println("  -qs, --query-strings-only      Only include URLs if they have query strings")
	fmt.Println("  -ne, --no-extensions ext1,ext2 Do not include URLs with specific extensions")
	fmt.Println("  -m, --mode mode1,mode2         Enable specific modes/filters (r,s,qs,ne)")
}

func showVersion() {
	fmt.Printf("dedupe %s\n", currentVersion)
}

func checkForUpdates() {
	resp, err := http.Get(versionURL)
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching version:", resp.Status)
		return
	}

	var latestVersion string
	scanner := bufio.NewScanner(resp.Body)
	if scanner.Scan() {
		latestVersion = strings.TrimSpace(scanner.Text())
	}

	if latestVersion == "" {
		fmt.Println("Unable to fetch the latest version.")
		return
	}

	if latestVersion == currentVersion {
		fmt.Println("You are using the latest version:", currentVersion)
	} else {
		fmt.Printf("A new version is available: %s (current: %s)\n", latestVersion, currentVersion)
		fmt.Println("Update from: https://github.com/0xpugal/dedupe")
	}
}

func normalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	hostname := parsed.Hostname()
	query := parsed.Query()
	var keys []string
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var normalizedQuery string
	for _, key := range keys {
		if normalizedQuery != "" {
			normalizedQuery += "&"
		}
		normalizedQuery += key
	}

	if normalizedQuery != "" {
		return fmt.Sprintf("%s?%s", hostname, normalizedQuery)
	}
	return hostname
}

func hasExtension(rawURL string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.Contains(rawURL, "."+ext) {
			return true
		}
	}
	return false
}

func isSimilarURL(rawURL string) bool {
	similarPatterns := []string{
		`\d+`, // Integers
		`\.(png|jpg|jpeg|gif|woff|woff2|ttf|otf|svg|ico)$`, // Image/font files
	}
	for _, pattern := range similarPatterns {
		if regexp.MustCompile(pattern).MatchString(rawURL) {
			return true
		}
	}
	return false
}

func deduplicateURLs(inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 64*1024) // 64 KB buffer
	scanner.Buffer(buffer, len(buffer))

	uniqueURLs := make(map[string]bool)
	for scanner.Scan() {
		rawURL := scanner.Text()
		if rawURL == "" {
			continue
		}

		if noExtensionsFlag != "" {
			extensions := strings.Split(noExtensionsFlag, ",")
			if hasExtension(rawURL, extensions) {
				continue
			}
		}

		if queryStringsOnly && !strings.Contains(rawURL, "?") {
			continue
		}

		if similarFlag && isSimilarURL(rawURL) {
			continue
		}

		normalized := normalizeURL(rawURL)
		if _, exists := uniqueURLs[normalized]; !exists {
			uniqueURLs[normalized] = true
			fmt.Println(rawURL)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
}

func main() {
	flag.BoolVar(&helpFlag, "h", false, "Usage/help info for dedupe")
	flag.BoolVar(&helpFlag, "help", false, "Usage/help info for dedupe")
	flag.BoolVar(&versionFlag, "V", false, "Get current version for dedupe")
	flag.BoolVar(&updateFlag, "U", false, "Check for updates")
	flag.StringVar(&urlsFlag, "u", "", "Filename containing URLs (use this if you don't pipe URLs via stdin)")
	flag.StringVar(&urlsFlag, "urls", "", "Filename containing URLs (use this if you don't pipe URLs via stdin)")
	flag.BoolVar(&regexParseFlag, "r", false, "Use regex parsing (slower but more thorough)")
	flag.BoolVar(&regexParseFlag, "regex-parse", false, "Use regex parsing (slower but more thorough)")
	flag.BoolVar(&similarFlag, "s", false, "Remove similar URLs (based on integers and image/font files)")
	flag.BoolVar(&similarFlag, "similar", false, "Remove similar URLs (based on integers and image/font files)")
	flag.BoolVar(&queryStringsOnly, "qs", false, "Only include URLs if they have query strings")
	flag.BoolVar(&queryStringsOnly, "query-strings-only", false, "Only include URLs if they have query strings")
	flag.StringVar(&noExtensionsFlag, "ne", "", "Do not include URLs with specific extensions")
	flag.StringVar(&noExtensionsFlag, "no-extensions", "", "Do not include URLs with specific extensions")
	flag.StringVar(&modeFlag, "m", "", "Enable specific modes/filters (r,s,qs,ne)")
	flag.StringVar(&modeFlag, "mode", "", "Enable specific modes/filters (r,s,qs,ne)")
	flag.Parse()

	if len(os.Args) == 1 {
		showBanner()
		return
	}

	if helpFlag {
		showHelp() // Show banner and help message
		return
	}

	if versionFlag {
		showVersion()
		return
	}

	if updateFlag {
		checkForUpdates()
		return
	}

	if modeFlag != "" {
		modes := strings.Split(modeFlag, ",")
		for _, mode := range modes {
			switch mode {
			case "r":
				regexParseFlag = true
			case "s":
				similarFlag = true
			case "qs":
				queryStringsOnly = true
			case "ne":
				noExtensionsFlag = "png,jpg,jpeg,gif,woff,woff2,ttf,otf,svg,ico"
			}
		}
	}

	inputFile := urlsFlag
	if inputFile == "" {
		inputFile = "/dev/stdin"
	}

	deduplicateURLs(inputFile)
}
