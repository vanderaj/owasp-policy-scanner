package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var currChapter string

type StatusLevelT int

const (
	Info   StatusLevelT = 0
	Low                 = 1
	Medium              = 2
	High                = 3
	Policy              = 4
)

type privacyStatusT int

const (
	notpresent     privacyStatusT = 0
	unknown                       = 1
	owasp                         = 2
	gdpr_violation                = 3
)

type serviceStatusT int

const (
	nonexistant serviceStatusT = 0
	inactive                   = 1
	active                     = 2
)

type chapterStatusT struct {
	AutoMigration          bool
	ConfigYml              bool
	DefaultText            bool
	ExampleTab             bool
	GitHub                 serviceStatusT
	GoogleForms            privacyStatusT
	Leaders                int
	Meetup                 serviceStatusT
	MeetupMetaData         serviceStatusT
	MeetupName             string
	MeetupPastMeetings     int
	MeetupUpcomingMeetings int
	OldDonate              bool
	OldGitIgnore           bool
	OldLink                bool
	OldPolicy              bool
	OldProjects            bool
	OldSpeaker             bool
	OldWiki                bool
	SitePresent            bool
}

var chapterStatus = map[string]*chapterStatusT{}

func writeJSON() {

	file, err := json.MarshalIndent(chapterStatus, "", " ")
	if err != nil {
		println("Error marshalling chapterStatus")
		return
	}
	err = ioutil.WriteFile("scanner_output.json", file, 0644)
	if err != nil {
		println("Error writing JSON to disk")
		return
	}
}

func printStatus(sl StatusLevelT, s string) {
	if config.policy && sl < Policy {
		return
	}

	switch sl {
	case Info:
		fmt.Print("Info: ")
		break

	case Low:
		fmt.Print("Low: ")
		break

	case Medium:
		fmt.Print("Medium: ")
		break

	case High:
		fmt.Print("High: ")
		break

	case Policy:
		fmt.Print("POLICY: ")
		break

	}

	fmt.Println(s)
}

// Out of date dependencies in _config.yml
func checkConfigYml(s string, d fs.DirEntry) {

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
			printStatus(Policy, fmt.Sprintf("Default text present in %s on line %d", filename, line))
			chapterStatus[currChapter].DefaultText = true
			return nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkDefaultText error: "+err.Error())
	}

	return nil
}

// Default tab tab_example.md is present
func checkDefaultExampleTab(filename string, d fs.DirEntry) error {
	if strings.Contains(filename, "tab_example.md") {
		printStatus(Low, "Example tab found at: "+filename)
		chapterStatus[currChapter].ExampleTab = true
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
			printStatus(Policy, fmt.Sprintf("Auto-Migration Headers active in %s on line %d", filename, line))
			chapterStatus[currChapter].AutoMigration = true
			return nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkDefaultMigrationHeader error: "+err.Error())
	}

	return nil
}

// Number of leaders < 2 or > 5
func checkLeaderCount(filename string, d fs.DirEntry) error {
	if !strings.HasSuffix(filename, "leaders.md") {
		return nil
	}

	// the leaders tab is not the official source of leadership information
	if strings.HasSuffix(filename, "tab_leaders.md") {
		return nil
	}

	// Check that only the top leaders.md file is parsed
	// Only www-chapter-<chaptername>/leaders.md counts as a leader
	reg := regexp.MustCompile(`chapters\/www-chapter-.*\/(.*)(/)leaders.md`)
	if reg.MatchString(filename) {
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

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {

		email := scanner.Text()
		email = strings.TrimLeft(email, "* ")
		email = strings.TrimLeft(email, "- ")
		email = strings.TrimSpace(email)

		// skip blank lines
		if email == "" {
			continue
		}

		// Skip headers
		if strings.Contains(email, "#") {
			continue
		}

		// Skip lines without an @ symbol
		if !strings.Contains(email, "@") {
			continue
		}

		// Check to see if it's a Markdown mailto: and trim accordingly
		if strings.Contains(email, "[") && strings.Contains(email, "mailto") {
			// trim to just the email address bit
			reg := regexp.MustCompile(`^.*mailto:`)
			email = reg.ReplaceAllString(email, "${1}")
			reg = regexp.MustCompile(`\).*$`)
			email = reg.ReplaceAllLiteralString(email, "")
			email = strings.Replace(email, "/", "", -1)
		}

		// Check to see if it's a Mardown without a mailto
		if strings.Contains(email, "[") && !strings.Contains(email, "mailto") {
			// trim to just the email address bit
			reg := regexp.MustCompile(`^.*\(`)
			email = reg.ReplaceAllString(email, "${1}")
			reg = regexp.MustCompile(`\).*$`)
			email = reg.ReplaceAllLiteralString(email, "")
			email = strings.Replace(email, "/", "", -1)
		}

		// Some emails are like this "* <demo text>" but aren't real

		if strings.HasPrefix(email, "* ") || strings.HasPrefix(email, "- ") {
			email = strings.Replace(email, "* ", "", 1)
			email = strings.Replace(email, "- ", "", 1)
		}

		email = strings.TrimSpace(email)

		_, err := mail.ParseAddress(email)
		if err != nil {
			if err.Error() == "mail: no angle-addr" {
				printStatus(Low, "checkLeaderCount leader has no email: "+email)
			} else {
				printStatus(Info, "checkLeaderCount error: "+err.Error())
			}
		} else {
			leaders++
		}
	}

	if leaders < 2 || leaders > 5 {
		printStatus(Policy, fmt.Sprintf("%s has %d leaders", currChapter, leaders))
		chapterStatus[currChapter].Leaders = leaders
	}

	chapterStatus[currChapter].Leaders = leaders

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkLeaderCount error: "+err.Error())
	}

	return nil
}

