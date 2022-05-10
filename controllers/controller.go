package controllers

import (
	"encoding/json"
	"fmt"
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

type JSONResponse struct {
	Status string `json:"status"`
}

func ProcessData(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		panic("not implemented")
	default:
		log.Println("error no 404")
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, "not found")
	}

}
