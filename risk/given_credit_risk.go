package risk

import "strconv"

func CalculateCreditRisk(age, numberOfCreditCard int) string {
	sum := age + numberOfCreditCard
	mod := sum % 3
	if mod == 0 {
		return "LOW"
	}
	if mod == 1 {
		return "MEDIUM"
	}
	return "HIGH"
}

func CalculateAMLScore(jobIndustryCode string) (int, error) {
	jobIndustryPrefix := jobIndustryCode[0:1]
	jobIndustryValue, err := strconv.Atoi(jobIndustryPrefix)
	if err != nil {
		return 999999, err
	}

	return jobIndustryValue * 100, nil
}
