package main

import (
	"context"
	"github.com/go-chi/chi"
	"property_transactions/property_transactions"
)

func main() {

	//port:=emv
	//db env
	ctx := context.Background()
	c := property_transactions.Config{}
	s, err := property_transactions.New(ctx, c)
	if err != nil {
		panic(err)
	}
	_ = s
	r := chi.NewRouter()
	_ = r
}