// Leaders in leaders.md doesn’t match Copper
func checkLeadersInCopper(s string, d fs.DirEntry) {
}

// Out of date .gitignore
func checkOldGitIgnore(filename string, d fs.DirEntry) error {
	if !strings.HasSuffix(filename, ".gitignore") {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	hasSite := false
	hasGemfile := false

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "_site") {
			hasSite = true
		}

		if strings.Contains(scanner.Text(), "Gemfile.lock") {
			hasGemfile = true
		}

		line++
	}

	if hasSite {
		printStatus(Info, ".gitignore does not have _site in file "+filename)
		chapterStatus[currChapter].OldGitIgnore = true
	}

	if hasGemfile {
		printStatus(Info, ".gitignore does not have Gemfile.lock in file "+filename)
		chapterStatus[currChapter].OldGitIgnore = true
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkMeetupMissingMetaData error: "+err.Error())
	}

	return nil
}

// _site/ being present
func checkIfSite(s string, d fs.DirEntry) {
	if d.IsDir() {
		if d.Name() == "_site" {
			printStatus(Low, "Site directory is present at "+s)
			chapterStatus[currChapter].SitePresent = true
		}
	}
}

type PagesRespT struct {
	Has_pages bool `json:"has_pages"`
}

func checkPagesStatus(chapterName string) error {
	if !config.pages {
		return nil
	}

	// check the group is exists and active
	reqUrl := fmt.Sprintf("https://api.github.com/repos/OWASP/%s", chapterName)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+config.githubkey)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		printStatus(Policy, "GitHub Pages does not exist for "+chapterName)
		chapterStatus[currChapter].GitHub = nonexistant
		return nil
	}

	if resp.StatusCode == 410 {
		printStatus(Policy, "GitHub Pages exists, but is disabled for "+chapterName)
		chapterStatus[currChapter].GitHub = inactive
		return nil
	}

	requestsLeft, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	requestsTimeOut, _ := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)

	// Sleep if necessary to slow things down
	if requestsLeft < 1 {
		resetTime := requestsTimeOut - time.Now().Unix()
		printStatus(Info, fmt.Sprintf("GitHub API limit reached, sleeping for %d seconds", resetTime))
		time.Sleep(time.Duration(resetTime) * time.Second)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var m PagesRespT
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		log.Fatalln(err)
	}

	if m.Has_pages {
		printStatus(Info, "GitHub Pages published for "+chapterName)
		chapterStatus[currChapter].GitHub = active
	} else {
		printStatus(Policy, "GitHub Pages are disabled for "+chapterName)
		chapterStatus[currChapter].GitHub = inactive
	}

	return nil
}

type meetupGroupRespT struct {
	Status               string `json:"status"`
	Upcoming_event_count int    `json:"upcoming_event_count"`
	Past_event_count     int    `json:"past_event_count"`
	Members              int    `json:"members"`
}

