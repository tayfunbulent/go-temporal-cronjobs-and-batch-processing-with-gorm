package activities

import (
	"context"
	"crypto/rand"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

var attempt int

func PrepareCoffee(ctx context.Context) error {
	attempt++

	log.Println("----------------")
	log.Println(attempt, "= attempt count")
	log.Println("----------------")
	if attempt%3 != 0 {
		log.Println("panic")
		log.Println("----------------")
		panic("Simulated panic")
	}

	log.Println("----------------")
	log.Println("no panic")
	log.Println("Coffee is preparing...")
	return nil
}

func GiveCoffee(ctx context.Context, customerName string) error {
	log.Printf("%s adlı müşteriye kahve teslim edildi.", customerName)
	return nil
}

func WriteAsDept(ctx context.Context, customerName string) error {
	log.Printf("%s borrowed a coffee", customerName)
	return nil
}

func LongRunningActivity(ctx context.Context) error {
	// Bu aktivite 10 saniye bekleyecek, ancak zaman aşımı süresi sadece 3 saniye
	time.Sleep(10 * time.Second)
	return nil
}

func SendEmail(ctx context.Context, email string) error {
	n, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		log.Println("Failed to generate random number:", err)
		return err
	}

	if n.Int64() == 0 {
		log.Println("panic")
		log.Println("----------------")
		panic("Simulated panic")
	}

	log.Println("******************")
	log.Println("no panic")
	log.Printf("%s", email)
	log.Println("******************")
	return nil
}

func SendCurlRequest(ctx context.Context) error {
	// HTTP GET isteği gönder
	resp, err := http.Get("http://localhost:3310/start-workflow-email")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(resp.StatusCode, "Request failed")
	}

	return nil
}
