package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"github.com/fnproject/fdk-go"
	"io"
	"net/url"
	"os"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	r := new(api.RequestPayload)

	err := json.NewDecoder(in).Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode STDIN, got error %v", err.Error())
		return
	}

	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()

	ok, err := twitterAPI.VerifyCredentials()
	if !ok {
		fmt.Fprint(os.Stderr, "unable to verify credentials\n")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to verify createndials, got error: %v\n", err.Error())
		return
	}
	v := url.Values{}
	v.Set("in_reply_to_status_id", r.TweetID)
	msg := fmt.Sprintf(`Hey, %v! Had to admit, i don't know where is it ¯\_(ツ)_/¯`, r.User)
	_, err = twitterAPI.PostTweet(msg, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to post tweet for %v, got error %v", r.User, err.Error())
		return
	}
}