// Meetup header present but no active Meetup for that chapter
func checkMeetupExists(filename string, d fs.DirEntry) error {
	if !config.meetup {
		return nil
	}

	if !strings.HasSuffix(filename, "index.md") {
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
		lineStr := scanner.Text()
		if strings.HasPrefix(lineStr, "meetup-group:") {
			// let's grab the meetup group name and then validate it via Meetup API

			meetupGroup := strings.Split(lineStr, ": ")

			if len(meetupGroup) == 1 {
				printStatus(Policy, "Meetup-group header is present but blank")
				chapterStatus[currChapter].Meetup = nonexistant
				return nil
			}

			if strings.Trim(meetupGroup[1], " ") == "" {
				printStatus(Policy, "Meetup-group header is present but blank with whitespace")
				chapterStatus[currChapter].Meetup = nonexistant
				return nil
			}

			// check the group is exists and active
			resp, err := http.Get("https://api.meetup.com/" + meetupGroup[1] + "?fields=past_event_count,upcoming_event_count")
			if err != nil {
				log.Fatalln(err)
			}

			if resp.StatusCode == 404 {
				printStatus(Policy, "Meetup Group does not exist for "+meetupGroup[1])
				chapterStatus[currChapter].Meetup = nonexistant
				return nil
			}

			if resp.StatusCode == 410 {
				printStatus(Policy, "Meetup exists, but is disabled for "+meetupGroup[1])
				chapterStatus[currChapter].Meetup = inactive
				return nil
			}

			requestsLeft, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
			requestsTimeOut, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))

			// We could go to 0, but let's not
			if requestsLeft == 1 {
				time.Sleep(time.Duration(requestsTimeOut) * time.Second)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			var m meetupGroupRespT
			err = json.Unmarshal([]byte(body), &m)
			if err != nil {
				log.Fatalln(err)
			}

			if m.Status == "active" && m.Past_event_count < 3 {
				printStatus(Policy, fmt.Sprintf("Low past meetings. Meetup %s exists, is active, %d members, %d upcoming events, %d past events", meetupGroup[1], m.Members, m.Upcoming_event_count, m.Past_event_count))
				chapterStatus[currChapter].Meetup = active
				chapterStatus[currChapter].MeetupName = meetupGroup[1]
				chapterStatus[currChapter].MeetupPastMeetings = m.Past_event_count
				chapterStatus[currChapter].MeetupUpcomingMeetings = m.Upcoming_event_count
				return nil
			}

			if m.Status == "active" {
				printStatus(Info, fmt.Sprintf("Meetup %s exists, is active, %d members, %d upcoming events, %d past events", meetupGroup[1], m.Members, m.Upcoming_event_count, m.Past_event_count))
				chapterStatus[currChapter].Meetup = active
				chapterStatus[currChapter].MeetupName = meetupGroup[1]
				chapterStatus[currChapter].MeetupPastMeetings = m.Past_event_count
				chapterStatus[currChapter].MeetupUpcomingMeetings = m.Upcoming_event_count
				return nil
			}

			printStatus(Info, "DEBUG Meetup Unknown Status: "+m.Status)
		}
		line++
	}

	return nil
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

	chapterStatus[currChapter].MeetupMetaData = nonexistant
	if hasHeader && !hasJavaScript {
		printStatus(Medium, "Has Meetup metadata, but JavaScript is not present in "+filename)
		chapterStatus[currChapter].MeetupMetaData = inactive
	}

	if !hasHeader && hasJavaScript {
		printStatus(Medium, "No Meetup metadata, but JavaScript is present in "+filename)
		chapterStatus[currChapter].MeetupMetaData = inactive
	}

	if hasHeader && hasJavaScript {
		printStatus(Info, "Meetup metadata and JavaScript present")
		chapterStatus[currChapter].MeetupMetaData = active
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkMeetupMissingMetaData error: "+err.Error())
	}

	return nil
}

