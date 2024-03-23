<div align="center">
    <h1 align="center">Altie</h1>
    <h3>A simple way to change themes and font for alacritty 💟</h3>
</div>
<div align="center">
  <a href="https://github.com/copydataai/altie/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/license-MIT-purple"></a>
</div>
<br />


## Todo
- [ ] Refactor to a TUI simpler than pterm
- [ ] implement tags to change font
- [ ] implement a select font

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
