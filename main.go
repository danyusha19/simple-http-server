package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var money = 1000
var bank = 0

var mtx = sync.Mutex{}

func payHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "fail to read HTTP body: " + err.Error()
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
		return
	}

	httpRequestBodyString := string(httpRequestBody)

	paymentAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		fmt.Println("fail to convert HTTP body to integer:", err)
		return
	}

	mtx.Lock()
	if money-paymentAmount >= 0 {
		money -= paymentAmount
		msg := "Оплата прошла успешно! Баланс: " + strconv.Itoa(money)
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
	} else {
		msg := "Не хватает " + strconv.Itoa(paymentAmount-money) + " USD для оплаты!"
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
	}
	mtx.Unlock()
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "fail to read HTTP body:" + err.Error()
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
		return
	}
	httpRequestBodyString := string(httpRequestBody)

	saveAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {
		msg := "fail to convert HTTP body to integer: " + err.Error()
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
		return
	}

	mtx.Lock()
	if money-saveAmount >= 0 {
		money -= saveAmount
		bank += saveAmount

		msg := "Успешно положили в копилку! Баланс кошелька:" + strconv.Itoa(money) + "\nБаланс копилки:" + strconv.Itoa(bank)
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
			return
		}
	} else {
		msg := "Не хватает" + strconv.Itoa(saveAmount-money) + "USD для пополнения копилки!"
		fmt.Println(msg)

		_, err = w.Write([]byte(msg))
		if err != nil {
			fmt.Println("fail to write HTTP response")
		}
	}
	mtx.Unlock()
}

func main() {
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("HTTP server error:", err)
	}
}
