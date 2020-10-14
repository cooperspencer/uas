package main

import (
	"fmt"
	"gitlab.com/buddyspencer/chameleon"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Found struct {
	pkgname     string
	confinement string
	name        string
	url         string
	platform    string
}

var (
	repoPath = "/tmp/.snap-repos"
	repo     = "https://github.com/cooperspencer/snap-repos"
	found    = []Found{}

	re *regexp.Regexp

	clear = kingpin.Flag("clear", "clear the cache").Short('c').Bool()
	arg   = kingpin.Arg("program", "the program you are looking for").String()
	update = kingpin.Flag("update", "update the images").Short('u').Bool()
)

func main() {
	kingpin.Parse()

	if *update {
		UpdateFlatpaks()
		UpdateSnaps()
		UpdateAppImages()
	} else {
		if *clear {
			err := os.RemoveAll(repoPath)
			if err != nil {
				fmt.Println("couldn't clear cache.\n Please try to remove " + repoPath + " on your own.")
				os.Exit(1)
			}
			fmt.Println("cleared cache")
		}
		if len(*arg) != 0 {
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
					URL:   repo,
					Depth: 0,
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

			fmt.Println("Searching for", chameleon.Lightblue(*arg))

			re = regexp.MustCompile(fmt.Sprintf("(?m)%s", strings.ToUpper(*arg)))

			SearchFlathub()
			SearchSnapcraft()
			SearchAppImage()

			if len(found) == 0 {
				fmt.Println("Nothing found!")
			} else {
				for i, f := range found {
					fmt.Printf("#%d: %s %s\n", i+1, f.platform, chameleon.Lightgreen(f.name))
				}

				var choice_s string
				fmt.Print("Select the number you'd like or 'q' to quit: ")
				_, err := fmt.Scan(&choice_s)
				if err != nil {
					fmt.Println("Invalid choice!")
					os.Exit(1)
				} else {
					if choice_s == "q" {
						fmt.Println("Quitting")
						os.Exit(0)
					}

					choice, err := strconv.Atoi(choice_s)
					if err != nil {
						fmt.Println("Invalid choice!")
						os.Exit(1)
					}
					choice--

					if choice < len(found) && choice >= 0 {
						switch found[choice].platform {
						case "snapcraft":
							InstallSnap(choice)
						case "flathub":
							InstallFlatpak(choice)
						case "appimage":
							InstallAppImage(choice)
						default:
							fmt.Println(found[choice].url)
						}
					} else {
						fmt.Println("Invalid choice!")
						os.Exit(1)
					}
				}
			}
		}
	}
	if len(*arg) == 0 && !*clear && !*update {
		fmt.Println("Please use the --help command")
	}
}

