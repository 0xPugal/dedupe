# Dedupe

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xpugal/dedupe)](https://goreportcard.com/report/github.com/0xpugal/dedupe)
[![GitHub stars](https://img.shields.io/github/stars/0xpugal/dedupe?style=social)](https://github.com/0xpugal/dedupe/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/0xpugal/dedupe)](https://github.com/0xpugal/dedupe/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/0xpugal/dedupe)](https://github.com/0xpugal/dedupe/pulls)

**Dedupe** is a high-performance, memory-optimized Go tool designed to deduplicate URLs from a list. It supports filtering based on query strings, extensions, and more. Perfect for cleaning up large lists of URLs for web scraping, penetration testing, or bug bounty.

---

## Features

- **Fast and Efficient**: Built in Go for high performance and low memory usage.
- **Deduplication**: Removes duplicate URLs based on hostname and query parameters.
- **Filtering**:
  - Exclude URLs with specific extensions (e.g., `.js`, `.css`, `.png`).
  - Include only URLs with query strings.
  - Remove similar URLs (e.g., `/api/user/1` and `/api/user/2`).
- **Flexible Modes**: Enable multiple filters using the `--mode` flag.
- **Input/Output**:
  - Accepts URLs from a file or stdin.
  - Outputs deduplicated URLs to stdout or a file.

---

## Installation
### Prerequisites
- Go 1.20 or higher.

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
# Deduplicate URLs from a file
./dedupe -u urls.txt

# Only include URLs with query strings
./dedupe -u urls.txt -qs

# Exclude URLs with specific extensions
./dedupe -u urls.txt -ne "js,css,png,jpg"

# Remove similar URLs (e.g., /api/user/1 and /api/user/2)
./dedupe -u urls.txt -s

# Enable multiple modes (query strings, similar URLs, no extensions)
./dedupe -u urls.txt --mode "qs,s,ne"
```

### Help menu
```
Options:
  -h, --help                     Usage/help info for dedupe
  -u, --urls <filename>          Filename containing URLs (use this if you don't pipe URLs via stdin)
  -V, --version                  Get current version for dedupe
  -U, --update                   Check for updates
  -r, --regex-parse              Use regex parsing (slower but more thorough)
  -s, --similar                  Remove similar URLs (based on integers and image/font files)
  -qs, --query-strings-only      Only include URLs if they have query strings
  -ne, --no-extensions ext1,ext2 Do not include URLs with specific extensions
  -m, --mode mode1,mode2         Enable specific modes/filters (r,s,qs,ne)
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
./dedupe -u urls.txt -ne "js,css,png" -qs

https://example.com/api/user/1?name=John
https://example.com/login.asp?redirect=/dashboard
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