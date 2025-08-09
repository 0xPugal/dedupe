package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "regexp"
    "sort"
    "strings"
    "sync"

    yaml "gopkg.in/yaml.v3"
)

var (
    // base flags
    helpFlag    bool
    versionFlag bool
    updateFlag  bool

    // I/O
    inputFlag  string
    outputFlag string

    // behavior
    queryStringOnlyFlag  bool // -qs
    filterExtensionsFlag string // -fe ext1,ext2 (exclude)
    matchExtensionsFlag  string // -me ext1,ext2 (include-only)
    regexParseFlag       bool   // -r
    langCountryFlag      bool   // -lc
)

// (colors removed for simpler cross-platform output)

// Current version of the tool
const currentVersion = "v0.2"

// URL to fetch the latest version
const versionURL = "https://raw.githubusercontent.com/0xpugal/dedupe/master/VERSION"

// Default language/country codes to normalize
var defaultLanguages = []string{
	"en", "en-us", "en-gb", "fr", "de", "pl", "nl", "fi", "sv", "it", "es", "pt", "ru", "pt-br", "es-mx", "zh-tw", "zh-cn", "ko", "ja", "tr", "ar",
}

// Compiled regex patterns for performance
var (
	guidPattern    = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	integerPattern = regexp.MustCompile(`\b\d+\b`)
	portPattern    = regexp.MustCompile(`:(80|443)(/|$)`)
)

// URLProcessor handles URL processing with thread-safe operations
type URLProcessor struct {
    seenURLs   map[string]struct{}
    mutex      sync.RWMutex
    filterExts map[string]bool
    matchExts  map[string]bool
    languages  map[string]bool
}

// NewURLProcessor creates a new URL processor with optimized memory usage
func NewURLProcessor(filterExts, matchExts, languages []string) *URLProcessor {
    p := &URLProcessor{
        seenURLs:  make(map[string]struct{}, 4096),
        filterExts: make(map[string]bool),
        matchExts:  make(map[string]bool),
        languages:  make(map[string]bool, len(languages)),
    }
    for _, ext := range filterExts {
        if ext = strings.TrimSpace(strings.ToLower(ext)); ext != "" {
            p.filterExts[ext] = true
        }
    }
    for _, ext := range matchExts {
        if ext = strings.TrimSpace(strings.ToLower(ext)); ext != "" {
            p.matchExts[ext] = true
        }
    }
    for _, lang := range languages {
        if lang = strings.TrimSpace(strings.ToLower(lang)); lang != "" {
            p.languages[lang] = true
        }
    }
    return p
}

func showBanner() {
    fmt.Println("dedupe - high performance URL deduplicator")
}

func showHelp() {
    showBanner()
    fmt.Println("Usage:")
    fmt.Println("  cat input.txt | dedupe --output output.txt")
    fmt.Println("  dedupe --input input.txt --output output.txt")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("  -i,  --input <file>            Input file (defaults to stdin)")
    fmt.Println("  -o,  --output <file>           Output file (defaults to stdout)")
    fmt.Println("  -qs, --query-string-only       Only include URLs that have query strings")
    fmt.Println("  -fe, --filter-extensions list  Exclude URLs with these extensions (csv: css,png,js)")
    fmt.Println("  -me, --match-extensions list   Include only URLs with these extensions (csv)")
    fmt.Println("  -r,  --regex-parse             Use regex normalization (GUIDs, integers)")
    fmt.Println("  -lc, --lang-country            Deduplicate by language/country codes")
    fmt.Println("  -V,  --version                 Show version")
    fmt.Println("  -U,  --update                  Check for updates and install latest")
    fmt.Println("  -h,  --help                    Show this help")
}

func showVersion() { fmt.Printf("dedupe %s\n", currentVersion) }

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

func checkForUpdatesAndInstall() {
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
        return
    }

    fmt.Printf("A new version is available: %s (current: %s)\n", latestVersion, currentVersion)
    fmt.Println("Attempting to install latest with: go install github.com/0xpugal/dedupe@latest")
    cmd := exec.Command("go", "install", "github.com/0xpugal/dedupe@latest")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        fmt.Println("Install failed:", err)
        fmt.Println("You can install manually with: go install github.com/0xpugal/dedupe@latest")
        return
    }
    fmt.Println("Installed latest successfully.")
}

