package bgs

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

// BGS config init
type BGS struct {
	Proj    string
	Skc     string
	Dataset string
	Table   string
	c       *bigquery.Client
}

// New new bigquery stream
func New(proj, skc, table, dataset string) *BGS {
	bg := &BGS{Proj: proj, Skc: skc, Table: table, Dataset: dataset}

	return bg.connect()
}

func (bgs *BGS) connect() *BGS {
	ctx := context.Background()
	opt := option.WithCredentialsFile(bgs.Skc)
	client, err := bigquery.NewClient(ctx, bgs.Proj, opt)
	if err != nil {
		fmt.Println("[bigquery]:", err)
	}
	bgs.c = client

	return bgs
}

// AddSchema create table and schema
func (bgs *BGS) AddSchema(schm interface{}) *BGS {

	dtset := bgs.c.Dataset(bgs.Dataset)
	if err := dtset.Create(context.TODO(), nil); err != nil {
		fmt.Println("[bigquery]:", err)
	}

	schema2, err := bigquery.InferSchema(schm)
	if err != nil {
		fmt.Println("[bigquery] inter schema2:", err)
	}

	err = dtset.Table(bgs.Table).Create(context.TODO(), &bigquery.TableMetadata{
		Schema: schema2,
	})
	if err != nil {
		fmt.Println("[bigquery] ", err)
	}
	return bgs
}

// AddRow put to table
func (bgs *BGS) AddRow(row interface{}) error {
	dtset := bgs.c.Dataset(bgs.Dataset)
	u := dtset.Table(bgs.Table).Uploader()
	err := u.Put(context.TODO(), row)
	if err != nil {
		fmt.Println("[bigquery] row: ", err)
	}
	return err
}
