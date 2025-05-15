# RSS Reader

A simple RSS feed reader application written in Go using the Fyne framework for creating a graphical user interface.

## Features

- Add, delete, and update RSS feeds
- Display articles from feeds
- View article content
- Open links in browser
- Mark articles as read
- Automatically save added feeds

## Requirements

- Go 1.16 or higher
- For compiling the GUI, Fyne dependencies are required:
  - On Linux: `gcc` and GTK3 dependencies
  - On macOS: Xcode Command Line Tools
  - On Windows: GCC (can be installed via MinGW or MSYS2)

## Installation

```bash
# Clone the repository
git clone https://github.com/romanitalian/rss-reader.git
cd rss-reader

# Build the project (using Makefile)
make build

# Run
make run

# Or manually
go build -o rssreader ./cmd/rssreader
./rssreader
```

## Using Makefile

The project includes a Makefile to simplify development tasks:

```bash
# Build the project
make build

# Run the application
make run

# Clean the build directory
make clean

# Build for all platforms (Linux, Windows, macOS)
make build-all

# Run tests
make test

# Update dependencies
make tidy

# Format code
make fmt

# Check code with linter
make lint

# Show help for commands
make help
```

## Usage

1. At first launch, the application will add several default RSS feeds
2. To add a new feed, click the "+" button at the bottom of the feed list
3. To update feeds, click the refresh icon button
4. To delete a feed, select it and click the trash icon button
5. To view articles, select a feed from the list on the left
6. To read an article, select it from the article list
7. To open an article in a browser, click the "Open in Browser" button

## Configuration

The configuration file is located at `~/.rssreader/config.json` and contains the following settings:

- `refresh_interval`: feed update interval in minutes
- `auto_refresh`: automatic feed updates
- `default_feeds`: default RSS feed list
- `data_dir`: data storage directory

## License

MIT
