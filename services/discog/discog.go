package discog

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mrjones/oauth"
)

const (
	consumerKey    = "qSrNiVZkWynmCmHddibm"
	consumerSecret = "yzUhZJQHHNCVEeNnWeTNkvxIfsuYvNEB"
)

type Client struct {
	c            *http.Client
	token        *oauth.AccessToken
	consumer     *oauth.Consumer
	requestToken *oauth.RequestToken
}

func New(ctx context.Context) (*Client, error) {
	d := &Client{}

	return d, nil
}

func (dc *Client) Authenticate(ctx context.Context) (string, error) {
	url, err := dc.authenticate(ctx)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (dc *Client) authenticate(ctx context.Context) (string, error) {
	dc.consumer = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.discogs.com/oauth/request_token",
			AuthorizeTokenUrl: "https://www.discogs.com/oauth/authorize",
			AccessTokenUrl:    "https://api.discogs.com/oauth/access_token",
		})
	var u string
	var err error
	dc.requestToken, u, err = dc.consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		return "", err
	}

	fmt.Println("(1) Go to: " + u)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	// verificationCode := ""
	// fmt.Scanln(&verificationCode)

	// access_token, err := dc.consumer.AuthorizeToken(requestToken, verificationCode)
	// if err != nil {
	// 	return err
	// }

	// return dc.makeClient(access_token)
	return u, nil
}

func (dc *Client) Register(verificationCode string) error {
	access_token, err := dc.consumer.AuthorizeToken(dc.requestToken, verificationCode)
	if err != nil {
		return err
	}

	return dc.makeClient(access_token)

}

func (dc *Client) makeClient(access_token *oauth.AccessToken) error {
	// Make a request to the authenticated endpoint
	var err error
	dc.c, err = dc.consumer.MakeHttpClient(access_token)
	if err != nil {
		return err
	}
	return nil
}

func (dc *Client) GetUser() error {
	resp, err := dc.c.Get("https://api.discogs.com/users/strawhat_sunny/collection/folders/0/releases")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	return nil
}
