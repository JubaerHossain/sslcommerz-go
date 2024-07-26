package sslcommerz

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "sort"
    "strings"
)

// HashPassword hashes the store password using MD5
func HashPassword(password string) string {
    hash := md5.New()
    hash.Write([]byte(password))
    return hex.EncodeToString(hash.Sum(nil))
}

// VerifyHash verifies the hash signature of the response
func VerifyHash(storePass string, postData map[string]interface{}) (bool, error) {
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

    newData["store_passwd"] = HashPassword(storePass)
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

// GetImage retrieves the image URL from the response data
func GetImage(gw string, source map[string]interface{}) (string, error) {
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
