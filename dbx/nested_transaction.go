package dbx

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type contextKey struct {
	name string
}

var nestedTransactionCtxKey = &contextKey{"db.nestedTransaction"}
var transactionCtxKey = &contextKey{"db.transaction"}

func GetNestedContextTransaction(ctx context.Context) *NestedTransaction {
	if ntx, ok := ctx.Value(nestedTransactionCtxKey).(*NestedTransaction); ok {
		return ntx
	}

	return nil
}

func GetContextTransaction(ctx context.Context) *sqlx.Tx {
	if ctx == nil {
		return nil
	}

	if tx, ok := ctx.Value(transactionCtxKey).(*sqlx.Tx); ok {
		return tx
	}

	return nil
}

func BeginContextTransaction(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, context.Context, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, nil, err
	}

	return tx, context.WithValue(ctx, transactionCtxKey, tx), nil
}

func BeginNestedContextTransaction(ctx context.Context, db *sqlx.DB) (*NestedTransaction, error) {
	var err error

	if ntx := GetNestedContextTransaction(ctx); ntx != nil {
		ntx.Begin()
		return ntx, nil
	}

	ntx := &NestedTransaction{
		//tx:           db.BeginTx(ctx, nil),
		openedCount:  1,
		needRollback: false,
	}
	ntx.tx, err = db.Beginx()
	if err != nil {
		return nil, err
	}

	ntx.Context = context.WithValue(ctx, nestedTransactionCtxKey, ntx)

	return ntx, nil
}

type NestedTransaction struct {
	tx           *sqlx.Tx
	openedCount  int
	needRollback bool
	Context      context.Context
}

func (nt *NestedTransaction) GetTx() *sqlx.Tx {
	if nt == nil {
		return nil
	}
	return nt.tx
}

func (nt *NestedTransaction) Begin() {
	nt.openedCount++
	fmt.Println("Begin transaction")
}

func (nt *NestedTransaction) Commit() error {
	nt.openedCount--
	fmt.Println("Commit - ", nt.openedCount)
	if nt.openedCount > 0 {
		return nil
	} else if nt.openedCount < 0 {
		return fmt.Errorf("Opened transaction count < 0 ")
	}

	if nt.needRollback {
		err := nt.tx.Rollback()
		return fmt.Errorf("Transaction mark as rollback. Rollback executed: %s ", err)
	}

	return nt.tx.Commit()
}

func (nt *NestedTransaction) Rollback() error {
	nt.needRollback = true
	nt.openedCount--
	fmt.Println("Rollback - ", nt.openedCount)
	if nt.openedCount > 0 {
		return nil
	} else if nt.openedCount < 0 {
		return fmt.Errorf("Opened transaction count < 0 ")
	}

	return nt.tx.Rollback()
}
