package sslcommerz

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/JubaerHossain/sslcommerz-go/client"
	"github.com/JubaerHossain/sslcommerz-go/config"
)

type SSLCommerz struct {
	storeID           string
	storePass         string
	sslcSubmitURL     string
	sslcValidationURL string
	sslcMode          string
	sslcData          map[string]interface{}
}

func NewSSLCommerz() *SSLCommerz {
	sslc := &SSLCommerz{
		storeID:   config.StoreID,
		storePass: config.StorePass,
	}

	if config.IS_SANDBOX == "false" {
		sslc.sslcMode = "securepay"
		sslc.sslcSubmitURL = "https://securepay.sslcommerz.com/gwprocess/v3/api.php"
		sslc.sslcValidationURL = "https://securepay.sslcommerz.com/validator/api/validationserverAPI.php"
	} else {
		sslc.sslcMode = "sandbox"
		sslc.sslcSubmitURL = "https://sandbox.sslcommerz.com/gwprocess/v4/api.php"
		sslc.sslcValidationURL = "https://sandbox.sslcommerz.com/validator/api/validationserverAPI.php"
	}

	return sslc
}

func (s *SSLCommerz) InitiatePayment(postData map[string]interface{}) (map[string]interface{}, error) {
	postData["store_id"] = s.storeID
	postData["store_passwd"] = s.storePass

	client := client.NewClient()
	response, err := client.MakeRequest("POST", s.sslcSubmitURL, postData)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}

	if status, ok := result["status"].(string); ok && status == "SUCCESS" {
		return result, nil
	} else if status, ok := result["status"].(string); ok && status == "FAILED" {
		return result, nil
	} else if status, ok := result["status"].(string); ok && status == "CANCELLED" {
		return result, nil
	}

	return nil, errors.New("failed to initiate payment")
}

func (s *SSLCommerz) SSLCOMMERZ_hash_varify(storePass string, postData map[string]interface{}) (bool, error) {
	verifySign, ok1 := postData["verify_sign"].(string)
	verifyKey, ok2 := postData["verify_key"].(string)
	if !ok1 || !ok2 {
		return false, errors.New("required data missing: verify_key, verify_sign")
	}

	preDefineKey := strings.Split(verifyKey, ",")
	newData := make(map[string]string)

	for _, key := range preDefineKey {
		if val, exists := postData[key]; exists {
			newData[key] = fmt.Sprintf("%v", val)
		}
	}

	newData["store_passwd"] = fmt.Sprintf("%x", md5.Sum([]byte(storePass)))

	keys := make([]string, 0, len(newData))
	for k := range newData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var hashString string
	for _, k := range keys {
		hashString += k + "=" + newData[k] + "&"
	}
	hashString = strings.TrimRight(hashString, "&")

	hash := md5.Sum([]byte(hashString))
	hashHex := hex.EncodeToString(hash[:])

	if hashHex == verifySign {
		return true, nil
	}

	return false, errors.New("verification signature not matched")
}

func (s *SSLCommerz) ValidateTransaction(tranID, amount, currency string, postData map[string]interface{}) (bool, error) {
	if tranID == "" || amount == "" || currency == "" {
		return false, errors.New("invalid transaction data")
	}

	postData["store_id"] = s.storeID
	postData["store_passwd"] = s.storePass

	valid, err := s.SSLCOMMERZ_hash_varify(s.storePass, postData)
	if err != nil {
		return false, err
	}

	if !valid {
		return false, errors.New("hash validation failed")
	}

	valID := url.QueryEscape(postData["val_id"].(string))
	storeID := url.QueryEscape(s.storeID)
	storePasswd := url.QueryEscape(s.storePass)

	validationURL := fmt.Sprintf("%s?val_id=%s&store_id=%s&store_passwd=%s&v=1&format=json", s.sslcValidationURL, valID, storeID, storePasswd)

	resp, err := http.Get(validationURL)
	if err != nil {
		return false, errors.New("failed to connect with SSLCOMMERZ")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("failed to connect with SSLCOMMERZ")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, err
	}

	status, ok := result["status"].(string)
	if !ok || (status != "VALID" && status != "VALIDATED") {
		return false, errors.New("transaction validation failed")
	}

	amountFloat, err := strconv.ParseFloat(fmt.Sprintf("%v", result["amount"]), 64)
	if err != nil {
		return false, err
	}

	currencyAmount, err := strconv.ParseFloat(fmt.Sprintf("%v", result["currency_amount"]), 64)
	if err != nil {
		return false, err
	}

	finalAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false, err
	}

	if currency == "BDT" {
		if tranID == postData["tran_id"].(string) && ((amountFloat-finalAmount) < 1) && currency == "BDT" {
			return true, nil
		}
	} else {
		if tranID == postData["tran_id"].(string) && ((currencyAmount-finalAmount) < 1) && currency == fmt.Sprintf("%v", result["currency_type"]) {
			return true, nil
		}
	}

	return false, errors.New("data has been tampered with")
}

func (s *SSLCommerz) HashVerify(storePass string, postData map[string]interface{}) (bool, error) {
	if postData == nil {
		return false, errors.New("post data is nil")
	}

	verifySign, ok := postData["verify_sign"].(string)
	if !ok || verifySign == "" {
		return false, errors.New("verify_sign missing or invalid")
	}

	verifyKey, ok := postData["verify_key"].(string)
	if !ok || verifyKey == "" {
		return false, errors.New("verify_key missing or invalid")
	}

	keys := strings.Split(verifyKey, ",")
	newData := make(map[string]string)

	for _, key := range keys {
		if val, ok := postData[key].(string); ok {
			newData[key] = val
		}
	}

	newData["store_passwd"] = s.hashPassword(storePass)
	var sortedKeys []string
	for key := range newData {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var hashString string
	for _, key := range sortedKeys {
		hashString += key + "=" + newData[key] + "&"
	}
	hashString = strings.TrimRight(hashString, "&")

	hash := md5.New()
	hash.Write([]byte(hashString))
	generatedHash := hex.EncodeToString(hash.Sum(nil))

	if generatedHash == verifySign {
		return true, nil
	}

	return false, errors.New("hash verification failed")
}

func (s *SSLCommerz) hashPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *SSLCommerz) GetImage(gw string, source map[string]interface{}) (string, error) {
	if source == nil {
		return "", errors.New("source data is nil")
	}

	desc, ok := source["desc"].([]interface{})
	if !ok || len(desc) == 0 {
		return "", errors.New("desc field missing or invalid")
	}

	for _, item := range desc {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if itemMap["gw"] == gw {
			logo, ok := itemMap["logo"].(string)
			if ok {
				return strings.Replace(logo, "/gw/", "/gw1/", 1), nil
			}
		}
	}

	return "", errors.New("logo not found")
}
