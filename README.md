# auto-clicker
automatic control of mouse position and click

## Requirements

This app uses the lib github.com/go-vgo/robotgo, so the external requirements can be found there: https://github.com/go-vgo/robotgo#requirements

### For Mac OS X:
```
Xcode Command Line Tools
```

### For Windows:
```
MinGW-w64 (Use recommended) or other GCC
```

### For everything else:
```
GCC, libpng

X11 with the XTest extension (also known as the Xtst library)

Event:

xcb, xkb, libxkbcommon
```

### Ubuntu:
```
sudo apt-get install gcc libc6-dev

sudo apt-get install libx11-dev xorg-dev libxtst-dev libpng++-dev

sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
sudo apt-get install libxkbcommon-dev

sudo apt-get install xsel xclip
```

### Fedora:
```
sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel

sudo dnf install libpng-devel

sudo dnf install xsel xclip
```

## Installation

```
go get github.com/kpacha/auto-clicker
go install github.com/kpacha/auto-clicker
```

## Usage

```
auto-clicker -h
Usage of auto-clicker:
  -f string
    	path of the file containning the animation (default "frames.json")
  -s duration
    	time to sleep between clicks (default 45s)
```
