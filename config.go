package main

import (
	"flag"
)

type configT struct {
	build           bool
	chapter         string
	gitPull         bool
	githubkey       string
	meetup          bool
	meetup_password string
	meetup_username string
	pages           bool
	policy          bool
}

var config configT

func processFlags() {
	flag.BoolVar(&config.build, "build", config.build, "Build Jekyll site (slow, may require super user privs)")
	flag.BoolVar(&config.gitPull, "gitpull", config.gitPull, "Update and force reset GitHub repos (slow)")
	flag.StringVar(&config.githubkey, "githubkey", config.githubkey, "Set a GitHub API access token")
	flag.BoolVar(&config.meetup, "meetup", config.meetup, "Show Meetup Group status (slow)")
	flag.BoolVar(&config.pages, "pages", config.pages, "Show chapter page status")
	flag.BoolVar(&config.policy, "policy", config.policy, "Only show potential policy violations")
	flag.StringVar(&config.chapter, "chapter", config.chapter, "Scan a single chapter")
	flag.StringVar(&config.meetup_password, "password", config.meetup_password, "Meetup Password")
	flag.StringVar(&config.meetup_username, "username", config.meetup_username, "Meetup Username")
	flag.Parse()
}

func loadConfig() configT {

	config = configT{}

	config.gitPull = true
	config.build = false
	config.meetup = false
	config.pages = false
	config.policy = false

	return config
}
