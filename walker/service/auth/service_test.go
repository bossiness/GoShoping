package auth_test

import (
	"os"
	"testing"

	"btdxcx.com/walker/model"
	"btdxcx.com/walker/repository"

	"btdxcx.com/walker/repository/account"
	"btdxcx.com/walker/repository/token"
	"btdxcx.com/walker/service/auth"
	"golang.org/x/net/context"
)

const (
	defaultHost = "localhost:27017"
)

func TestAuthService(t *testing.T) {

	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := repository.CreateSession(host)

	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	defer session.Close()

	if err != nil {
		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error.
		t.Fatal("Could not connect to datastore with host ", host, err)
	}

	arepo := account.Repository{Session: session}
	trepo := token.Repository{Session: session}
	auth := auth.Service{
		ARepo: arepo,
		TRepo: trepo,
	}

	out := model.Token{}
	if err := auth.Create(context.TODO(), &model.AuthRequest{}, &out); err != nil {
		t.Fatal("expected: auth Create error:", err)
	}

	

}
