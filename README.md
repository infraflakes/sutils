# Sane Utils

Sane Utils (`sn`) is an opinionated CLI suite to streamline many command line work.

## Features

*   **Music Conversion:** Convert audio files (FLAC, Opus) to MP3 format with embedded cover art, and format M3U playlists for compatibility.
*   **Archive Operations:** Compress and extract files using 7z, with support for password protection.
*   **YouTube Audio Download:** Download audio from YouTube URLs using `yt-dlp` with embedded thumbnails and metadata.
*   **File and Directory Search:** Search for words within files, or locate files and directories by name, with optional deletion of matched items.
*   **Todo Management:** An interactive terminal-based application for managing todo lists with contexts, priorities, and more.
*   **Smart CD:** Change directory with TUI and aliases.

## Installation

### Quick Try (Run without Installation)

If you want to quickly try it without installing it permanently:

1.  **Ensure Nix is installed** on your system with flake support enabled.
2.  **Run the CLI directly from GitHub:**
    You can run the stable build:
    ```bash
    nix run github:infraflakes/sutils -- [args]
    ```
    (Replace `[args]` with any command and its arguments, e.g., `nix run github:infraflakes/sutils -- music convert mp3 /path/to/dir`)

### Binary Distribution (For Non-Nix Users)

For users not using Nix, the CLI can be downloaded as a single executable binary.

1.  **Download the latest release:**
    Visit the [GitHub Releases page](https://github.com/infraflakes/sutils/releases) and download the wanted binary.

2.  **Make the binary executable:**
    ```bash
    chmod +x sn
    ```

3.  **Move the binary to your PATH (optional but recommended):**
    ```bash
    sudo mv sn /usr/local/bin/
    ```

### Manual Installation (from source)

If you have a Go environment set up, you can build from source.

1.  **Clone the repo:**
    ```bash
    git clone https://github.com/infraflakes/sutils
    cd sutils
    ```

2.  **Build the binary:**
    The included `Makefile` provides an easy way to build the application:
    ```bash
    make build
    ```
    Alternatively, you can use the standard Go command:
    ```bash
    go build -o sn .
    ```

#### Caution!

**In order for `sn cd` to work you need to generate shell functions:**
    Add these to your shell config:

    Bash:
    ```
    eval "$(sn cd init bash)"
    ```

    Zsh:
    ```
    eval "$(sn cd init zsh)"
    ```

    Fish:
    ```
    sn cd init fish | source
    ```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