// normalizeURL implements urless-like normalization logic
func (p *URLProcessor) normalizeURL(rawURL string) (string, bool) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}

	// Remove ports 80 and 443
	parsed.Host = portPattern.ReplaceAllString(parsed.Host, "$2")

    // Extension filtering
    if len(p.matchExts) > 0 {
        if !p.pathHasAnyExtension(parsed.Path, p.matchExts) {
            return "", false
        }
    } else if len(p.filterExts) > 0 {
        if p.pathHasAnyExtension(parsed.Path, p.filterExts) {
            return "", false
        }
    }

    // Remove trailing slash (except root)
    if strings.HasSuffix(parsed.Path, "/") && len(parsed.Path) > 1 {
        parsed.Path = strings.TrimSuffix(parsed.Path, "/")
    }
    // no param removal; use as-is

	// Query strings only filter
    if queryStringOnlyFlag && len(parsed.Query()) == 0 {
		return "", false
	}

    // drop fragments
    parsed.Fragment = ""

	// Create normalized URL for deduplication
	normalized := p.createNormalizedKey(parsed)
	return normalized, true
}

func (p *URLProcessor) pathHasAnyExtension(path string, extMap map[string]bool) bool {
    lowerPath := strings.ToLower(path)
    for ext := range extMap {
        if strings.HasSuffix(lowerPath, "."+ext) {
            return true
        }
    }
    return false
}

func (p *URLProcessor) createNormalizedKey(parsed *url.URL) string {
    if len(parsed.Query()) > 0 {
        var keys []string
        for key := range parsed.Query() {
            keys = append(keys, key)
        }
        sort.Strings(keys)
        normalizedPath := p.normalizePath(parsed.Path)
        return fmt.Sprintf("%s://%s%s?%s", parsed.Scheme, parsed.Host, normalizedPath, strings.Join(keys, "&"))
    }
    normalizedPath := p.normalizePath(parsed.Path)
    return fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, normalizedPath)
}

func (p *URLProcessor) normalizePath(path string) string {
    if strings.HasSuffix(path, "/") && len(path) > 1 {
        path = strings.TrimSuffix(path, "/")
    }
    if regexParseFlag {
        path = guidPattern.ReplaceAllString(path, "{GUID}")
        path = integerPattern.ReplaceAllString(path, "{INT}")
    }
    if langCountryFlag {
        path = p.normalizeLanguageCodes(path)
    }
    return path
}

func (p *URLProcessor) normalizeLanguageCodes(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if p.languages[strings.ToLower(part)] {
			parts[i] = "{LANG}"
			break // Only replace the first language code
		}
	}
	return strings.Join(parts, "/")
}

func (p *URLProcessor) shouldIncludeURL(normalized string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

    if _, exists := p.seenURLs[normalized]; exists {
		return false
	}

    p.seenURLs[normalized] = struct{}{}
	return true
}

func deduplicateURLs(r io.Reader, w io.Writer, processor *URLProcessor) error {
    scanner := bufio.NewScanner(r)
    scanner.Buffer(make([]byte, 64*1024), 1024*1024)
    out := bufio.NewWriter(w)
    defer out.Flush()

    for scanner.Scan() {
        rawURL := strings.TrimSpace(scanner.Text())
        if rawURL == "" {
            continue
        }
        normalized, ok := processor.normalizeURL(rawURL)
        if !ok {
            continue
        }
        if processor.shouldIncludeURL(normalized) {
            if _, err := out.WriteString(rawURL + "\n"); err != nil {
                return err
            }
        }
    }
    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}

// Config file (optional)
type Config struct {
    QueryStringOnly  bool     `yaml:"query_string_only"`
    FilterExtensions []string `yaml:"filter_extensions"`
    MatchExtensions  []string `yaml:"match_extensions"`
    LangCountry      bool     `yaml:"lang_country"`
}

