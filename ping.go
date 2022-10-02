package helper

import (
	"context"

	"cloud.google.com/go/spanner"
	ers "github.com/tys-muta/go-ers"
	"google.golang.org/api/iterator"
)

func Ping(ctx context.Context, client *spanner.Client) error {
	iter := client.Single().Query(ctx, spanner.Statement{SQL: "SELECT 1"})
	defer iter.Stop()

	row, err := iter.Next()
	if err == iterator.Done {
		return nil
	} else if err != nil {
		return ers.W(err)
	}

	var i int64
	if err := row.Columns(&i); err != nil {
		return ers.W(err)
	}

	return nil
}
