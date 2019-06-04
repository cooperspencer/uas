package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
	"gitlab.com/buddyspencer/chameleon"
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

type Snapcraft struct {
	Embedded struct {
		ClickindexPackage []struct {
			Aliases             interface{}   `json:"aliases"`
			AnonDownloadURL     string        `json:"anon_download_url"`
			Apps                []string      `json:"apps"`
			Architecture        []string      `json:"architecture"`
			BinaryFilesize      int           `json:"binary_filesize"`
			Channel             string        `json:"channel"`
			CommonIds           []interface{} `json:"common_ids"`
			Confinement         string        `json:"confinement"`
			Contact             string        `json:"contact"`
			Content             string        `json:"content"`
			DatePublished       time.Time     `json:"date_published"`
			Deltas              []interface{} `json:"deltas"`
			Description         string        `json:"description"`
			DeveloperID         string        `json:"developer_id"`
			DeveloperName       string        `json:"developer_name"`
			DeveloperValidation string        `json:"developer_validation"`
			DownloadSha3384     string        `json:"download_sha3_384"`
			DownloadSha512      string        `json:"download_sha512"`
			DownloadURL         string        `json:"download_url"`
			Epoch               string        `json:"epoch"`
			GatedSnapIds        []interface{} `json:"gated_snap_ids"`
			IconURL             string        `json:"icon_url"`
			LastUpdated         time.Time     `json:"last_updated"`
			License             string        `json:"license"`
			Name                string        `json:"name"`
			Origin              string        `json:"origin"`
			PackageName         string        `json:"package_name"`
			Prices              struct {
			} `json:"prices"`
			Private        bool        `json:"private"`
			Publisher      string      `json:"publisher"`
			RatingsAverage float64     `json:"ratings_average"`
			Release        []string    `json:"release"`
			Revision       int         `json:"revision"`
			ScreenshotUrls []string    `json:"screenshot_urls"`
			SnapID         string      `json:"snap_id"`
			Summary        string      `json:"summary"`
			SupportURL     string      `json:"support_url"`
			Title          string      `json:"title"`
			Version        string      `json:"version"`
			Website        interface{} `json:"website"`
			Base           string      `json:"base,omitempty"`
		} `json:"clickindex:package"`
	} `json:"_embedded"`
}

type Appimage struct {
	Version     int    `json:"version"`
	HomePageURL string `json:"home_page_url"`
	FeedURL     string `json:"feed_url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Favicon     string `json:"favicon"`
	Expired     bool   `json:"expired"`
	Items       []struct {
		Name        string   `json:"name"`
		Description string   `json:"description,omitempty"`
		Categories  []string `json:"categories"`
		Authors     []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"authors"`
		License interface{} `json:"license"`
		Links   []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"links"`
		Icons       []string `json:"icons"`
		Screenshots []string `json:"screenshots"`
	} `json:"items"`
}

type Found struct {
	name     string
	url      string
	platform string
}

var (
	repoPath = "/tmp/.snap-repos"
	repo     = "https://gitlab.com/buddyspencer/snap-repos"
	found    = []Found{}
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL: repo,
			})
			if err != nil {
				panic(err)
			}
		} else {
			r, err := git.PlainOpen(repoPath)
			if err != nil {
				panic(err)
			}
			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				if err.Error() != "already up-to-date" {
					panic(err)
				}
			}
		}

		fmt.Println("Searching for", chameleon.Lightblue(args[0]))

		var re= regexp.MustCompile(fmt.Sprintf("(?m)%s", strings.ToUpper(args[0])))

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
				found = append(found, Found{fmt.Sprintf("%s %s", flat.Name, flat.CurrentReleaseVersion), fmt.Sprintf("https://flathub.org/repo/upstream/%s", flat.FlatpakAppID), "flathub"})
			}
		}
		flathub = Flatpak{}

		snapcraftfile, err := ioutil.ReadFile(fmt.Sprintf("%s/snapcraft.json", repoPath))
		if err != nil {
			panic(err)
		}
		snapcraft := Snapcraft{}
		err = json.Unmarshal(snapcraftfile, &snapcraft)
		if err != nil {
			panic(err)
		}
		for _, snap := range snapcraft.Embedded.ClickindexPackage {
			if len(re.FindAllString(strings.ToUpper(snap.Name), -1)) > 0 {
				found = append(found, Found{fmt.Sprintf("%s %s", snap.Name, snap.Version), snap.DownloadURL, "snapcraft"})
			}
		}
		snapcraft = Snapcraft{}

		appimagefile, err := ioutil.ReadFile(fmt.Sprintf("%s/appimage.json", repoPath))
		if err != nil {
			panic(err)
		}
		appimage := Appimage{}
		err = json.Unmarshal(appimagefile, &appimage)
		if err != nil {
			panic(err)
		}
		for _, app := range appimage.Items {
			if len(re.FindAllString(strings.ToUpper(app.Name), -1)) > 0 {
				for _, link := range app.Links {
					if link.Type == "Download" {
						found = append(found, Found{fmt.Sprintf("%s", app.Name), link.URL, "appimage"})
					}
				}
			}
		}
		appimage = Appimage{}

		if len(found) == 0 {
			fmt.Println("Nothing found!")
		} else {
			for i, f := range found {
				fmt.Printf("#%d: %s %s\n", i, f.platform, chameleon.Lightgreen(f.name))
			}

			var choice int
			fmt.Print("Select the number you'd like: ")
			_, err = fmt.Scan(&choice)
			if err != nil {
				fmt.Println("Invalid choice!")
				os.Exit(1)
			} else {
				if choice < len(found) && choice >= 0 {
					fmt.Println(found[choice].url)
				} else {
					fmt.Println("Invalid choice!")
					os.Exit(1)
				}
			}
		}
	}
}
