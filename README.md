# 🌍 GeoGuessr Discord RPC

GeoGuessr Discord RPC is a project that integrates GeoGuessr with Discord Rich Presence, allowing you to display your current game status on your Discord profile.

## 📜 Features

-   **Automatic Detection**: Automatically detects when you start or stop playing GeoGuessr.
-   **Game Mode Display**: Shows the current game mode you are playing (e.g., Duels, Team Duels, Battle Royale).
-   **Rich Presence**: Updates your Discord status with detailed information about your game.

## 🚀 Getting Started

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.24 or later)
-   [Chrome Browser](https://www.google.com/chrome/)

### Installation

1. **Clone the repository**:

    ```sh
    git clone https://github.com/RafaeloxMC/GeoGuessrRPC.git
    cd GeoGuessrRPC
    ```

2. **Install Go dependencies**:

    ```sh
    go mod tidy
    ```

3. **Build the Go server**:

    ```sh
    go build -o geoguessrrpc main.go
    ```

4. **Set up the Chrome extension**:
    - Open Chrome and go to `chrome://extensions/`.
    - Enable "Developer mode" (toggle in the top right).
    - Click "Load unpacked" and select the `extension` directory from this repository.

### Usage

1. **Run the Go server**:

    ```sh
    ./geoguessrrpc
    ```

2. **Start playing GeoGuessr**:
    - Open GeoGuessr in your browser.
    - Your Discord status will automatically update with your current game mode.

## 🛠️ How It Works

### Server (`main.go`)

-   The Go server listens for HTTP POST requests from the Chrome extension.
-   It parses the URL and determines the current game mode.
-   It updates the Discord Rich Presence using the [`rich-go`](https://github.com/hugolgst/rich-go/) library.

### Chrome Extension (`extension/background.js`)

-   Listens for tab updates and removals.
-   Sends the current GeoGuessr URL to the Go server when a GeoGuessr tab is opened or closed.

## 🤝 Contributing

Contributions are welcome! Please fork this repository and submit pull requests. If you have any problems, feel free to open an issue!

## 📄 License

This project is licensed under the GNU General Public License - see the [LICENSE](LICENSE) file for details.

---

Enjoy playing GeoGuessr with enhanced Discord Rich Presence! 🌍🎮
