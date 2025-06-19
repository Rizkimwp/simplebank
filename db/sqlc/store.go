package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	Queries *Queries
	db      *sql.DB
}

// NewStore creates a new Store instance with the given database connection
func NewStore(db *sql.DB) *Store {
	store := &Store{
		Queries: New(db),
		db:      db,
	}
	return store
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64         `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTxParams contains the input parameters for the transfer transaction
func (store *Store) transferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult
	
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID:  sql.NullInt64{Int64: arg.FromAccountID, Valid: true} , 
			ToAccountID:    sql.NullInt64{Int64: arg.ToAccountID, Valid: true}, 
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create From Entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:  sql.NullInt64{Int64:arg.FromAccountID, Valid: true},
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create To Entry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Get account 1 for update")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "Update Account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Get Account 2 For Update")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "Update Account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}
		
		return nil
	})

	return result, err
}
