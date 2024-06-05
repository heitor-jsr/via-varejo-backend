package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}

type DataPoint struct {
	Data  string `json:"data"`
	Valor string `json:"valor"`
}

func GetInterestRates(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var dataPoints []DataPoint
	err = json.NewDecoder(resp.Body).Decode(&dataPoints)
	if err != nil {
		fmt.Println("erro aqui")
		return 0, err
	}

	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	var lastRate float64

	for _, dataPoint := range dataPoints {
		date, err := time.Parse("02/01/2006", dataPoint.Data)
		if err != nil {
			return 0, err
		}

		if date.Month() == currentMonth && date.Year() == currentYear {
			rate, err := strconv.ParseFloat(dataPoint.Valor, 64)
			if err != nil {
				return 0, err
			}
			return rate, nil
		} else if date.Year() == currentYear {
			lastRate, err = strconv.ParseFloat(dataPoint.Valor, 64)
			if err != nil {
				return 0, err
			}
		}
	}

	if lastRate == 0 {
		return 0, fmt.Errorf("no interest rate available for the current year")
	}

	return lastRate, nil
}

func GetInstallmentsValue(price float64, downpayment float64, installments int, interestRate float64) float64 {
	principal := price - downpayment

	installmentValue := (principal * interestRate) / (1 - (1 / (math.Pow(1+interestRate, float64(installments)))))
	totalAmountWithInterest := installmentValue*float64(installments) + downpayment
	totalAmountWithInterest, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmountWithInterest), 64)

	return totalAmountWithInterest
}
