package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/RafaeloxMC/richer-go/client"
)

var startTime *time.Time = func() *time.Time { now := time.Now(); return &now }()
var playing bool = false
var loggedin bool = false
var debug *bool = nil

type URLData struct {
	URL    string `json:"url"`
	Action string `json:"action"`
}

type GameMode struct {
	Mode           string `json:"mode"`
	SmallImage     string `json:"smallImage"`
	SmallImageText string `json:"smallImageText"`
	Playing        bool   `json:"playing"`
}

func logDebug(log string) {
	if *debug {
		fmt.Println("[DEBUG] " + log)
	}
}

var currentMode GameMode

func setDiscordRPC(mode GameMode, startTime *time.Time) {
	if currentMode != (GameMode{}) && mode == currentMode {
		log.Println("Discord RPC already up to date")
		return
	}
	state := "In the menus"
	if mode.Playing {
		state = "Playing"
	}
	err := client.SetActivity(client.Activity{
		State:      state,
		Details:    mode.Mode,
		LargeImage: "large",
		LargeText:  "GeoGuessr by @xvcf",
		SmallImage: mode.SmallImage,
		SmallText:  mode.SmallImageText,
		Timestamps: &client.Timestamps{
			Start: startTime,
		},
	})

	if err != nil || !loggedin {
		log.Printf("Error updating Discord RPC: %v", err.Error())
		log.Println("Attempting to re-login")
		logout()
		time.Sleep(500 * time.Millisecond)
		login_err := login()
		if login_err != nil {
			log.Printf("Failed to re-login: %v", login_err.Error())
			return
		}
		time.Sleep(1 * time.Second)
		setDiscordRPC(mode, startTime)
	} else {
		log.Println("Discord RPC updated")
	}
	currentMode = mode
	logDebug(fmt.Sprintf("Current mode: %s", mode.Mode))
}

func extractGameMode(url string) GameMode {
	re := regexp.MustCompile(`https://www\.geoguessr\.com(?:/[a-z]{2})?/([a-z-]+)(?:/([a-zA-Z0-9-]+))?`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return GameMode{
			Mode:           "Not in a game",
			SmallImage:     "",
			SmallImageText: "",
			Playing:        false,
		}
	}

	mainPath := matches[1]
	subPath := ""
	if len(matches) > 2 {
		subPath = matches[2]
	}

	logDebug(fmt.Sprintf("Main path: %s, Sub path: %s", mainPath, subPath))

	switch mainPath {
	case "singleplayer":
		if subPath == "game" {
			return GameMode{
				Mode:           "Singleplayer Game",
				SmallImage:     "singleplayer",
				SmallImageText: "Singleplayer Game",
				Playing:        true,
			}
		}
	case "game":
		{
			return GameMode{
				Mode:           "Singleplayer Game",
				SmallImage:     "singleplayer",
				SmallImageText: "Singleplayer Game",
				Playing:        true,
			}
		}
	case "maps":
		{
			return GameMode{
				Mode:           "Selecting a map",
				SmallImage:     "",
				SmallImageText: "",
				Playing:        false,
			}
		}
	case "duels":
		if subPath != "" {
			return GameMode{
				Mode:           "Duels Game",
				SmallImage:     "ranked",
				SmallImageText: "Ranked Duels Game",
				Playing:        true,
			}
		}
		return GameMode{
			Mode:           "Duels",
			SmallImage:     "ranked",
			SmallImageText: "Ranked Duels",
			Playing:        true,
		}
	case "team-duels":
		if subPath != "" {
			return GameMode{
				Mode:           "Team Duels Game",
				SmallImage:     "rankedteams",
				SmallImageText: "Ranked Team Duels Game",
				Playing:        true,
			}
		}
		return GameMode{
			Mode:           "Team Duels",
			SmallImage:     "rankedteams",
			SmallImageText: "Team Duels (Ranked)",
			Playing:        true,
		}
	case "teams":
		if subPath != "" {
			return GameMode{
				Mode:           "Team Duels Game",
				SmallImage:     "teamduels",
				SmallImageText: "Team Duels Game",
				Playing:        true,
			}
		}
		return GameMode{
			Mode:           "Team Duels",
			SmallImage:     "teamduels",
			SmallImageText: "Team Duels",
			Playing:        true,
		}
	case "multiplayer":
		switch subPath {
		case "teams":
			return GameMode{
				Mode:           "Team Duels",
				SmallImage:     "rankedteams",
				SmallImageText: "Team Duels (Ranked)",
				Playing:        false,
			}
		case "battle-royale-distance":
			return GameMode{
				Mode:           "Battle Royale Distance",
				SmallImage:     "br-distance",
				SmallImageText: "Battle Royale Distance",
				Playing:        false,
			}
		case "battle-royale-countries":
			return GameMode{
				Mode:           "Battle Royale Countries",
				SmallImage:     "br-countries",
				SmallImageText: "Battle Royale Countries",
				Playing:        false,
			}
		case "unranked-teams":
			return GameMode{
				Mode:           "Unranked Team Duels",
				SmallImage:     "teamduels",
				SmallImageText: "Unranked Team Duels",
				Playing:        false,
			}
		default:
			return GameMode{
				Mode:           "Duels",
				SmallImage:     "ranked",
				SmallImageText: "Ranked Duels",
				Playing:        false,
			}
		}
	case "battle-royale":
		if subPath != "" {
			return GameMode{
				Mode:           "Battle Royale Game",
				SmallImage:     "br-distance",
				SmallImageText: "Battle Royale Distance Game",
				Playing:        true,
			}
		}
		return GameMode{
			Mode:           "Battle Royale Distance",
			SmallImage:     "br-distance",
			SmallImageText: "Battle Royale Distance",
			Playing:        true,
		}
	default:
		return GameMode{
			Mode:           "Not in a game",
			SmallImage:     "",
			SmallImageText: "",
			Playing:        false,
		}
	}
	return GameMode{
		Mode:           "Not in a game",
		SmallImage:     "",
		SmallImageText: "",
		Playing:        false,
	}
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logDebug("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data URLData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logDebug("Failed to parse request body")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if data.Action == "close" {
		logDebug("Game closed, clearing Discord RPC")
		playing = false
		client.Logout()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Discord RPC cleared")
		return
	}

	mode := extractGameMode(data.URL)

	w.WriteHeader(http.StatusOK)

	if !loggedin {
		login_err := login()
		if login_err != nil {
			log.Printf("Failed to login: %v", login_err.Error())
			return
		}
	}

	if !playing {
		playing = true
		startTime = func() *time.Time { now := time.Now(); return &now }()
		logDebug("Game started, now playing")
	}

	setDiscordRPC(mode, startTime)
}

func login() error {
	logDebug("Sending login request to Discord")
	err := client.Login("1341072184113762305")
	if err != nil {
		loggedin = false
		return err
	} else {
		log.Println("Logged in to Discord")
		loggedin = true
	}
	return nil
}

func logout() {
	if loggedin {
		client.Logout()
		log.Println("Logged out of Discord")
		loggedin = false
	} else {
		log.Println("Already logged out of Discord")
	}
}

func main() {
	debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		logDebug("--------------------- Debug mode enabled ---------------------")
	}

	http.HandleFunc("/", urlHandler)

	port := "7777"
	logDebug(fmt.Sprintf("Server is running on port %s", port))
	log.Println("Starting...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
