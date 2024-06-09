package controller

import (
	"app/config"
	"app/dto/request"
	"app/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type paymentController struct {
	paymentUtils utils.PaymentUtils
}

type PaymentController interface {
	CreateBillOnline(w http.ResponseWriter, r *http.Request)
}

func (c *paymentController) CreateBillOnline(w http.ResponseWriter, r *http.Request) {
	var payload request.CreateBillOnlineRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ipAddr := strings.Join([]string{
		r.Header.Get("x-forwarded-for"),
		r.RemoteAddr,
	}, ",")

	tmnCode := config.GetVnpTmncode()
	secretKey := config.GetVnpHashsecret()
	vnpUrl := config.GetVnpUrl()
	returnUrl := config.GetVnpReturnUrl()

	date := time.Now()
	createDate := date.Format("20060102150405")
	bankCode := payload.BankCode
	expireDate := date.Add(15 * time.Minute).Format("20060102150405")
	orderId, errUUID := uuid.NewV6()

	if errUUID != nil {
		internalServerError(w, r, errUUID)
		return
	}

	amount := payload.Amount * 100
	orderInfo := payload.OrderDescription
	orderType := payload.OrderType
	locale := "vn"
	currCode := "VND"

	vnpParams := map[string]interface{}{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    tmnCode,
		"vnp_Locale":     locale,
		"vnp_CurrCode":   currCode,
		"vnp_TxnRef":     orderId.String(),
		"vnp_OrderInfo":  orderInfo,
		"vnp_OrderType":  orderType,
		"vnp_Amount":     amount,
		"vnp_ReturnUrl":  returnUrl,
		"vnp_IpAddr":     ipAddr,
		"vnp_CreateDate": createDate,
		"vnp_ExpireDate": expireDate,
	}

	if bankCode != "" {
		vnpParams["vnp_BankCode"] = bankCode
	}

	vnpParams = c.paymentUtils.SortMap(vnpParams)

	signData := url.Values{}
	for key, value := range vnpParams {
		convert := fmt.Sprint(value)
		signData.Add(key, convert)
	}

	signature := c.paymentUtils.GenerateSignature(signData.Encode(), secretKey)
	signData.Add("vnp_SecureHash", signature)

	redirectURL := vnpUrl + "?" + signData.Encode()

	res := Response{
		Data:    redirectURL,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewPaymentController() PaymentController {
	return &paymentController{
		paymentUtils: utils.NewPaymentUtils(),
	}
}
