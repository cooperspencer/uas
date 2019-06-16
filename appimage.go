package main

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"gitlab.com/buddyspencer/chameleon"
	"gopkg.in/cheggaaa/pb.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

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

type AppImageSave struct {
	Appimage [] AppImageSlot `yaml:"appimage"`
}

type AppImageSlot struct {
	Program string `yaml:"program"`
	File    string `yaml:"file"`
}

func SearchAppImage() {
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
				if link.Type == "uas" {
					found = append(found, Found{app.Name, "", app.Name, link.URL, "appimage"})
				}
			}
		}
	}
	appimage = Appimage{}
}

func InstallAppImage(choice int) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	var appimagesave *AppImageSave
	appimage_dir := fmt.Sprintf("%s/.appimage", home)
	appimage_save := fmt.Sprintf("%s/.appimage.yml", appimage_dir)

	fmt.Println("saving to", appimage_dir)
	if _, err = os.Stat(appimage_dir); os.IsNotExist(err) {
		err = os.MkdirAll(appimage_dir, 0700)
		fmt.Printf("%s not found. Trying to create it.\n", appimage_dir)
		if err != nil {
			fmt.Println("could not create " + appimage_dir + ".")
			os.Exit(1)
		}

		fmt.Println("Created it.")
	}

	if _, err = os.Stat(appimage_save); os.IsNotExist(err) {
		appimagesave = &AppImageSave{}
	} else {
		appimagesave, _ = ReadAppImageSaveFile(appimage_save)
	}

	var bar *pb.ProgressBar
	var started bool
	progress := func(current, total int64) {
		if !started {
			bar = pb.New(int(total)).SetUnits(pb.U_BYTES)
			bar.Start()
			started = true
		}
		bar.Set(int(current))
	}

	r, err := req.Get(found[choice].url, req.DownloadProgress(progress))
	if err != nil {
		fmt.Println("could not download", found[choice].url)
		os.Exit(1)
	}
	splittedLink := strings.Split(found[choice].url, "/")
	filename := fmt.Sprintf("%s/%s", appimage_dir, splittedLink[len(splittedLink) - 1])
	err = r.ToFile(filename)
	if err != nil {
		fmt.Println("could not save file.")
		os.Exit(1)
	}
	fmt.Printf("Downloaded %s\n", filename)
	err = os.Chmod(filename, 0775)
	if err != nil {
		fmt.Println("could not set permissions on ", filename)
		os.Exit(1)
	}
	fmt.Printf("Set execute permission on %s\n", filename)

	f := false

	for x, slot := range appimagesave.Appimage {
		if found[choice].name == slot.Program {
			appimagesave.Appimage[x].File = splittedLink[len(splittedLink) - 1]
			f = true
			break
		}
	}

	if !f {
		appimagesave.Appimage = append(appimagesave.Appimage, AppImageSlot{found[choice].name, splittedLink[len(splittedLink) - 1]})
	}

	WriteAppImageSaveFile(appimage_save, appimagesave)
}

func ReadAppImageSaveFile(configfile string) (*AppImageSave, bool) {
	cfgdata, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Panic("Cannot open config file from " + configfile)
	}

	t := AppImageSave{}

	err = yaml.Unmarshal([]byte(cfgdata), &t)

	if err != nil {
		log.Panic("Cannot map yml config file to interface, possible syntax error")
		log.Panic(err)
	}

	return &t, true
}

func WriteAppImageSaveFile(configfile string, config *AppImageSave) {
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile(configfile, d, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func UpdateAppImages() {
	fmt.Println("updating appimages...")
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	appimage_dir := fmt.Sprintf("%s/.appimage", home)
	appimage_save := fmt.Sprintf("%s/.appimage.yml", appimage_dir)

	toDelete := []string{}
	if _, err = os.Stat(appimage_save); !os.IsNotExist(err) {
		appimagesave, _ := ReadAppImageSaveFile(appimage_save)
		files, _ := ioutil.ReadDir(appimage_dir)
		for _, f := range files {
			for _, a := range appimagesave.Appimage {
				if f.Name() == a.File {
					re = regexp.MustCompile(fmt.Sprintf("(?m)%s", strings.ToUpper(a.Program)))
					SearchAppImage()
					dfile := strings.Split(found[len(found)-1].url, "/")
					if dfile[len(dfile)-1] == f.Name() {
						found = found[:len(found)-1]
					} else {
						toDelete = append(toDelete, f.Name())
					}
				}
			}
		}

		for k := range found {
			fmt.Println("updating", chameleon.BLightblue(found[k].name))
			InstallAppImage(k)
		}
		for _, del := range toDelete {
			err := os.Remove(fmt.Sprintf("%s/%s", appimage_dir, del))
			if err != nil {
				fmt.Println("couldn't delete", del)
			}
		}
	}
}