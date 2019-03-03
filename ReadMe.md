# Installation



## Install Go runtime

**MAC**

`brew update && brew install golang` OR ([Download](https://golang.org/doc/install?download=go1.12.darwin-amd64.pkg) and install the package manually)

**Windows**

[Download](https://golang.org/doc/install?download=go1.12.windows-amd64.msi) the MSI package and go through the wizard

### Set GOPATH enviromnet variable

Create a folder on your hard drive where you are going to put your Go projects in. 

**MAC**

Put `export GOPATH=[The Directory That You Have Just Created Above]` in your init file (ie. ~/.bashrc or ~/.zshrc)

**Windows**

Click [here](https://github.com/golang/go/wiki/SettingGOPATH#windows) for the instructions

**Note**

Before proceeding to the next step open a Terminal and validatie your installation:

```shell
go version
echo $GOPATH (echo %GOPATH% on Windows)
```



## Prerequisites

In order to be able to play around with one of the steps in the workshop, you would need to install `mplayer` on your machine.

**MAC**

`brew install mplayer`



## Install Gophobotics

`go get -u -v github.com/xitonix/gophobotics/...`