func loadConfigIfPresent() (*Config, error) {
    f, err := os.Open("config.yml")
    if err != nil {
        // no config present is not an error
        return &Config{}, nil
    }
    defer f.Close()
    data, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }
    var c Config
    if err := yaml.Unmarshal(data, &c); err != nil {
        return nil, err
    }
    return &c, nil
}

func main() {
    flag.BoolVar(&helpFlag, "h", false, "Show help")
    flag.BoolVar(&helpFlag, "help", false, "Show help")
    flag.BoolVar(&versionFlag, "V", false, "Show version")
    flag.BoolVar(&updateFlag, "U", false, "Check for updates and install latest")

    flag.StringVar(&inputFlag, "i", "", "Input file (defaults to stdin)")
    flag.StringVar(&inputFlag, "input", "", "Input file (defaults to stdin)")
    flag.StringVar(&outputFlag, "o", "", "Output file (defaults to stdout)")
    flag.StringVar(&outputFlag, "output", "", "Output file (defaults to stdout)")

    flag.BoolVar(&queryStringOnlyFlag, "qs", false, "Only include URLs that have query strings")
    flag.BoolVar(&queryStringOnlyFlag, "query-string-only", false, "Only include URLs that have query strings")
    flag.StringVar(&filterExtensionsFlag, "fe", "", "Exclude URLs with these extensions (csv)")
    flag.StringVar(&filterExtensionsFlag, "filter-extensions", "", "Exclude URLs with these extensions (csv)")
    flag.StringVar(&matchExtensionsFlag, "me", "", "Include only URLs with these extensions (csv)")
    flag.StringVar(&matchExtensionsFlag, "match-extensions", "", "Include only URLs with these extensions (csv)")
    flag.BoolVar(&regexParseFlag, "r", false, "Enable regex-based normalization (GUIDs, integers)")
    flag.BoolVar(&regexParseFlag, "regex-parse", false, "Enable regex-based normalization (GUIDs, integers)")
    flag.BoolVar(&langCountryFlag, "lc", false, "Normalize language/country codes in paths")
    flag.BoolVar(&langCountryFlag, "lang-country", false, "Normalize language/country codes in paths")

    flag.Parse()

    if helpFlag || len(os.Args) == 1 {
        showHelp()
        return
    }
    if versionFlag {
        showVersion()
        return
    }
    if updateFlag {
        checkForUpdatesAndInstall()
        return
    }

    // load optional config
    cfg, err := loadConfigIfPresent()
    if err != nil {
        fmt.Println("Error loading config.yml:", err)
        return
    }
    if cfg.QueryStringOnly && !queryStringOnlyFlag {
        queryStringOnlyFlag = true
    }
    if cfg.LangCountry && !langCountryFlag {
        langCountryFlag = true
    }
    filterExts := []string{}
    matchExts := []string{}
    if len(cfg.FilterExtensions) > 0 {
        filterExts = append(filterExts, cfg.FilterExtensions...)
    }
    if len(cfg.MatchExtensions) > 0 {
        matchExts = append(matchExts, cfg.MatchExtensions...)
    }
    if filterExtensionsFlag != "" {
        filterExts = strings.Split(filterExtensionsFlag, ",")
    }
    if matchExtensionsFlag != "" {
        matchExts = strings.Split(matchExtensionsFlag, ",")
    }

    processor := NewURLProcessor(filterExts, matchExts, defaultLanguages)

    // input selection
    var in io.ReadCloser
    if inputFlag != "" {
        f, err := os.Open(inputFlag)
        if err != nil {
            fmt.Println("Error opening input file:", err)
            return
        }
        in = f
    } else {
        fi, _ := os.Stdin.Stat()
        if (fi.Mode() & os.ModeCharDevice) == 0 {
            in = os.Stdin
        } else {
            showHelp()
            return
        }
    }
    defer func() { if in != os.Stdin { _ = in.Close() } }()

    // output selection
    var out io.WriteCloser
    if outputFlag != "" {
        f, err := os.Create(outputFlag)
        if err != nil {
            fmt.Println("Error creating output file:", err)
            return
        }
        out = f
    } else {
        out = os.Stdout
    }
    if out != os.Stdout {
        defer out.Close()
    }

    if err := deduplicateURLs(in, out, processor); err != nil {
        fmt.Println("Error:", err)
    }
}
