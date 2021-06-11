# owasp-policy-scanner

This quick and dirty tool is the start of an API backend. It checks OWASP GitHub Repos for policy requirements and leading practices, and produces a JSON file with the results. It also dumps to the screen, but I'm assuming this will be headless.

## Initial Setup 

### Clone repos 

This tool does a lot of the heavy lifting using OWASP's GitHub repos and essentially grepping them for issues or obtaining metadata.

```
chmod +x clonerepos.sh
./clonerepos.sh
```

This script queries GitHub's public search API sufficient times for up to 400 repos, and creates a file called "chapters_all.txt" and a folder called "chapters". Once all repos have been de-duped, it fires off git to clone all these repos. Once the repos are cloned, you can delete the chapters_all.txt file. 

### Get a GitHub key

GitHub APIs have a low number of API requests in a period before you get slowed down. You're gonna need a lot more. Login to your GitHub account, and obtain an oAuth token for API access, which will give you 5000 requests in an hour. You will need to copy this token somewhere safe like a Password Manager, because you're never gonna see it again. Do not check this token in, provide it via a command line switch. A future version of this tool will accept this value via an environment variable, but that's not currently implemented. 

### Compile the tool

Install Go from the usual places for your platform

```
go build
```

This produces a binary called "scanner". 

``` 
./scanner -help
Usage of ./scanner:
  -build
        Build Jekyll site (slow, may require super user privs)
  -chapter string
        Scan a single chapter
  -githubkey string
        Set a GitHub API access token
  -gitpull
        Update and force reset GitHub repos (slow) (default true)
  -meetup
        Show Meetup Group status (slow)
  -pages
        Show chapter page status
  -password string
        Meetup Password
  -policy
        Only show potential policy violations
  -username string
        Meetup Username
```

The tool doesn't use so many Meetup queries (yet) to need a Meetup API key, but it will pause when it runs out of requests. This pause is not long, so no message will be shown. If you run the tool A LOT, you will notice that GitHub forces the tool to sleep for up to 60 minutes at a time. So run it like once a day with the -meetup or -pages flag, otherwise the tool will never finish. 

## Usage

### Comprehensive scan with all the bells and whistles

This will take a LOT of time and need a GitHub API token. The results will be saved in the scanner_results.json file, but you can watch progress on the console or go make an espresso.

```
% ./scanner -githubkey xxxxxxxx -meetup -gitpull -pages
OWASP Policy Scanner Tool

Scanning chapter  www-chapter-london
Info: Updating www-chapter-london
Already up to date.
Info: GitHub Pages published for www-chapter-london
Info: .gitignore does not have _site in file chapters/www-chapter-london/.gitignore
Info: .gitignore does not have Gemfile.lock in file chapters/www-chapter-london/.gitignore
Info: Meetup OWASP-London exists, is active, 1202 members, 0 upcoming events, 25 past events
Info: Meetup metadata and JavaScript present
High: Old individual membership link in chapters/www-chapter-london/info.md on line 2
Low: Old wiki link found in chapters/www-chapter-london/info.md on line 2
High: Old donate mechanism in chapters/www-chapter-london/tab_pastevents.md on line 394
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 287
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 424
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 533
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 605
Low: Old wiki link found in chapters/www-chapter-london/tab_pastevents.md on line 287
```
### Running it on a single chapter

You can run the tool across a single chapter:

```
% ./scanner -githubkey xxxxxxxx -meetup -gitpull -pages -chapter www-chapter-london
OWASP Policy Scanner Tool

Scanning chapter  www-chapter-london
Info: Updating www-chapter-london
Already up to date.
Info: GitHub Pages published for www-chapter-london
Info: .gitignore does not have _site in file chapters/www-chapter-london/.gitignore
Info: .gitignore does not have Gemfile.lock in file chapters/www-chapter-london/.gitignore
Info: Meetup OWASP-London exists, is active, 1202 members, 0 upcoming events, 25 past events
Info: Meetup metadata and JavaScript present
High: Old individual membership link in chapters/www-chapter-london/info.md on line 2
Low: Old wiki link found in chapters/www-chapter-london/info.md on line 2
High: Old donate mechanism in chapters/www-chapter-london/tab_pastevents.md on line 394
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 287
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 424
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 533
High: Old conference policy in chapters/www-chapter-london/tab_pastevents.md on line 605
Low: Old wiki link found in chapters/www-chapter-london/tab_pastevents.md on line 287
```
### Just the policy, ma'am

The Chapter Policy contains things chapter leaders should be doing right. This is not easy with Jekyll and often people forget. So this flag just outputs policy requirements. 

```
% ./scanner -githubkey xxxxxx -meetup -gitpull -pages -chapter www-chapter-ankara -policy 
OWASP Policy Scanner Tool

Scanning chapter  www-chapter-ankara
Already up to date.
POLICY: GitHub Pages are disabled for www-chapter-ankara
POLICY: www-chapter-ankara has 1 leaders
POLICY: Meetup Group does not exist for OWASP-Ankara-Chapter
POLICY: www-chapter-ankara has 0 leaders
```

### Quick and Dirty Incremental scan

Run the tool with no flags

```
% ./scanner
[ lots and lots of output truncated]
```

This mode does not update the GitHub repos, and makes no API calls, so it's fast. 






