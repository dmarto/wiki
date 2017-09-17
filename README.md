# wiki

wiki is a self-hosted wiki engine/content management system that lets
you create and share content in Markdown format (for now).

### Source

```#!bash
$ go get github.com/dmarto/wiki
```

## Usage

Run wiki:

```#!bash
$ wiki
```

Visit: http://localhost:8000/

Start creating/editing content!

## Configuration

By default wiki pages are stored in `./data` in the local directory. This can
be changed by supplying the `-data /path/to/data` option.

## License

MIT
