package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// const      mongoConnUrl string = "mongodb://localhost:27017/chapters"

var currChapter string

//    ctx
//    client *mongo.Client
//    chaptersCollection

// Store results
func storeResult(chapter string, key string, result string) {

}

// Out of date dependencies in _config.yml
func checkConfigYml() {

}

// Default text in index.md
func checkDefaultText(filename string, d fs.DirEntry) error {
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
		if strings.Contains(scanner.Text(), "Standard Chapter Page Template") {
			fmt.Printf("Default text present in %s on line %d\n", filename, line)
			return nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkDefaultText error: ", err.Error())
	}

	return nil
}

// Default tab tab_example.md is present
func checkDefaultExampleTab(filename string, d fs.DirEntry) error {
	if strings.Contains(filename, "tab_example.md") {
		fmt.Println("Example tab found at: ", filename)
	}

	return nil
}

// Automigration metadata is present and set to 1
func checkDefaultMigrationHeader(filename string, d fs.DirEntry) error {
	if !strings.Contains(filename, "index.md") {
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
		if strings.Contains(scanner.Text(), "auto-migrated: 1") {
			fmt.Println("Auto-Migration Headers active in ", filename, "on line", line)
			return nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkDefaultText error:", err.Error())
	}

	return nil
}

// Number of leaders < 2 or > 5
func checkLeaderCount(filename string, d fs.DirEntry) error {
	if !strings.HasSuffix(filename, "leaders.md") {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	leaders := 0
	var emailRegex = regexp.MustCompile("[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*")

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "* ") && emailRegex.Match(scanner.Bytes()) {
			leaders++
		}

		line++
	}

	if leaders < 2 || leaders > 5 {
		fmt.Println(currChapter, "has", leaders, "leaders")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkLeaderCount error:", err.Error())
	}

	return nil
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
			fmt.Println(s)
		}
	}

}

// Meetup header present but no active Meetup for that chapter
func checkMeetupExists(filename string, d fs.DirEntry) error {
	if !strings.HasSuffix(filename, ".md") {
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
		if strings.Contains(scanner.Text(), "meetup-group:") {
			// let's grab the meetup group name and then validate it via Meetup API
		}
		line++
	}

	return nil
}

// Meetup header present and less than three meetings since Feb 2020
func checkPastMeetups() {

}

// Meetup header present and Link to Meetup in info.md but no metadata JavaScript for automated (warning)
func checkMeetupMissingMetaData(filename string, d fs.DirEntry) error {
	if !strings.HasSuffix(filename, ".md") {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	hasHeader := false
	hasJavaScript := false

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "meetup-group:") {
			hasHeader = true
		}

		if strings.Contains(scanner.Text(), "include chapter_events.html group=page.meetup-group") {
			hasJavaScript = true
		}

		line++
	}

	if hasHeader && !hasJavaScript {
		fmt.Println("Has Meetup metadata, but JavaScript is not present in", filename)
	}

	if !hasHeader && hasJavaScript {
		fmt.Println("No Meetup metadata, but JavaScript is present in", filename)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkMeetupMissingMetaData error:", err.Error())
	}

	return nil
}

// Old Wiki links are present (a warning not a breakage)
func checkForOldWiki(filename string, d fs.DirEntry) error {
	if !strings.Contains(filename, ".md") {
		return nil
	}

	if strings.Contains(filename, "migrated_content.md") {
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
			fmt.Println("Old wiki link found in ", filename, " on line ", line)
			fmt.Println(scanner.Text())
			return nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkForOldPolicy error:", err.Error())
	}

	return nil
}

func checkForDonate(filename string, d fs.DirEntry) error {

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
		if strings.Contains(scanner.Text(), "PayPal") || strings.Contains(scanner.Text(), "Paypal") {
			fmt.Println("Old donate mechanism in", filename, "on line", line)
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkForDonate error:", err.Error())
	}

	return nil
}

// Old policy links are present (a warning not a breakage)
func checkForOldPolicy(filename string, d fs.DirEntry) error {

	if !strings.Contains(filename, ".md") {
		return nil
	}

	if strings.Contains(filename, "migrated_content.md") {
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
		if strings.Contains(scanner.Text(), "Speaker_Agreement") {
			fmt.Println("Old Speaker Agreement in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "Conference_Policies") {
			fmt.Println("Old conference policy in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "Local_Chapter_Supporter") {
			fmt.Println("Old local chapter supporter policy in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "Chapter_Rules") {
			fmt.Println("Old local chapter rules in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "index.php/Membership") {
			fmt.Println("Old individual membership link in ", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "index.php/Corporate_Membership") {
			fmt.Println("Old corporate membership link in", filename, "on line ", line)
		}

		if strings.Contains(scanner.Text(), "OWASP_Project") {
			fmt.Println("Old projects link in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "About_OWASP") {
			fmt.Println("Old About OWASP link in", filename, "on line", line)
		}

		if strings.Contains(scanner.Text(), "docs.google.com/forms") ||
			strings.Contains(scanner.Text(), "goo.gl/forms") ||
			strings.Contains(scanner.Text(), "forms.gle") {
			fmt.Println("Google Forms link in ", filename, "on line", line)
			fmt.Println(scanner.Text())
		}

		re := regexp.MustCompile(`docs.google.com/a/.*/forms`)

		if re.Match(scanner.Bytes()) && strings.Contains(scanner.Text(), "owasp.org") {
			fmt.Println("Google Forms link in ", filename, "on line", line)
			fmt.Println(scanner.Text())
		}

		if re.Match(scanner.Bytes()) && !strings.Contains(scanner.Text(), "owasp.org") {
			fmt.Println("POLICY VIOLATION: Non-GDPR Google Forms link in ", filename, "on line", line)
			fmt.Println(scanner.Text())
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("checkForOldPolicy error:", err.Error())
	}

	return nil
}

// check if not meetup, then we manually look for other platforms (ConnPass, etc)
func checkNonAutomatedPlatforms() {
	// ConnPass, EventBrite, Facebook Groups, etc
}

// Tab filename and title metadata is incorrect
func checkTabTags() {

}

// Jekyll bundle fails to build
func checkJekyllBuilds() {

}

// Update git repos
func updateGit(s string, d fs.DirEntry) {

	// if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
	// 	fmt.Println("Updating ", d.Name())
	// 	cmd := exec.Command("git", "pull")
	// 	cmd.Dir = s
	// 	output, err := cmd.Output()
	// 	cmd.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%s", output)
	// }
}

func walk(s string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}

	// Directory Checks
	if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
		currChapter = d.Name()
		fmt.Println()
		fmt.Println("Scanning chapter ", currChapter)
		updateGit(s, d)
		checkIfSite(s, d)
	}

	// File Checks
	if !d.IsDir() {
		checkLeaderCount(s, d)
		checkForDonate(s, d)
		checkDefaultText(s, d)
		checkForOldPolicy(s, d)
		checkForOldWiki(s, d)
		checkDefaultExampleTab(s, d)
	}
	return nil
}

func main() {
	fmt.Println("OWASP Policy Scanner Tool")

	// client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnUrl))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)

	// chaptersDatabase := client.Database("chapters")
	// chaptersCollection := chaptersDatabase.Collection("chapters")

	// if chaptersCollection.Name() == "chapters" {
	// 	fmt.Println("Connected to MongoDB")
	// }

	filepath.WalkDir("chapters/", walk)
}
