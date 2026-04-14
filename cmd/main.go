package main

import (
	"log"
	"pb-shop/internal/hooks"
	"pb-shop/internal/routes"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		hooks.RegisterProductsHooks(app)
		hooks.RegisterCartHooks(app)
		hooks.RegisterOrderHooks(app)

		routes.Register(e)

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
