package main

import (
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const filename = "projects.json"

type History []Session

type Session struct {
	Begin    time.Time
	End      time.Time
	Duration time.Duration
	Commits  int
}

type Project struct {
	Name        string
	Description string
	Created     time.Time
	Duration    time.Duration
	History     History
	Commits     int
	Id          int
}

type ProjectList struct {
	List  []Project
	maxId int
}

func (p *ProjectList) save() error {
	data, err := json.MarshalIndent(p, "", "	")
	if err != nil {
		log.Fatalf("error during json marshaling : %s", err)
	}
	return ioutil.WriteFile(filename, []byte(data), 0600)
}

func loadProjects() ProjectList {
	var projects ProjectList
	// if the file doesn't exist, the project list is empty
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return projects
	} else {
		// process data
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("error reading projects list : #{err}")
		}
		if err := json.Unmarshal(data, &projects); err != nil {
			log.Fatalf("error during json unmarshaling: %s", err)
		}
		return projects
	}
}

func (p *Project) Add(s Session) {
	p.Duration += s.Duration
	p.Commits += s.Commits
	p.History = append(p.History, s)
}

func initCreateGUI() {
	form := ui.NewForm()
	form.SetPadded(true)
	titleEntry := ui.NewEntry()
	form.Append("Enter project name", titleEntry, false)
	descriptionEntry := ui.NewMultilineEntry()
	form.Append("Enter the description of the project", descriptionEntry, true)
	button := ui.NewButton("\n\n\n\n Save this project \n\n\n\n")
	form.Append("", button, false)

	projects := loadProjects()
	maxId := projects.maxId + 1
	window := ui.NewWindow("Create a project", 1200, 600, true)
	window.SetChild(form)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	button.OnClicked(func(b *ui.Button) {
		var history History
		var duration time.Duration
		title := titleEntry.Text()
		//description := descriptionEntry.Text()
		description := "a"
		project := Project{title, description, time.Now(), duration, history, 0, maxId}
		projects.List = append(projects.List, project)
		projects.maxId = maxId
		projects.save()
		window.Destroy()
		ui.Quit()
	})
	window.Show()

}

func main() {
	// Send and receive times for tracking
	beginTimes := make(chan time.Time)
	endTimes := make(chan time.Time)
	sessions := make(chan Session)
	//Send and receive projects ids
	selectedProject := make(chan int)
	go func() {
		project := <-selectedProject
		fmt.Println(project)
	}()

	go func() {
		for {
			session := <-sessions
			fmt.Println(session)
		}
	}()

	go func() {
		for {
			beginTime := <-beginTimes
			endTime := <-endTimes
			duration := endTime.Sub(beginTime)
			session := Session{beginTime, endTime, duration, 0}
			sessions <- session
		}
	}()
	ui.Main(initCreateGUI)
	//Play/Pause button
	err := ui.Main(func() {
		box := ui.NewVerticalBox()
		button := ui.NewButton("Play")
		box.Append(button, true)
		window := ui.NewWindow("Hello", 400, 200, false)
		window.SetMargined(true)
		window.SetChild(box)
		button.OnClicked(func(b *ui.Button) {
			if b.Text() == "Play" {
				b.SetText("Pause")
				beginTimes <- time.Now()
			} else {
				b.SetText("Play")
				endTimes <- time.Now()
			}
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}

}
