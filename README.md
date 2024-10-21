# 📥 lrprev-extract-go

## 📝 Summary of Project
`lrprev-extract-go` is a Go-based command-line tool designed for extracting the largest JPEG images embedded within Adobe Lightroom's `.lrprev` files. In addition to extracting images, the tool can also utilize Lightroom's catalog database (`.lrcat`) to ensure that the JPG files are stored in a structured way according to their original paths. 🚀

This project aims to facilitate the management of your Lightroom previews and is especially useful for photographers looking to backup or organize their image assets efficiently. With simple command-line options or a configuration file, users can quickly extract images from directories of Lightroom previews or individual files.

## ⚙️ How to Use

### Prerequisites
- Go 1.23.2 or later
- Access to a Lightroom catalog (`.lrcat`) if you want to structure your output by original paths.

### Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/harperreed/lrprev-extract-go.git
    cd lrprev-extract-go
    ```

2. Compile the code:
    ```bash
    go build -o lrprev-extract ./cmd/lrprev-extract
    ```

### Commands
The main executable is `lrprev-extract`. You can invoke it from the command line with the following options:

```bash
./lrprev-extract [-c <path-to-config-file>] [-d <path-to-lightroom-directory> | -f <path-to-lrprev-file>] [-o <output-directory>] [-l <path-to-lrcat>] [-include-size]
```

- `-c`: Specify the path to a configuration file (YAML format).
- `-d`: Specify the path to a directory containing `.lrdata` files.
- `-f`: Specify the path to an individual `.lrprev` file.
- `-o`: Specify the output directory where the extracted JPEGs should be saved.
- `-l`: Specify the path to your Lightroom catalog (.lrcat) [Optional].
- `-include-size`: Include the size of the images in the filename of the output JPEGs [Optional].

### Configuration File
You can use a YAML configuration file to specify default settings. Here's an example:

```yaml
input_dir: "/path/to/lightroom/directory"
input_file: "/path/to/specific/file.lrprev"
output_directory: "/path/to/output"
lightroom_db: "/path/to/lightroom.lrcat"
include_size: true
```

Command-line arguments will override settings in the configuration file.

### Example Usage
To use a configuration file:
```bash
./lrprev-extract -c /path/to/config.yaml
```

To extract images from a directory:
```bash
./lrprev-extract -d /path/to/lightroom -o /path/to/output
```

To extract images from a single `.lrprev` file:
```bash
./lrprev-extract -f /path/to/file.lrprev -o /path/to/output
```

## 🛠️ Tech Info
- **Language**: Go
- **Dependencies**:
  - `github.com/mattn/go-sqlite3`: A pure Go SQLite driver.
  - `gopkg.in/yaml.v3`: YAML support for Go.

### Directory Structure
```plaintext
lrprev-extract-go/
├── README.md          # Documentation file
├── cmd                # Command line interface code
│   └── lrprev-extract # Main executable for the tool
│       └── main.go    # Entry point of the application
├── go.mod             # Go module file for dependencies
├── internal           # Internal logic for the application
│   ├── config         # Configuration handling
│   │   ├── config.go
│   │   └── config_test.go
│   ├── database       # Database interaction logic
│   │   └── database.go
│   ├── extractor      # Extraction logic for JPEGs
│   │   └── extractor.go
│   └── utils          # Utility functions
│       └── utils.go
```

### File Descriptions
- **`main.go`**: The main application entry point that handles command-line arguments, configuration file parsing, and invokes the appropriate functions for file processing.
- **`config.go`**: Implements the logic for parsing and loading the YAML configuration file.
- **`database.go`**: Contains functions for interacting with the Lightroom catalog database to retrieve original file paths.
- **`extractor.go`**: Implements the logic to read `.lrprev` files, extract JPEGs, and manage the output.
- **`utils.go`**: Contains utility functions, including the extraction of UUIDs from filenames.

Feel free to contribute! We welcome any improvements or bug fixes. 😊

---

For issues and feature requests, please create an issue on the [GitHub Issues Page](https://github.com/harperreed/lrprev-extract-go/issues) or submit a pull request! Happy coding! 💻🔥
