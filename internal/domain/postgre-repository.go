// domain/postgre-repository.go

package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func InsertNewPurchaseToPostgres(purchase Purchase) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection pool is nil")
	}

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	query := `
        INSERT INTO purchases (id, product_name, product_price, product_code, down_payment_amount, installments, total_amount, purchase_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, query,
		purchase.ID,
		purchase.ProductInfo.Name,
		purchase.ProductInfo.Price,
		purchase.ProductInfo.ProductCode,
		purchase.PaymentInfo.DownPaymentAmount,
		purchase.PaymentInfo.Installments,
		purchase.ProductInfo.Price, // Assuming total_amount is the same as product price for this example
		purchase.PurchaseDate,
	)

	if err != nil {
		return err
	}

	return nil
}
