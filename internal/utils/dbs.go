package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// var xmlErrors = map[string]string{
// 	"A001": "Organisation ID is incorrect",
// 	"A002": "Maximum transaction transmission is exceeded",
// 	"A003": "Invalid Request",
// 	"A004": "Security credential is incorrect",
// 	"A005": "Transaction has timed out",
// 	"A006": "Gateway System Error",
// 	"A009": "Internal Server Error",
// 	"A010": "Security Check failed",
// 	"A011": "Invalid API Key",
// 	"A012": "User is not authorized to access this API",
// }

// func GetDummyXMLError(code string) error {
// 	if msg, ok := xmlErrors[code]; ok {
// 		return errors.New("XML Error [" + code + "]: " + msg)
// 	}
// 	return errors.New("XML Error: Unknown")
// }

// func CallDummyDBSAPI(req dto.StatementRequest) (interface{}, error) {
// 	r := rand.Intn(3)

// 	switch r {
// 	case 0:
// 		return getSuccessResponse(req), nil
// 	case 1:
// 		return map[string]interface{}{
// 			"header": map[string]string{
// 				"msgId":     req.Header.MsgId,
// 				"orgId":     req.Header.OrgId,
// 				"timeStamp": time.Now().Format(time.RFC3339),
// 				"ctry":      "SG",
// 			},
// 			"txnEnqResponse": map[string]string{
// 				"enqStatus":            "RJCT",
// 				"enqRejectCode":        "I301",
// 				"enqStatusDescription": "No available Statement",
// 			},
// 		}, nil
// 	default:
// 		return nil, GetDummyXMLError("A001") // simulate XML error
// 	}
// }

func ParseAPIResponse(status int, body []byte) (interface{}, string, error) {
	if status == 200 {
		var success map[string]interface{}
		_ = json.Unmarshal(body, &success)
		return success, "ACSP", nil
	}

	// if status is known gateway error
	if status >= 400 {
		var failure map[string]interface{}
		_ = json.Unmarshal(body, &failure)
		return failure, "RJCT", errors.New("rejected or gateway error")
	}

	return nil, "RJCT", errors.New("unknown error")
}

func CallExternalAPI(url string, payload interface{}) ([]byte, int, error) {
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, res.StatusCode, nil
}
