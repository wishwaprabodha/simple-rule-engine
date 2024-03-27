package rules

import (
	"github.com/honestbank/tech-assignment-backend-engineer/config"
	"github.com/honestbank/tech-assignment-backend-engineer/risk"
	"log"
	"os"
	"strconv"
	"strings"
)

func validateCondition(configName string, condition func(int) bool) bool {
	configValue, exist := os.LookupEnv(configName)
	if !exist {
		log.Fatalf("no %s config specified", configName)
	}
	configInt, err := strconv.Atoi(configValue)
	if err != nil {
		log.Fatalf("error converting %s config to number: %v", configName, err)
	}
	return condition(configInt)
}

func validateConditionList(configName string, condition func([]string) bool) bool {
	configValue, exist := os.LookupEnv(configName)
	if !exist {
		log.Fatalf("no %s config specified", configName)
	}
	configList := strings.Split(configValue, ",")
	return condition(configList)
}

func ValidateIncome(income int) bool {
	return validateCondition("MINIMUM_INCOME", func(configIncome int) bool {
		return income > configIncome
	})
}

func ValidateAge(age int) bool {
	return validateCondition("MINIMUM_INCOME", func(configAge int) bool {
		return age >= configAge
	})
}

func ValidateCreditCardCount(creditCardCount int, age int) bool {
	return validateCondition("MAX_CREDIT_CARDS", func(maxCreditCardCount int) bool {
		return creditCardCount <= maxCreditCardCount && risk.CalculateCreditRisk(age, creditCardCount) == os.Getenv("ACCEPTED_CREDIT_RISK")
	})
}

func ValidateAreaCode(phone string) bool {
	return validateConditionList("ACCEPTED_AREA_CODES", func(areaCodeList []string) bool {
		for _, areaCode := range areaCodeList {
			if strings.HasPrefix(string(phone[0]), areaCode) {
				return true
			}
		}
		return false
	})
}

func OverrideValidation(phone string) bool {
	numberList := config.ConfigInstance.Config.PhoneNumbers
	for _, number := range numberList {
		if phone == number {
			return true
		}
	}
	return false
}
