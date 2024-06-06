# dedupe
A simple bash script to de-duplicate the urls based on hostnames and parameters.

## Installation
```
curl -sSL https://raw.githubusercontent.com/0xPugal/dedupe/master/dedupe -o dedupe && chmod +x dedupe && sudo mv dedupe /usr/bin/
```

## Usage
+ ``dedupe input.txt`` - Specify the input file
+ ``cat input.txt | dedupe`` - Read input as stdin
