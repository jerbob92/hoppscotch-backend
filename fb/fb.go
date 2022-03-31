package fb

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var FBApp *firebase.App
var FBAuth *auth.Client

func Initialize() error {
	opt := option.WithCredentialsFile(viper.GetString("firebase.serviceAccountFile"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	FBApp = app
	authClient, err := FBApp.Auth(context.Background())
	if err != nil {
		return err
	}

	FBAuth = authClient

	return nil
}
