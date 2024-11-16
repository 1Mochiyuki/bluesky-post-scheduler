package app

import (
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp() App {
	return App{
		ctx: context.Background(),
	}
}
