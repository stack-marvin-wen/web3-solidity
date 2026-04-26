package main

import "fmt"

type PayServiceInterface interface {
	Pay(amount float64, from string, to string) error
}

func NewPayService(payType string) PayServiceInterface {
	switch payType {
	case "credit_card":
		return &CraditCardPayService{}
	case "wechat":
		return &WechatPayService{}
	case "alipay":
		return &AliPayService{}
	default:
		return nil
	}
}

type CraditCardPayService struct {
}

func (c *CraditCardPayService) Pay(amount float64, from string, to string) error {
	fmt.Println("使用信用卡支付：", amount, "从", from, "到", to)
	return nil
}

type WechatPayService struct {
}

func (w *WechatPayService) Pay(amount float64, from string, to string) error {
	fmt.Println("使用微信支付：", amount, "从", from, "到", to)
	return nil
}

type AliPayService struct {
}

func (a *AliPayService) Pay(amount float64, from string, to string) error {
	fmt.Println("使用支付宝支付：", amount, "从", from, "到", to)
	return nil
}
