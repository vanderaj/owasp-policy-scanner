package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoConnUrl string = "mongodb://localhost:27017/chapters"

// var {
//    ctx
//    client *mongo.Client
//    chaptersCollection
// }

// Store results
func storeResult(chapter string, key string, result string) {

}

// Out of date dependencies in _config.yml
func checkConfigYml() {

}

// Default text in index.md
func checkDefaultText() {

}

// Default tab tab_example.md is present
func checkDefaultExampleTab() {

}

// Automigration metadata is present and set to 1
func checkDefaultMigrationHeader() {

}

// Number of leaders < 2 or > 5
func checkLeaderCount() {
	// Is this file leaders.md? If not, return

	// Count the number of lines starting with * that also contain an email address

	// Are there less than two?

	// Are there more than five?
}

// Leaders in leaders.md doesn’t match Copper
func checkLeadersInCopper() {
}

// Out of date .gitignore
func checkOldGitIgnore(s string, d fs.DirEntry) {
	//
}

// _site/ being present
func checkIfSite(s string, d fs.DirEntry) {
	if d.IsDir() {
		if d.Name() == "_site" {
			println(s)
		}
	}

}

// Meetup header present but no active Meetup for that chapter
func checkMeetupExists() {

}

// Meetup header present and less than three meetings since Feb 2020
func checkPastMeetups() {

}

// Meetup header present and Link to Meetup in info.md but no metadata JavaScript for automated (warning)
func checkMeetupMissingMetaData() {

}

// Old Wiki links are present (a warning not a breakage)
func checkForOldWiki() {

}

// Old policy links are present (a warning not a breakage)
func checkForOldPolicy(filename string, d fs.DirEntry) error {

	if !strings.Contains(filename, ".md") {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "index.php") {
			fmt.Println(filename, "Old wiki found on line ", line)
			return nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		println("checkForOldPolicy error:", err.Error())
	}

	return nil
}

// Tab filename and title metadata is incorrect
func checkTabTags() {

}

// Jekyll bundle fails to build
func checkJekyllBuilds() {

}

// Update git repos
func updateGit(s string, d fs.DirEntry) {

	if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
		println("Updating ", d.Name())
		cmd := exec.Command("git", "pull")
		cmd.Dir = s
		output, err := cmd.Output()
		cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", output)
	}
}

func walk(s string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}

	// Directory Checks
	if d.IsDir() {
		updateGit(s, d)
		checkIfSite(s, d)
	}

	// File Checks
	if !d.IsDir() {
		checkForOldPolicy(s, d)
	}
	return nil
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	chaptersDatabase := client.Database("chapters")
	chaptersCollection := chaptersDatabase.Collection("chapters")

	if chaptersCollection.Name() == "chapters" {
		println("Connected to MongoDB")
	}

	filepath.WalkDir("../chapters/", walk)
}
