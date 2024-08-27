package activities

import (
	"context"
)

func SendEmail(ctx context.Context, email string) error {
	/*
		// random number generation for panic simulation
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			log.Println("Failed to generate random number:", err)
			return err
		}

		// panic if n is 0
		if n.Int64() == 0 {
			panic("Simulated panic")
		}
	*/

	// log email
	// log.Printf("%s", email)
	return nil
}
