package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"via-varejo/helpers"
	"via-varejo/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func CreateNewPurchaseSummary(w http.ResponseWriter, r *http.Request) {
	var newPurchaseSummary domain.PurchaseSummary

	if erro := json.NewDecoder(r.Body).Decode(&newPurchaseSummary); erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusBadRequest)
		return
	}

	newPurchaseSummaryResponse, purchaseID, err := processPurchaseSummary(newPurchaseSummary)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newPurchaseSummaryResponse)
	http.Redirect(w, r, "/purchases/"+purchaseID.String(), http.StatusSeeOther)
}

func processPurchaseSummary(newPurchaseSummary domain.PurchaseSummary) (domain.PurchaseSummaryResponse, uuid.UUID, error) {
	var response domain.PurchaseSummaryResponse
	id, _ := uuid.NewV4()

	if newPurchaseSummary.PaymentInfo.Installments > 6 {
		interestRate, err := helpers.GetInterestRates("https://api.bcb.gov.br/dados/serie/bcdata.sgs.4390/dados?formato=json")
		if err != nil {
			return response, id, err
		}

		totalAmount := helpers.GetInstallmentsValue(newPurchaseSummary.ProductInfo.Price, newPurchaseSummary.PaymentInfo.DownPaymentAmount, newPurchaseSummary.PaymentInfo.Installments, interestRate)
		purchaseSummary := createPurchaseSummary(newPurchaseSummary, id, totalAmount)

		err = insertPurchaseSummaryToRedis(id.String(), purchaseSummary)
		if err != nil {
			return response, id, err
		}

		response = createPurchaseSummaryResponse(newPurchaseSummary, totalAmount, interestRate)
	} else {
		purchaseSummary := createPurchaseSummary(newPurchaseSummary, id, newPurchaseSummary.ProductInfo.Price)

		err := insertPurchaseSummaryToRedis(id.String(), purchaseSummary)
		if err != nil {
			return response, id, nil
		}

		response = createPurchaseSummaryResponse(newPurchaseSummary, newPurchaseSummary.ProductInfo.Price, 0)
	}
	return response, id, nil
}

func createPurchaseSummary(newPurchaseSummary domain.PurchaseSummary, id uuid.UUID, totalAmount float64) domain.PurchaseSummary {
	return domain.PurchaseSummary{
		ProductInfo: newPurchaseSummary.ProductInfo,
		PaymentInfo: domain.PaymentMethod{
			DownPaymentAmount: newPurchaseSummary.PaymentInfo.DownPaymentAmount,
			Installments:      newPurchaseSummary.PaymentInfo.Installments,
			TotalAmount:       totalAmount,
		},
		ID: id,
	}
}

func insertPurchaseSummaryToRedis(key string, purchaseSummary domain.PurchaseSummary) error {
	return domain.InsertRedisPurchaseSummary(key, purchaseSummary)
}

func createPurchaseSummaryResponse(newPurchaseSummary domain.PurchaseSummary, totalAmount float64, interestRate float64) domain.PurchaseSummaryResponse {
	return domain.PurchaseSummaryResponse{
		Installments: newPurchaseSummary.PaymentInfo.Installments,
		TotalAmount:  totalAmount,
		InterestRate: interestRate,
	}
}

func GetRedisPurchaseSummary(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	purchaseSummary, err := domain.FindByIDRedisPurchaseSumary(id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, purchaseSummary)

}

func CreateConfirmNewPurchase(w http.ResponseWriter, r *http.Request) {
	var purchase domain.Purchase
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.ErrorJSON(w, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	purchaseSummary, err := domain.FindByIDRedisPurchaseSumary(id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println(purchaseSummary)

	purchase = domain.Purchase{
		ID: purchaseSummary.ID,
		ProductInfo: domain.Product{
			Name:        purchaseSummary.ProductInfo.Name,
			Price:       purchaseSummary.ProductInfo.Price,
			ProductCode: purchaseSummary.ProductInfo.ProductCode,
		},
		PaymentInfo: domain.PaymentMethod{
			DownPaymentAmount: purchaseSummary.PaymentInfo.DownPaymentAmount,
			Installments:      purchaseSummary.PaymentInfo.Installments,
		},
		PurchaseDate: time.Now(),
	}

	fmt.Println(purchase, "purchase")

	err = domain.InsertNewPurchaseToPostgres(purchase)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, purchaseSummary)

}
