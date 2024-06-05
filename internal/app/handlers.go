package app

import (
	"encoding/json"
	"net/http"
	"via-varejo/helpers"
	"via-varejo/internal/domain"

	"github.com/gofrs/uuid"
)

func CreateNewPurchaseSummary(w http.ResponseWriter, r *http.Request) {
	var newPurchaseSummary domain.PurchaseSummary

	if erro := json.NewDecoder(r.Body).Decode(&newPurchaseSummary); erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusBadRequest)
		return
	}

	newPurchaseSummaryResponse, err := processPurchaseSummary(newPurchaseSummary)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newPurchaseSummaryResponse)
}

func processPurchaseSummary(newPurchaseSummary domain.PurchaseSummary) (domain.PurchaseSummaryResponse, error) {
	var response domain.PurchaseSummaryResponse
	id, _ := uuid.NewV4()

	if newPurchaseSummary.PaymentInfo.Installments > 6 {
		interestRate, err := helpers.GetInterestRates("https://api.bcb.gov.br/dados/serie/bcdata.sgs.4390/dados?formato=json")
		if err != nil {
			return response, err
		}

		totalAmount := helpers.GetInstallmentsValue(newPurchaseSummary.ProductInfo.Price, newPurchaseSummary.PaymentInfo.DownPaymentAmount, newPurchaseSummary.PaymentInfo.Installments, interestRate)
		purchaseSummary := createPurchaseSummary(newPurchaseSummary, id, totalAmount)

		err = insertPurchaseSummaryToRedis(id.String(), purchaseSummary)
		if err != nil {
			return response, err
		}

		response = createPurchaseSummaryResponse(newPurchaseSummary, totalAmount, interestRate)
	} else {
		purchaseSummary := createPurchaseSummary(newPurchaseSummary, id, newPurchaseSummary.ProductInfo.Price)

		err := insertPurchaseSummaryToRedis(id.String(), purchaseSummary)
		if err != nil {
			return response, err
		}

		response = createPurchaseSummaryResponse(newPurchaseSummary, newPurchaseSummary.ProductInfo.Price, 0)
	}

	return response, nil
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
