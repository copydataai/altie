# 🌈 Alacritty-themes built with Go 🐻 

An alternative to [alacritty-themes](https://github.com/rajasegar/alacritty-themes) but in Go. 

## TODO
- [ ] Implement an append to change just the colors and leave the default config

## Installation
### Online using go install

``` sh
# Please always use the latest version
go install github.com/copydataai/altie/cmd/altie@latest
```

### Cloning the repository

```sh
git clone https://github.com/copydataai/altie

cd altie

# using go install 
# Don't forget to add the GOPATH to PATH
go install ./cmd/altie

# using go build
go build -o altie ./cmd/altie/main.go

# and move to binary dir prefered
sudo mv altie /usr/bin/

altie
```

## License
This project is using the MIT license.
