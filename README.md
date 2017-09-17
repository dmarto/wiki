# wiki

wiki is a self-hosted wiki engine for static markdown files

### Source

```#!bash
$ go get github.com/dmarto/wiki
```

## Usage

Run wiki:

```#!bash
$ wiki
```

## Configuration

By default wiki pages are stored in `./data` in the local directory.
This can be changed by supplying the `-data /path/to/data` option.

By default wiki binds to `0.0.0.0:8000`.
This can be changed by supplying the `-bind [addr]:<port>` option.

## License

MIT
