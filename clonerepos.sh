#!/bin/sh

# Find a list of all chapters, a page at a time with these queries:

# First page (Chapter 1-100)
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=1

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c1.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=1'

# Second page (chapters 101-200)
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=2

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c2.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=2'

# Third page (chapters 201-300)
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=3

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c3.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=3'

# Fourth page (Chapters 301-400)
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=4

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c4.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=4'

# Get a list of all raw chapters from that json file:

grep -Po '(?<="html_url": ")[^"]*' c1.json > chapters_raw.txt
grep -Po '(?<="html_url": ")[^"]*' c2.json >> chapters_raw.txt
grep -Po '(?<="html_url": ")[^"]*' c3.json >> chapters_raw.txt
grep -Po '(?<="html_url": ")[^"]*' c4.json >> chapters_raw.txt

# Remove dupes
grep www-chapter chapters_raw.txt | sort -u > chapters_all.txt
rm c?.json chapters_raw.txt

mkdir -p chapters
cd chapters

# Clone the repos
while read IN; do
  # Don't process blank lines
  if [ ! -z "$IN" ] 
  then 
    CHAPTER=${IN#"https://github.com/OWASP/"}
    CHAPTER=${CHAPTER%"/"}
    # Check to see if we've already cloned the repo
    if [ -d $CHAPTER ]
    then
      # repo already cloned, let's reset & pull it
      echo "Updating $CHAPTER"
      cd $CHAPTER
      git fetch --all

      MAIN=false
      for BRANCH in `git branch | grep "[^* ]+" -Eo`;
      do 
        if [[ "$BRANCH" == *"main"* ]]; then
          echo "Main branch found: $BRANCH"
          MAIN=true
        fi      
      done

      if [ $MAIN == true ]
      then
        echo "Resetting to main"
        # Newer repos and those that are fixed now use main
        git reset --hard origin/main || true
        # Switch to main in case it has master & main
        git checkout main || true
      else
        # Older repos use "master", grab that first
        git reset --hard origin/master || true
      fi

      # Update the contents to head
      git config pull.rebase true 
      git pull
      cd ..
    else
      # new repo, which will be the pristine state we expect
      echo "Cloning $IN"
      git clone "$IN"
    fi
  fi
done < ../chapters_all.txt

cd ..
