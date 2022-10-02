package txn

import (
	"context"
	"regexp"
	"sync"

	"cloud.google.com/go/spanner"
	"github.com/tys-muta/go-ers"
)

type contextKey string

type DebugLogger interface {
	Debugf(format string, args ...interface{})
}

const clientKey = contextKey("client")

var (
	loggerMu sync.Mutex
	logger   DebugLogger
)

type Txn interface {
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
}

type TxnHandle func(ctx context.Context) error

func SetLogger(v DebugLogger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger = v
}

func ReadOnly(ctx context.Context, client *spanner.Client, handle TxnHandle) error {
	if client != nil {
		rot := client.ReadOnlyTransaction()
		defer rot.Close()
		ctx = ctxWithTxn(ctx, rot)
	}

	if err := handle(ctx); err != nil {
		return ers.W(err)
	}

	return nil
}

func ReadWrite(ctx context.Context, client *spanner.Client, handle TxnHandle) error {
	if client != nil {
		if _, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
			ctx = ctxWithTxn(ctx, rwt)
			if err := handle(ctx); err != nil {
				return ers.W(err)
			}
			return nil
		}); err != nil {
			return ers.W(err)
		}
		return nil
	}

	if err := handle(ctx); err != nil {
		return ers.W(err)
	}

	return nil
}

func Transaction(ctx context.Context) (Txn, bool) {
	if ctx == nil {
		return nil, false
	}
	v := ctx.Value(clientKey)
	if v, ok := v.(Txn); ok {
		return v, true
	} else {
		return nil, false
	}
}

func Query(ctx context.Context, stmt spanner.Statement) (*spanner.RowIterator, error) {
	txn, ok := Transaction(ctx)
	if !ok {
		return nil, ers.ErrInternal.New("transaction is not found")
	}

	if logger != nil {
		logger.Debugf(
			"[SPANNER] SQL: %s\n\tParams: %v",
			regexp.MustCompile("[ |\n|\t]+").ReplaceAllString(stmt.SQL, " "),
			stmt.Params,
		)
	}

	return txn.Query(ctx, stmt), nil
}

func Read[Record any](ctx context.Context, stmt spanner.Statement) ([]*Record, error) {
	records := []*Record{}

	iter, err := Query(ctx, stmt)
	if err != nil {
		return nil, ers.W(err)
	}
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		record := new(Record)
		if err := row.ToStructLenient(record); err != nil {
			return ers.W(err)
		}
		records = append(records, record)
		return nil
	}); err != nil {
		return nil, ers.W(err)
	}

	return records, nil
}

func Write(ctx context.Context, stmt spanner.Statement) (affectedRows int64, err error) {
	iter, err := Query(ctx, stmt)
	if err != nil {
		return 0, ers.W(err)
	}
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return nil
	}); err != nil {
		return 0, ers.W(err)
	}

	return iter.RowCount, nil
}

func ctxWithTxn(ctx context.Context, txn Txn) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, clientKey, txn)
	return ctx
}
