package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"sort"
)

type paymentUtils struct{}

type PaymentUtils interface {
	SortMap(inputMap map[string]interface{}) map[string]interface{}
	GenerateSignature(data, secretKey string) string
}

func (u *paymentUtils) SortMap(inputMap map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	sortedMap := make(map[string]interface{})
	for _, key := range keys {
		sortedMap[key] = inputMap[key]
	}
	return sortedMap
}

func (u *paymentUtils) GenerateSignature(data, secretKey string) string {
	hash := hmac.New(sha512.New, []byte(secretKey))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func NewPaymentUtils() PaymentUtils {
	return &paymentUtils{}
}
