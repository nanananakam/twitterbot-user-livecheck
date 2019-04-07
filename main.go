package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"math"
	"net/url"
	"os"
	"time"
)

func main() {
	lambda.Start(Main)
}

func Main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	twitterApi := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	v := url.Values{}
	v.Set("screen_name", os.Getenv("TWITTER_TARGET_SCREEN_NAME"))
	v.Set("count", "1")
	tweets, err := twitterApi.GetUserTimeline(v)
	if err != nil {
		panic(err)
	}
	lastTweetTime, err := tweets[0].CreatedAtTime()
	if err != nil {
		panic(err)
	}
	if lastTweetTime.Add(time.Duration(12) * time.Hour).Before(time.Now()) {
		duration := time.Now().Sub(lastTweetTime)
		tweet := fmt.Sprintf(
			".@%s 生きてますか～（前回のツイートから約%d.%d時間経過しています",
			os.Getenv("TWITTER_TARGET_SCREEN_NAME"),
			int(math.Floor(duration.Hours())),
			int(math.Ceil(duration.Minutes()))%60*10/60,
		)
		if _, err := twitterApi.PostTweet(tweet, nil); err != nil {
			panic(err)
		}
	}
}
