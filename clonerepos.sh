#!/bin/sh

# Find a list of all chapters, a page at a time with these queries:

# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=1
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=2
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=3
# https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=4

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c1.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=1'

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c2.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=2'

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c3.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=3'

curl \
  -H "Accept: application/vnd.github.v3+json" \
  -o c4.json \
  'https://api.github.com/search/repositories?q=in%3Aname+www-chapter+org%3AOWASP&per_page=100&page=4'

# Get a list of all raw chapters from that json file:
# grep -Po '(?<="html_url":")[^"]*' c1.json > chapters_raw.txt

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
while read in; do git clone "$in"; done < ../chapters_all.txt

cd ..
