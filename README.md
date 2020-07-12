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

