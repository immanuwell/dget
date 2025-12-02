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
- [amd64](https://github.com/immanuwell/dget/tree/main/binaries/linux/amd64/dget
)
- [arm64](https://github.com/immanuwell/dget/tree/main/binaries/linux/arm64/dget
)

**MacOS:**
- [amd64](https://github.com/immanuwell/dget/tree/main/binaries/windows/amd64/dget
)
- [arm64](https://github.com/immanuwell/dget/tree/main/binaries/windows/arm64/dget
)

**Windows:**
- [amd64](https://github.com/immanuwell/dget/tree/main/binaries/macos/amd64/dget
)
- [arm64](https://github.com/immanuwell/dget/tree/main/binaries/macos/arm64/dget
)


## Example of usage

![](media/screenshot-1.png)

![](media/screenshot-2.png)

![](media/screenshot-3.png)

![](media/screenshot-4.png)

![](media/dget_screencast.gif)


## You can support me, if you want)

<a href="https://www.buymeacoffee.com/immanuwell" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>




