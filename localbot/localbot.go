package main

import (
	"fmt"
	"github.com/scord/reddit-youtube-bot/bot"
)

//add your own API Key
const apiKey = ""

func main() {

	channel, subreddit, user, pass := "", "", "", ""

	fmt.Println("Enter Channel Username:")
	fmt.Scanf("%s", &channel)
	fmt.Println("Enter Subreddit")
	fmt.Scanf("%s", &subreddit)

	fmt.Println("Enter Bot Username:")
	fmt.Scanf("%s", &user)
	fmt.Println("Enter Bot Password")
	fmt.Scanf("%s", &pass)

	err := bot.Run(channel, subreddit, user, pass, apiKey)

	if err != nil {
		fmt.Println(err)
	}
}
