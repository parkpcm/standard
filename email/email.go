package email

import (
	"context"
	"encoding/json"
	"fmt"
	mg "github.com/mailgun/mailgun-go/v4"
	"github.com/parkpcm/standard/secret"
	"os"
)

type clientCreds struct {
	Domain string `json:"mailgun_domain"`
	Key    string `json:"mailgun_key"`
}

type Client struct {
	Domain     string
	apiKey     string
	credLoaded bool
	DataPath   string
}

func (i *Client) Client(ctx context.Context) (*mg.MailgunImpl, error) {

	if err := i.load(ctx); err != nil {
		return nil, err
	}

	return mg.NewMailgun(i.Domain, i.apiKey), nil
}

func (i *Client) load(ctx context.Context) error {

	var data []byte
	var err error

	if !i.credLoaded {

		if len(os.Getenv("SECRET_PATH")) > 0 {
			if data, err = secret.Get(ctx, os.Getenv("SECRET_PATH")); err != nil {
				return fmt.Errorf("%w error getting data from SECRET_PATH", err)
			}
		} else if len(i.DataPath) > 0 {
			if data, err = secret.GetFromVolume(i.DataPath); err != nil {
				return fmt.Errorf("%w error getting data from %s", err, i.DataPath)
			}
		} else {
			return fmt.Errorf("unable to find a credentials source")
		}

		mi := clientCreds{}

		if err = json.Unmarshal(data, &mi); err != nil {
			return fmt.Errorf("%w error loading data into information", err)
		}

		if len(mi.Domain) == 0 {
			return fmt.Errorf("missing mailgun domain")
		}

		if len(mi.Key) == 0 {
			return fmt.Errorf("missing mailgun API key")
		}

		i.apiKey = mi.Key
		i.Domain = mi.Domain

		i.credLoaded = true
	}

	return nil
}
