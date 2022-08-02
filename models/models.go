package models

import (
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"context"
)

type App struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}
