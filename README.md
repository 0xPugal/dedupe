# dedupe
`dedupe` is a URL deduplication tool that normalizes and deduplicates URLs based on hostname and query parameters. It offers options to filter URLs by query strings and exclude URLs with specific file extensions.


[![made-with-bash](https://img.shields.io/badge/Made%20with-Bash-1f425f.svg)](https://www.gnu.org/software/bash/) [![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/0xPugal/dedupe/graphs/commit-activity) [![Latest release](https://badgen.net/github/release/0xPugal/dedupe?sort=semver&label=version)](https://github.com/0xPugal/dedupe/releases) [![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/0xPugal/dedupe)

![dedupe](https://github.com/0xPugal/dedupe/assets/75373225/77754d3d-1e9f-4250-a2b7-58c36c60ae0c)

## Installation

```sh
curl -sSL https://raw.githubusercontent.com/0xPugal/dedupe/master/dedupe -o dedupe && chmod +x dedupe && sudo mv dedupe /usr/bin/
```

## Features

- **Normalize and Deduplicate URLs:** Based on hostname and sorted query parameter names.
- **Filter by Query Strings:** Optionally include only URLs with query strings.
- **Exclude Specific Extensions:** Remove URLs with specified file extensions.

## Help/Options

```sh
Usage: dedupe [options] [<input_file>]
Options:
  -h, --help                    Usage/help info for dedupe
  -u, --urls <filename>         Filename containing URLs (use this if you don't pipe URLs via stdin)
  -V, --version                 Get current version for dedupe
  -qs, --query-strings-only     Only include URLs if they have query strings
  -ne, --no-extensions <ext>    Do not include URLs with specific extensions
```

## Usage

- `cat urls.txt | dedupe -qs` or `dedupe -u urls.txt -qs` to get only parameterized URLs
- `cat urls.txt | dedupe -ne css,png,js` or `dedupe -u urls.txt -ne css,png,js` to remove URLs with these extensions
- Chain with other tools `echo example.com | gau | dedupe -qs -ne css,png,jpg,gif | anew output.txt`

**Before:**

```sh
$ cat test.txt
https://test.com/api/users/123
https://test.com/api/users/222
https://test.com/api/users/412/profile
https://test.com/users/photos/photo.jpg
https://test.com/users/photos/myPhoto.jpg
https://demo.com/photo.png
https://google.com/home?qs=fuzz
https://google.com/home?qs=new&second=old
https://google.com/home?qs=asd&xyz=das
https://bing.com/test
https://bing.com/test.php?x=y&y=z
```

**Only URLs with query strings:**

```sh
$ ./dedupe -u test.txt -qs
https://google.com/home?qs=fuzz
https://google.com/home?qs=new&second=old
https://google.com/home?qs=asd&xyz=das
https://bing.com/test.php?x=y&y=z
```

**Remove URLs with certain extensions:**

```sh
$ ./dedupe -u test.txt -ne jpg,png,gif
https://test.com/api/users/123
https://test.com/api/users/222
https://test.com/api/users/412/profile
https://google.com/home?qs=fuzz
https://google.com/home?qs=new&second=old
https://google.com/home?qs=asd&xyz=das
https://bing.com/test
https://bing.com/test.php?x=y&y=z
```

## Similar Tools

- [uro](https://github.com/s0md3v/uro)
- [urless](https://github.com/xnl-h4ck3r/urless)
- [durl](https://github.com/j3ssie/durl)
- [urldedupe](https://github.com/ameenmaali/urldedupe)
```
