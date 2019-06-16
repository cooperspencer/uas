package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Flatpak []struct {
	FlatpakAppID          string      `json:"flatpakAppId"`
	Name                  string      `json:"name"`
	Summary               string      `json:"summary"`
	IconDesktopURL        string      `json:"iconDesktopUrl"`
	IconMobileURL         string      `json:"iconMobileUrl"`
	CurrentReleaseVersion string      `json:"currentReleaseVersion"`
	CurrentReleaseDate    interface{} `json:"currentReleaseDate"`
	InStoreSinceDate      time.Time   `json:"inStoreSinceDate"`
	Rating                float64     `json:"rating"`
	RatingVotes           int         `json:"ratingVotes"`
}

func SearchFlathub() {
	flatpakfile, err := ioutil.ReadFile(fmt.Sprintf("%s/flathub.json", repoPath))
	if err != nil {
		panic(err)
	}
	flathub := Flatpak{}
	err = json.Unmarshal(flatpakfile, &flathub)
	if err != nil {
		panic(err)
	}
	for _, flat := range flathub {
		if len(re.FindAllString(strings.ToUpper(flat.Name), -1)) > 0 {
			found = append(found, Found{flat.FlatpakAppID, "", fmt.Sprintf("%s %s", flat.Name, flat.CurrentReleaseVersion), fmt.Sprintf("https://flathub.org/repo/upstream/%s", flat.FlatpakAppID), "flathub"})
		}
	}
	flathub = Flatpak{}

}

func InstallFlatpak(choice int) {
	if CheckIfCommandExists("flatpak") {
		RunCommand([]string{"flatpak", "install", "-y", "flathub", found[choice].pkgname})
	} else {
		fmt.Println("Please install flatpak")
		os.Exit(1)
	}
}

func UpdateFlatpaks() {
	if CheckIfCommandExists("flatpak") {
		fmt.Println("updating flatpaks...")
		RunCommand([]string{"flatpak", "update"})
	}
}