// Old Wiki links are present (a warning not a breakage)
func checkForOldWiki(filename string, d fs.DirEntry) error {
	if !strings.Contains(filename, ".md") {
		return nil
	}

	// If it's in the migrated content, don't care as it's not shown
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
		if strings.Contains(scanner.Text(), "www.owasp.org/index.php") {
			if !config.policy {
				printStatus(Low, fmt.Sprintf("Old wiki link found in %s on line %d", filename, line))
				chapterStatus[currChapter].OldWiki = true
			}
			return nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkForOldPolicy error: "+err.Error())
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
			printStatus(High, fmt.Sprintf("Old donate mechanism in %s on line %d", filename, line))
			chapterStatus[currChapter].OldDonate = true
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkForDonate error: "+err.Error())
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
			printStatus(High, fmt.Sprintf("Old Speaker Agreement in %s on line %d", filename, line))
			chapterStatus[currChapter].OldSpeaker = true
		}

		if strings.Contains(scanner.Text(), "Conference_Policies") {
			printStatus(High, fmt.Sprintf("Old conference policy in %s on line %d", filename, line))
			chapterStatus[currChapter].OldPolicy = true
		}

		if strings.Contains(scanner.Text(), "Local_Chapter_Supporter") {
			printStatus(High, fmt.Sprintf("Old local chapter supporter policy in %s on line %d", filename, line))
			chapterStatus[currChapter].OldPolicy = true
		}

		if strings.Contains(scanner.Text(), "Chapter_Rules") || strings.Contains(scanner.Text(), "Chapter_Handbook") {
			printStatus(High, fmt.Sprintf("Old local chapter rules or handbook in %s on line %d", filename, line))
			chapterStatus[currChapter].OldPolicy = true
		}

		if strings.Contains(scanner.Text(), "index.php/Membership") {
			printStatus(High, fmt.Sprintf("Old individual membership link in %s on line %d", filename, line))
			chapterStatus[currChapter].OldLink = true
		}

		if strings.Contains(scanner.Text(), "index.php/Corporate_Membership") {
			printStatus(High, fmt.Sprintf("Old corporate membership link in %s on line %d", filename, line))
			chapterStatus[currChapter].OldLink = true
		}

		if strings.Contains(scanner.Text(), "OWASP_Project") {
			printStatus(Low, fmt.Sprintf("Old projects link in %s on line %d", filename, line))
			chapterStatus[currChapter].OldLink = true
		}

		if strings.Contains(scanner.Text(), "About_OWASP") {
			printStatus(Low, fmt.Sprintf("Old About OWASP link in %s on line %d", filename, line))
			chapterStatus[currChapter].OldLink = true
		}

		re := regexp.MustCompile(`docs.google.com/a/.*/forms`)

		if strings.Contains(scanner.Text(), "docs.google.com/forms") ||
			strings.Contains(scanner.Text(), "goo.gl/forms") ||
			strings.Contains(scanner.Text(), "forms.gle") {
			printStatus(High, fmt.Sprintf("High: Google Forms link in %s on line %d", filename, line))
			chapterStatus[currChapter].GoogleForms = unknown
		}

		if re.Match(scanner.Bytes()) && strings.Contains(scanner.Text(), "owasp.org") {
			printStatus(High, fmt.Sprintf("Info: OWASP Google Forms link in %s on line %d", filename, line))
			chapterStatus[currChapter].GoogleForms = owasp
		}

		if re.Match(scanner.Bytes()) && !strings.Contains(scanner.Text(), "owasp.org") {
			printStatus(Policy, fmt.Sprintf("Non-GDPR Google Forms link in %s on line %d", filename, line))
			chapterStatus[currChapter].GoogleForms = gdpr_violation
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		printStatus(Info, "checkForOldPolicy error: "+err.Error())
	}

	return nil
}

// check if not meetup, then we manually look for other platforms (ConnPass, etc)
func checkNonAutomatedPlatforms(s string, d fs.DirEntry) {
	// ConnPass, EventBrite, Facebook Groups, etc
}

// Tab filename and title metadata is incorrect
func checkTabTags(s string, d fs.DirEntry) {
	// find the tag in index.md

	// if no tag, but tab_files exist, display an error and exit

	// search all tab_name.md files for tags

	// show any tabs that aren't tagged correctly

}

// Jekyll bundle fails to build
func checkJekyllBuilds(s string, d fs.DirEntry) {
	if config.build {
		if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
			printStatus(Info, "Building "+d.Name())
			cmd := exec.Command("bundle", "install")
			cmd.Dir = s
			output, err := cmd.Output()
			cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", output)

			cmd = exec.Command("bundle", "exec jekyll serve")
			cmd.Dir = s
			output, err = cmd.Output()
			cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", output)

		}
	}
}

// Update git repos
func updateGit(s string, d fs.DirEntry) {

	if config.gitPull {
		if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
			printStatus(Info, "Updating "+d.Name())
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
}

var dirsInspected int = 0

func walk(s string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}

	// If we are only processing one chapter, let's only do that one
	if len(config.chapter) > 0 && !strings.Contains(s, config.chapter) {
		return nil
	}

	// Directory Checks
	if d.IsDir() && strings.HasPrefix(d.Name(), "www-chapter") {
		currChapter = d.Name()
		blankStatus := &chapterStatusT{}
		chapterStatus[currChapter] = blankStatus
		fmt.Println()
		fmt.Println("Scanning chapter ", currChapter)
		updateGit(s, d)
		checkPagesStatus(currChapter)
		checkIfSite(s, d)
		checkJekyllBuilds(s, d)
		dirsInspected++
	}

	// File Checks
	if !d.IsDir() {
		checkLeaderCount(s, d)
		checkMeetupExists(s, d)
		checkMeetupMissingMetaData(s, d)
		checkLeadersInCopper(s, d)
		checkDefaultMigrationHeader(s, d)
		checkDefaultText(s, d)
		checkDefaultExampleTab(s, d)
		checkConfigYml(s, d)
		checkForDonate(s, d)
		checkForOldPolicy(s, d)
		checkForOldWiki(s, d)
		checkOldGitIgnore(s, d)
		checkNonAutomatedPlatforms(s, d)
		checkTabTags(s, d)
	}
	return nil
}

func main() {
	fmt.Println("OWASP Policy Scanner Tool")

	chapterStatus = make(map[string]*chapterStatusT)

	config = loadConfig()
	processFlags()

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

	if dirsInspected == 0 {
		fmt.Println("No chapters scanned")
	} else {
		writeJSON()
	}
}
