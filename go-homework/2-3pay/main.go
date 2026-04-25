package main

func main() {
	pay := NewPayService("wechat")
	err := pay.Pay(100.0, "Alice", "Bob")
	if err != nil {
		panic(err)
	}
	pay = NewPayService("credit_card")
	err = pay.Pay(100.0, "Alice", "Bob")
	if err != nil {
		panic(err)
	}
	pay = NewPayService("alipay")
	err = pay.Pay(100.0, "Alice", "Bob")
	if err != nil {
		panic(err)
	}
}
