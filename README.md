# gorename

gorename is a tool for moving packages to the new import paths.

## Installation

The go get command will automatically fetch the dependencies listed above, compile the binary and place it in your `$GOPATH/bin` directory.  
```shell script
go get -u github.com/Orogenesis/gorename/cmd/gorename
```

## Usage

```text
Usage:
  gorename [flags]

Flags:
  -path             Package path you want to replace
  -new-path         New path of the package
  -root-dir         Root directory path of the project (default ".")
  -use-modules      Replace module name in go.mod
  -print-result     Display progress
```

## Example

```shell script
gorename -path gitlab.com/orogenesis/gorename -new-path github.com/orogenesis/gorename
```

## License

The MIT License (MIT).
