# Filemux

`filemux` is a command-line tool written in Go that compiles the contents of multiple files into a single, formatted output suitable for pasting into LLM (Large Language Model) chats. It supports glob patterns, directories, and individual files, with options for clipboard output and handling large or binary files.

## Features

- **File Compilation**: Combines contents of multiple files into a Markdown-formatted output with headers and code blocks.
- **Input Flexibility**: Accepts glob patterns (e.g., `*.txt`), directory paths (recursively processes all files), and specific file names.
- **Size Limit**: Enforces a 500KB max file size, with a `-f` flag to force processing of larger files.
- **Binary Detection**: Warns if a file appears to be binary (may not display well in text-based chats).
- **Clipboard Support**: Copies output to the clipboard using `-c` or `-clipboard` flags (works on macOS, Windows, and WSL with proper setup).
- **Error Handling**: Exits with an error if any specified file or directory doesn't exist or if a glob pattern matches nothing.

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) (version 1.16 or later) installed to build from source.
- For clipboard support:
  - **macOS**: No additional setup needed.
  - **Windows**: No additional setup needed.
  - **WSL**: Requires a clipboard integration tool like `xclip` or Windows clipboard access configured.

### Pre-built Binary
Download the latest pre-built binary for your platform from the [Releases](https://github.com/ragulmathawa/filemux/releases) page:
1. Choose the appropriate binary (e.g., `filemux-linux-amd64`, `filemux-darwin-amd64`, `filemux-windows-amd64.exe`).
2. Download and make it executable (Linux/macOS):
   ```bash
   chmod +x filemux-<platform>
   ```
3. Move to a directory in your PATH (optional, user-specific):
   ```bash
   mkdir -p ~/.local/bin && mv filemux-<platform> ~/.local/bin/filemux
   ```

Alternatively, use these single-line commands to install directly into `~/.local/bin` (no `sudo` required):

#### For macOS (Intel-based):
```bash
curl -sL https://api.github.com/repos/ragulmathawa/filemux/releases/latest | grep "browser_download_url.*filemux-darwin-amd64" | cut -d'"' -f4 | xargs curl -L -o ~/.local/bin/filemux && chmod +x ~/.local/bin/filemux
```

#### For WSL/Linux (x86_64):
```bash
curl -sL https://api.github.com/repos/ragulmathawa/filemux/releases/latest | grep "browser_download_url.*filemux-linux-amd64" | cut -d'"' -f4 | xargs curl -L -o ~/.local/bin/filemux && chmod +x ~/.local/bin/filemux
```

**Note**: Ensure `~/.local/bin` is in your PATH. You can check with `echo $PATH | grep ~/.local/bin`. If it’s not, add it to your shell configuration (e.g., `~/.bashrc` or `~/.zshrc`):
```bash
export PATH="$HOME/.local/bin:$PATH"
```
Then, reload your shell: `source ~/.bashrc` or `source ~/.zshrc`.


### Build from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/ragulmathawa/filemux.git
   cd filemux
   ```
2. Install the required dependency:
   ```bash
   go get github.com/atotto/clipboard
   ```
3. Build the tool:
   ```bash
   go build filemux.go
   ```
4. (Optional) Move the binary to a directory in your PATH:
   ```bash
   sudo mv filemux /usr/local/bin/
   ```

## Usage

```
Usage: filemux [flags] <file-patterns-or-dirs...>
```

### Flags
- `-c`: Copy output to clipboard (alternative to `-clipboard`).
- `-clipboard`: Copy output to clipboard instead of printing.
- `-f`: Force processing of files larger than 500KB.

### Examples
1. Compile all `.txt` files in the current directory:
   ```bash
   filemux *.txt
   ```
2. Process all files in a directory and copy to clipboard:
   ```bash
   filemux -c ./docs
   ```
3. Force processing of a large file:
   ```bash
   filemux -f largefile.txt
   ```
4. Mix of patterns, files, and directories:
   ```bash
   filemux *.go main.go ./src -f
   ```
5. Attempt to process a non-existent file (will error):
   ```bash
   filemux nonexistent.txt
   # Output: Error: File or directory 'nonexistent.txt' not found
   ```

### Output Format
The output is formatted as Markdown for easy pasting into LLM chats:
```
Here's the compiled content from multiple files:

### File: path/to/file1.txt
```
Content of file1
```

---

### File: path/to/file2.txt
```
Content of file2
```
```

## Notes
- **File Size Limit**: Files larger than 500KB will trigger an error unless `-f` is used.
- **Binary Files**: A warning is printed for files with non-printable characters, but they are still included.
- **Clipboard in WSL**: May require additional setup (e.g., `xclip` or Windows integration) to work properly.

## Contributing
Feel free to submit issues or pull requests to the [GitHub repository](https://github.com/ragulmathawa/filemux). Contributions are welcome!

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Author
- Ragul Mathawa ([ragulmathawa](https://github.com/ragulmathawa))
