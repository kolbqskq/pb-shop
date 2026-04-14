package routes

import (
	"net/http"
	"pb-shop/internal/services"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func Register(se *core.ServeEvent) {
	g := se.Router.Group("/api/myshop")
	g.Bind(apis.RequireAuth())

	g.POST("/checkout", handleCheckout)
}

func handleCheckout(e *core.RequestEvent) error {
	orderService := services.NewOrderService(e.App)

	order, err := orderService.CreateFromCart(e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, order)
}
