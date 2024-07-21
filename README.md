# Simple Rule Engine with Golang
This project implements a simple rule engine in Golang to evaluate applications based on a set of predefined rules. The engine is triggered by a POST request, processes the application data, and returns a JSON response indicating whether the application is "approved" or "declined".
## Features
* **Rule-Based Evaluation:**  Evaluates applications against a set of configurable rules.
* **HTTP Endpoint:**  Exposes a POST endpoint to receive application data in JSON format.
* **JSON Response:**  Returns a JSON response with the application status ("approved" or "declined").
* **Error Handling:**  Handles errors gracefully without interrupting the process.
* **External Data Integration:**  Retrieves credit risk scores from an external risk module.
* **Runtime Updatable Pre-approved List:**  Allows updating the list of pre-approved phone numbers without restarting the application.

## Specifications

### Request Body
| Field                  | Type    | Description                                                              |
|------------------------|---------|--------------------------------------------------------------------------|
|
 `income`               | number  | Applicant's income                                                        |
| `number_of_credit_cards` | number  | Number of credit cards held by the applicant                             |
| `age`                  | number  | Applicant's age                                                           |
| `politically_exposed`   | bool    | Indicates if the applicant is politically exposed (PEP)                 |
| `job_industry_code`    | string | Code representing the applicant's job industry                           |
| `phone_number`          | string | Applicant's phone number                                                  |


**Example:**

json { "income": 82428, "number_of_credit_cards": 3, "age": 9, "politically_exposed": true, "job_industry_code": "2-930 - Exterior Plants", "phone_number": "486-356-0375" }

### Response Body
| Field   | Type   | Description                               |
|---------|--------|-------------------------------------------|
| `status` | string | Application status: "approved" or "declined" |

##### Example

###### Approved:

```json
{
  "status": "approved"
}
```

###### Declined:

```json
{
  "status": "declined"
}
```

#### Decision Rules

The application is approved if it evaluates as `true` on the following rules:

1. The applicant must earn more than 100000.
1. The applicant must be at least 18 years old.
1. The applicant must not hold more than 3 credit cards and their `credit_risk_score` must be `LOW`.
1. The applicant must not be involved in any political activities (must not be a Politically Exposed Person or PEP).
1. The applicant's phone number must be in an area that is allowed to apply for this product. The area code is denoted by first digit of phone number. The allowed area codes are `0`, `2`, `5`, and `8`.
1. A pre-approved list of phone numbers should cause the application to be automatically approved without evaluation of the above rules. This list must be able to be updated at runtime without needing to restart the process.

#### External Data Sources

Values for the `credit_risk_score` field can be retrieved by calling the existing functions in the provided `risk` module.