# Dedupe

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xpugal/dedupe)](https://goreportcard.com/report/github.com/0xpugal/dedupe)
[![GitHub stars](https://img.shields.io/github/stars/0xpugal/dedupe?style=social)](https://github.com/0xpugal/dedupe/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/0xpugal/dedupe)](https://github.com/0xpugal/dedupe/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/0xpugal/dedupe)](https://github.com/0xpugal/dedupe/pulls)

![dedupe](https://github.com/user-attachments/assets/e6c6af2f-a9d2-4742-9884-11f0c7f7cbf7)

**Dedupe** is a high-performance, memory-optimized Go tool designed to deduplicate URLs from a list. It supports filtering based on query strings, extensions, and more. Perfect for cleaning up large lists of URLs for web scraping, penetration testing, or bug bounty.

---

## Features

- **Fast and Efficient**: Built in Go for high performance and low memory usage.
- **Deduplication**: Removes duplicate URLs based on hostname, path (normalized), and query parameter keys.
- **Filtering**:
  - Exclude URLs with specific extensions (`-fe css,png,js`) or include only certain extensions (`-me php,js`).
  - Include only URLs with query strings (`-qs`).
  - Optional regex-based normalization (`-r`) treats integers and GUIDs in paths as placeholders, effectively deduplicating similar URLs like `/api/user/1` vs `/api/user/2`.
  - Optional language/country code normalization in paths (`-lc`).
- **Combine flags freely**: Toggle any combination without a separate `--mode`.
- **Input/Output**:
  - Accepts URLs from a file (`-i`) or stdin.
  - Writes to stdout or to a file via `-o`.

---

## Installation
### Prerequisites
- Go 1.23 or higher.

### Install from Source
```bash
git clone https://github.com/0xpugal/dedupe.git
cd dedupe
go build -o dedupe
```

### Install via `go install`
```bash
go install github.com/0xpugal/dedupe@latest
```

---

## Usage

### Basic Usage
```bash
# stdin → stdout
cat urls.txt | dedupe -qs

# stdin → file (Linux/Mac)
cat urls.txt | dedupe -qs -o output.txt

# stdin → file (Windows PowerShell)
type urls.txt | .\dedupe.exe -qs -o output.txt

# file → file
dedupe -i urls.txt -o output.txt -qs

# exclude extensions
dedupe -i urls.txt -o output.txt -fe "js,css,png,jpg"

# include only these extensions
dedupe -i urls.txt -o output.txt -me "php,asp"

# normalize integers/GUIDs in paths (treat /1 and /2 as same)
dedupe -i urls.txt -o output.txt -r

# also normalize language/country codes in paths
dedupe -i urls.txt -o output.txt -r -lc
```

### Help menu
```
Usage:
  cat input.txt | dedupe --output output.txt
  dedupe --input input.txt --output output.txt

Options:
  -i,  --input <file>            Input file (defaults to stdin)
  -o,  --output <file>           Output file (defaults to stdout)
  -qs, --query-string-only       Only include URLs that have query strings
  -fe, --filter-extensions list  Exclude URLs with these extensions (css,png,js,jpg)
  -me, --match-extensions list   Include only URLs with these extensions (php,aspx,jsp)
  -r,  --regex-parse             Use regex normalization (GUIDs, integers)
  -lc, --lang-country            Deduplicate by language/country codes
  -V,  --version                 Show version
  -U,  --update                  Check for updates and install latest
  -h,  --help                    Show this help
```

---

## 📋 Example

### Input (`urls.txt`)
```
https://example.com/api/user/1
https://example.com/api/user/2
https://example.com/api/user/1?name=John
https://example.com/static/js/main.js
https://example.com/static/css/style.css
https://example.com/images/logo.png
https://example.com/index.php
https://example.com/login.asp?redirect=/dashboard
```

### Command
```bash
dedupe -i urls.txt -fe "js,css,png" -qs

https://example.com/api/user/1?name=John
https://example.com/login.asp?redirect=/dashboard
```

---

## Configuration file (optional)

Place a `config.yml` next to the binary or in the working directory to set defaults. CLI flags override config values.

Example `config.yml`:

```yaml
query_string_only: true
filter_extensions: [js, css, png]
match_extensions: []
lang_country: true
```
---

### Contribute
Contributions are welcome! Please open an issue or submit a pull request.

---

## 📄 License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Similar tools
- [uro](https://github.com/s0md3v/uro)
- [urless](https://github.com/xnl-h4ck3r/urless)
- [durl](https://github.com/j3ssie/durl)
- [urldedupe](https://github.com/ameenmaali/urldedupe)

---

## Support
If you find this project useful, please give it a star on [GitHub](https://github.com/0xpugal/dedupe)!
