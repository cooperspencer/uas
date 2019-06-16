package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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

func SearchSnapcraft() {
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
			found = append(found, Found{snap.PackageName, snap.Confinement, fmt.Sprintf("%s %s", snap.Name, snap.Version), snap.DownloadURL, "snapcraft"})
		}
	}
	snapcraft = Snapcraft{}
}

func InstallSnap(choice int) {
	if CheckIfCommandExists("snap") {
		confinement := ""
		if len(found[choice].confinement) > 0 {
			confinement =fmt.Sprintf("--%s", found[choice].confinement)
		}
		RunCommand([]string{"snap", "install", found[choice].pkgname, confinement})
	} else {
		fmt.Println("Please install snap")
		os.Exit(1)
	}
}

func UpdateSnaps() {
	if CheckIfCommandExists("snap") {
		fmt.Println("updating snaps...")
		RunCommand([]string{"snap", "refresh"})
	}
}
