package api

import (
	"app/config"
	"app/grpc/proto"
	"app/utils"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type billGrpc struct {
	db           *gorm.DB
	paymentUtils utils.PaymentUtils
	proto.UnsafeBillServiceServer
}

func (g *billGrpc) CreateBill(ctx context.Context, req *proto.CreateBillReq) (*proto.CreateBillRes, error) {
	tmnCode := config.GetVnpTmncode()
	secretKey := config.GetVnpHashsecret()
	vnpUrl := config.GetVnpUrl()
	returnUrl := config.GetVnpReturnUrl()

	date := time.Now()
	createDate := date.Format("20060102150405")
	expireDate := time.Unix(req.ExpireDate, 0).Format("20060102150405")
	orderId, errUUID := uuid.NewV6()

	if errUUID != nil {
		return nil, errUUID
	}

	amount := int(req.Amount) * 100
	orderInfo := req.OrderDescription
	orderType := req.OrderType
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
		"vnp_IpAddr":     req.IpAddr,
		"vnp_CreateDate": createDate,
		"vnp_ExpireDate": expireDate,
	}

	vnpParams = g.paymentUtils.SortMap(vnpParams)

	signData := url.Values{}
	for key, value := range vnpParams {
		convert := fmt.Sprint(value)
		signData.Add(key, convert)
	}

	signature := g.paymentUtils.GenerateSignature(signData.Encode(), secretKey)
	signData.Add("vnp_SecureHash", signature)

	redirectURL := vnpUrl + "?" + signData.Encode()

	res := proto.CreateBillRes{
		UuidBill: orderId.String(),
		HrefVnp:  redirectURL,
	}

	return &res, nil
}

func NewBillGrpc() proto.BillServiceServer {
	return &billGrpc{
		db:           config.GetDB(),
		paymentUtils: utils.NewPaymentUtils(),
	}
}
