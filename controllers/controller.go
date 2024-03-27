package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/honestbank/tech-assignment-backend-engineer/rules"
	"log"
	"net/http"
)

type RecordData struct {
	Income              int    `json:"income"`
	NumberOfCreditCards int    `json:"number_of_credit_cards"`
	Age                 int    `json:"age"`
	PoliticallyExposed  bool   `json:"politically_exposed"`
	JobIndustryCode     string `json:"job_industry_code"`
	PhoneNumber         string `json:"phone_number"`
}

var APPROVED = "approved"
var DECLINED = "declined"

type JSONResponse struct {
	Status string `json:"status"`
}

func sendResponse(resp http.ResponseWriter, validationStatus bool) {
	response := JSONResponse{}

	if validationStatus {
		response.Status = APPROVED
	} else {
		response.Status = DECLINED
	}

	data, _ := json.Marshal(&response)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(data)
}

func ProcessData(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		isRequestValidated, inputParams, err := validateInput(req)
		if err != nil || !isRequestValidated {
			log.Fatal("Request Body Validation Failed")
		}
		validationStatus := evaluateRules(inputParams.Income, inputParams.NumberOfCreditCards, inputParams.Age, inputParams.PoliticallyExposed, inputParams.PhoneNumber)
		sendResponse(resp, validationStatus)
	default:
		log.Println("error no 404")
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, "not found")
	}

}

func validateInput(req *http.Request) (bool, *RecordData, error) {
	var inputParams RecordData
	if err := json.NewDecoder(req.Body).Decode(&inputParams); err != nil {
		return false, nil, err
	}
	return true, &inputParams, nil
}

func evaluateRules(income int, ownedCreditCards int, age int, pep bool, phone string) bool {
	// Custom Rules Order can be swapped respective to information gain
	if rules.OverrideValidation(phone) {
		return true
	} else if pep {
		return false
	} else if !rules.ValidateIncome(income) {
		return false
	} else if !rules.ValidateAge(age) {
		return false
	} else if !rules.ValidateCreditCardCount(ownedCreditCards, age) {
		return false
	} else if !rules.ValidateAreaCode(phone) {
		return false
	}
	return true
}
