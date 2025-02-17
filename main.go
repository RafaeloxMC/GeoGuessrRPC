package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hugolgst/rich-go/client"
)

var startTime *time.Time = func() *time.Time { now := time.Now(); return &now }()
var playing bool = false

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

func setDiscordRPC(mode GameMode, err error, startTime *time.Time) {
	state := "In the menus"
	if mode.Playing {
		state = "Playing"
	}
	err = client.SetActivity(client.Activity{
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

	if err != nil {
		panic(err)
	} else {
		log.Println("Discord RPC updated")
	}
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

	fmt.Printf("Main Path: %s\n", mainPath)
	fmt.Printf("Sub Path: %s\n", subPath)

	switch mainPath {
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
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data URLData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %+v\n", data)

	if data.Action == "close" {
		playing = false
		client.Logout()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Discord RPC cleared")
		return
	}

	mode := extractGameMode(data.URL)
	log.Printf("Received URL: %s\n", data.URL)
	log.Printf("Inferred Game Mode: %s\n", mode.Mode)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "URL received. Game Mode: %s", mode.Mode)

	if !playing {
		playing = true
		startTime = func() *time.Time { now := time.Now(); return &now }()
		err := client.Login("1341072184113762305")
		if err != nil {
			panic(err)
		} else {
			log.Println("Logged in to Discord")
		}
	}

	setDiscordRPC(mode, nil, startTime)
}

func main() {
	http.HandleFunc("/", urlHandler)

	port := "7777"
	log.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
