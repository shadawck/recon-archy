# recon-archy
Linkedin Tools (and maybe later other source) to reconstruct a company hierarchy from scraping relations and jobs title

## Dependencies 
First, make sure you have xdfb and openjdk-11-jre installed : 
```
sudo apt-get install xvfb openjdk-11-jre
```

If you don't have ``openjdk-11-jre`` for your system with a package manager, just download it manually from [AdoptOpenJDK](https://adoptopenjdk.net/releases.html):
```
wget https://github.com/AdoptOpenJDK/openjdk11-binaries/releases/download/jdk-11.0.7%2B10/OpenJDK11U-jdk_x64_linux_hotspot_11.0.7_10.tar.gz
sudo tar -xvf OpenJDK11U-jdk_x64_linux_hotspot_11.0.7_10.tar.gz -C /usr/
# and add java to your path 
echo "export PATH=$PATH:/usr/jdk:/usr/jdk/bin:/usr/jdk/lib/:/usr/jdk/jre:/usr/jdk/jre/bin/:/usr/jdk/jre/lib/" >> ~/.bashrc
source ~/.bashrc
```
Then if you don't have it, download firefox
```
sudo apt install firefox
```

## Requirements
### For an Installation **with** golang

ReconArchy need golang if you want to install it with go. If go is not install on your system refer to [golang documentation](https://golang.org/doc/install) to install it. Then go to [Installation with golang](##installation)

### For an installation **without** Golang

Golang need Geckodriver (the WebDriver for firefox) and a selenium server.   
So download the last version of [Selenium Server (Grid)](https://www.selenium.dev/downloads/)
```sh
wget https://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar
mv selenium-server-standalone-3.141.59.jar selenium-server-standalone
```
And the last version of [Geckodriver](https://github.com/mozilla/geckodriver/releases) for your architecture.

## Installation
## Golang 

To install ``recon-archy`` just run :
```
go get github.com/remiflavien1/recon-archy
```
Next we need to install the dependencies :
- The Selenium server 
- And the Geckodriver. 

Fortunatly the [tebeka/selenium](https://github.com/tebeka/selenium) (which is a internal dependencies of ``recon-archy``) provide everything for us :
```sh
cd $GOPATH/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/
go run init.go --alsologtostderr  --download_browsers --download_latest
```
That's it, you're good to go to [usage](##usage)

### Binaries
You can download the precompiled binaries in the [release](https://github.com/remiflavien1/recon-archy/releases) section.


### From source
Assuming your environnement is well configured (GOPATH, GOROOT...): 
```sh
git clone https://github.com/remiflavien1/recon-archy
cd recon-archy
go build
go install
```

## Usage

First you need to start as many server as worker you want. 
*Notes: The number of workers is limited to 4 for the time being. So at max launch 4 standalone server if you want to work with 4 worker. **This part will be automated in the future***
```sh
# launch 4 standalone servers
./init-server 4
```

Next add your linkedin credential in `.creds` (*interactive mode will be added in the future*)

And then launch ReconArchy (With the previous example, here you will use 4 Workers/Threads)
```sh
recon-archy crawl -t <WORKERS> -c <COMPANY>
```
For example : 
```sh
recon-archy crawl -t 4 -c redhat
```

You can use help menu on command and subcommand for more information. But for now there is not much to cover.
```sh 
$ recon-archy --help
NAME:
   ReconArchy - Crawl 1000 employees of a choosen company and build their organizational chart

USAGE:
   recon-archy [global options] command [command options] [arguments...]

COMMANDS:
   crawl    crawl employees specific to a company
   analyse  Perform analysis on collected data.
   build    Build organisational chart of the company
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

```sh
$ recon-archy crawl --help
NAME:
   recon-archy crawl - crawl employees specific to a company

USAGE:
   recon-archy crawl [command options] [arguments...]

OPTIONS:
   --threads value, -t value  Adjust number of crawling worker (default: "1")
   --company value, -c value  Name of the target company
   --help, -h                 show help (default: false)
```

Crawl result can be retrieve in `/data/`