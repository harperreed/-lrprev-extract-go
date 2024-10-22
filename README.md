# 📥 lrprev-extract-go

## 📝 Summary of Project
`lrprev-extract-go` is a Go-based command-line tool designed for extracting the largest JPEG images embedded within Adobe Lightroom's `.lrprev` files. In addition to extracting images, the tool can also utilize Lightroom's catalog database (`.lrcat`) to ensure that the JPG files are stored in a structured way according to their original paths. 🚀

This project aims to facilitate the management of your Lightroom previews and is especially useful for photographers looking to backup or organize their image assets efficiently. With simple command-line options, users can quickly extract images from directories of Lightroom previews or individual files. 

## ⚙️ How to Use

### Prerequisites
- Go 1.23.2 or later
- Access to a Lightroom catalog (`.lrcat`) if you want to structure your output by original paths.
- For Linux users:
  - OpenGL development libraries
  - X11 development libraries

### Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/harperreed/lrprev-extract-go.git
    cd lrprev-extract-go
    ```

2. For Linux users, install the required development packages:
    ```bash
    # For Ubuntu/Debian:
    sudo apt-get update
    sudo apt-get install libgl1-mesa-dev xorg-dev

    # For Fedora:
    sudo dnf install mesa-libGL-devel libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel

    # For Arch Linux:
    sudo pacman -S mesa libx11 libxcursor libxrandr libxinerama libxi
    ```

3. Compile the code:
    ```bash
    go build -o lrprev-extract ./cmd/lrprev-extract
    ```

### Commands
The main executable is `lrprev-extract`. You can invoke it from the command line with the following options:

```bash
./lrprev-extract [-d <path-to-lightroom-directory> | -f <path-to-lrprev-file>] [-o <output-directory>] [-l <path-to-lrcat>] [-include-size] [-help]
```

- `-d`: Specify the path to a directory containing `.lrdata` files.
- `-f`: Specify the path to an individual `.lrprev` file.
- `-o`: Specify the output directory where the extracted JPEGs should be saved.
- `-l`: Specify the path to your Lightroom catalog (.lrcat) [Optional].
- `-include-size`: Include the size of the images in the filename of the output JPEGs [Optional].
- `-help`: Display help information and usage examples.

If you don't provide the required arguments, the tool will prompt you for the necessary information interactively.

### Example Usage
1. To extract images from a directory:
```bash
./lrprev-extract -d /path/to/lightroom -o /path/to/output
```

2. To extract images from a single `.lrprev` file:
```bash
./lrprev-extract -f /path/to/file.lrprev -o /path/to/output
```

3. To extract images and include size information in the filename:
```bash
./lrprev-extract -d /path/to/lightroom -o /path/to/output -include-size
```

4. To extract images and use the Lightroom catalog for original file paths:
```bash
./lrprev-extract -d /path/to/lightroom -o /path/to/output -l /path/to/catalog.lrcat
```

5. To use the interactive mode:
```bash
./lrprev-extract
```
This will prompt you for the necessary information step by step.

6. To display help information:
```bash
./lrprev-extract -help
```

## 🛠️ Tech Info
- **Language**: Go
- **Dependencies**:
  - `github.com/mattn/go-sqlite3`: A pure Go SQLite driver.
  - `github.com/schollz/progressbar/v3`: A progress bar for console applications.
  - `github.com/go-gl/gl/v2.1/gl`: OpenGL bindings for Go.
  - `github.com/go-gl/glfw/v3.3/glfw`: GLFW bindings for Go.

### Directory Structure
```plaintext
lrprev-extract-go/
├── README.md          # Documentation file
├── cmd                # Command line interface code
│   └── lrprev-extract # Main executable for the tool
│       └── main.go    # Entry point of the application
├── go.mod             # Go module file for dependencies
├── internal           # Internal logic for the application
│   ├── cli            # CLI interaction logic
│   │   └── cli.go
│   ├── database       # Database interaction logic
│   │   └── database.go
│   ├── extractor      # Extraction logic for JPEGs
│   │   └── extractor.go
│   └── utils          # Utility functions
│       └── utils.go
```

### File Descriptions
- **`main.go`**: The main application entry point that handles command-line arguments, interactive prompts, and invokes the appropriate functions for file processing.
- **`cli.go`**: Contains functions for interactive prompts and input validation.
- **`database.go`**: Contains functions for interacting with the Lightroom catalog database to retrieve original file paths.
- **`extractor.go`**: Implements the logic to read `.lrprev` files, extract JPEGs, and manage the output with detailed progress reporting.
- **`utils.go`**: Contains utility functions, including the extraction of UUIDs from filenames.
- **`utils_test.go`**: Contains unit tests for the utility functions in `utils.go`.

### Testing
The project includes unit tests for the utility functions. To run the tests, use the following command:

```bash
go test ./internal/utils
```

Feel free to contribute! We welcome any improvements or bug fixes. 😊

---

For issues and feature requests, please create an issue on the [GitHub Issues Page](https://github.com/harperreed/lrprev-extract-go/issues) or submit a pull request! Happy coding! 💻🔥
