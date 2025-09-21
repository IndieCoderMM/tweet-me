package main

import "html/template"

type TweetData struct {
    User      string
    Handle    string
    Text      string
    Timestamp string
    Retweets  string
    Quotes    string
    Likes     string
    Avatar    template.URL
    BodyClass string
}
