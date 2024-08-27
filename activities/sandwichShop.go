package activities

import (
	"context"
	"log"
)

func MakeSandwich(ctx context.Context) error {
	log.Printf("Sandwich is preparing...")
	return nil
}

func ServeSandwich(ctx context.Context, customerName string) error {
	log.Printf("%s has taken a sandwich after payment.", customerName)
	return nil
}

func BorrowSandwich(ctx context.Context, customerName string) error {
	log.Printf("%s borrowed a sandwich", customerName)
	return nil
}
