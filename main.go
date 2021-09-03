package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/alecthomas/kong"
	"google.golang.org/api/iterator"
)

type Entity struct {
	TimeStamp int64
	Value     string
	Key       *datastore.Key
}

type Globals struct {
	Kind      string `arg name:"kind" help:"Firestore Entity Kind."`
	ProjectId string `help:"GCP Project ID. Overrides value taken from environment variable GCP_PROJECT_ID."`
}

type GetCmd struct {
	Globals
	Key string `arg name:"key" help:"Firestory Entity Key name" type:"key"`
}

func (gc *GetCmd) Run(projectId string) error {
	if gc.ProjectId == "" && projectId == "" {
		return fmt.Errorf("project ID is not set")
	}
	if gc.ProjectId == "" {
		gc.ProjectId = projectId
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, gc.ProjectId)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()
	dbKey := datastore.NameKey(gc.Kind, gc.Key, nil)
	entity := Entity{}

	if err := client.Get(ctx, dbKey, &entity); err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	fmt.Println(entity.Value)
	return nil
}

type PutCmd struct {
	Globals
	Key   string `arg name:"key" help:"Firestory Entity Key name."`
	Value string `arg name:"value" help:"Value to store."`
}

func (pc *PutCmd) Run(projectId string) error {
	if pc.ProjectId == "" && projectId == "" {
		return fmt.Errorf("project ID is not set")
	}
	if pc.ProjectId == "" {
		pc.ProjectId = projectId
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, pc.ProjectId)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()
	dbKey := datastore.NameKey(pc.Kind, pc.Key, nil)
	entity := Entity{
		TimeStamp: time.Now().Unix(),
		Value:     pc.Value,
		Key:       dbKey,
	}

	if _, err := client.Put(ctx, dbKey, &entity); err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	return nil
}

type CleanCmd struct {
	Globals
	Timestamp string `arg name:"timestamp" help:"Unix epoch time." type:"timestamp"`
}

func (cc *CleanCmd) Run(projectId string) error {
	if cc.ProjectId == "" && projectId == "" {
		return fmt.Errorf("project ID is not set")
	}
	if cc.ProjectId == "" {
		cc.ProjectId = projectId
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, cc.ProjectId)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	query := datastore.NewQuery(cc.Kind).Filter("TimeStamp <", cc.Timestamp)
	entities := client.Run(ctx, query)
	var keys []*datastore.Key
	for {
		var e Entity
		key, err := entities.Next(&e)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		return err
	}
	return nil
}

type Cli struct {
	Get   GetCmd   `cmd help:"Get a Firestore entity of a Kind with a provided Key."`
	Put   PutCmd   `cmd help:"Store a Firestore Entity of a Kind with provided Key and Value."`
	Clean CleanCmd `cmd help:"Clean all Firestore Entities of a Kind created/updated before provided timestamp."`
}

func main() {
	cli := Cli{}
	ctx := kong.Parse(&cli,
		kong.Description("CLI tool to use Google Firestore/Datastore as a key/value storage."),
		kong.UsageOnError(),
	)
	projectId := os.Getenv("GCP_PROJECT_ID")
	err := ctx.Run(projectId)
	ctx.FatalIfErrorf(err)
}
