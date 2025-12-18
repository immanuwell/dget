# dget — Domain Availability Checker


## Features:

- it's really fast ⚡️ (and you can make it even faster by changing `maxConcurrent` (which is 50 by default))
- it uses 3 methods to check is this domain already registered (= 3x more reliable)


## Usage

```bash
        /$$                       /$$
       | $$                      | $$
   /$$$$$$$  /$$$$$$   /$$$$$$  /$$$$$$
  /$$__  $$ /$$__  $$ /$$__  $$|_  $$_/
 | $$  | $$| $$  \ $$| $$$$$$$$  | $$
 | $$  | $$| $$  | $$| $$_____/  | $$ /$$
 |  $$$$$$$|  $$$$$$$|  $$$$$$$  |  $$$$/
  \_______/ \____  $$ \_______/   \___/
            /$$  \ $$
           |  $$$$$$/
            \______/

dget [options] [domain]

Options (with shorthand):
  -b, --base <name>          Base domain name (can be comma-separated or used multiple times)
  -t, --tld <tld>            Top-level domain (can be comma-separated or used multiple times)
  -df, --domains-file <file> File containing base domain names (one per line)
  -tf, --tld-file <file>     File containing TLDs (one per line)
  -h, --help                 Show this help message

Examples:
  dget example.com
  
  dget --tld com example
  
  dget --base example --tld com
  
  dget --base example,test --tld com,net
  
  dget --base example --tld com --tld net
  
  dget --domains-file domains-file.txt --tld com
  
  dget --domains-file domains-file.txt --tld-file tlds-file.txt

File Formats:
  domains-file.txt:
    example
    test-domain
    mysite
    ...

  tlds-file.txt:
    com
    net
    org
    ...

```


## Install or download compiled binary

### Simply install with `go install`

```bash
go install github.com/immanuwell/dget@latest
```

### Or you can download compiled binary

Choose your OS and your architecture

**Linux:**
- [amd64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)
- [arm64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)

**MacOS:**
- [amd64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)
- [arm64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)

**Windows:**
- [amd64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)
- [arm64](https://github.com/immanuwell/dget/releases/tag/v1.0.0)


## Example of usage

![](media/screenshot-1.png)

![](media/screenshot-2.png)

![](media/screenshot-3.png)

![](media/screenshot-4.png)

![](media/dget_screencast.gif)


## You can support me, if you want)

<a href="https://www.buymeacoffee.com/immanuwell" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>



<!-- 

# 64-bit amd64
GOOS=linux GOARCH=amd64 go build -o binaries/linux/amd64/dget

# ARM64
GOOS=linux GOARCH=arm64 go build -o binaries/linux/arm64/dget


# 64-bit amd64
GOOS=windows GOARCH=amd64 go build -o binaries/windows/amd64/dget

# Windows on ARM64 (if supported)
GOOS=windows GOARCH=arm64 go build -o binaries/windows/arm64/dget


# 64-bit Intel
GOOS=darwin GOARCH=amd64 go build -o binaries/macos/amd64/dget

# ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o binaries/macos/arm64/dget

 -->
