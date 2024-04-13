package main

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logrus.Fatalf("failed to init app: %s", err.Error())
	}

	a.Run()
}
