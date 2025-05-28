# mdash

A lightweight web server for serving local markdown files using the goldmark
parser.

## Features

- Serves markdown files with GitHub Flavored Markdown support
- Directory browsing with clean UI
- Responsive design
- Security: restricts access to served directory only

## Installation

### From Source

```bash
go install github.com/xdg-go/mdash@latest
```

## Usage

If an `index.md` file exists in a directory, it will be served instead of a
directory listing.

```bash
# Serve current directory on port 3000
mdash

# Serve specific directory on specific port
mdash -dir /path/to/docs -port 8080
```

## License

Apache 2.0 License - see LICENSE file for details.
