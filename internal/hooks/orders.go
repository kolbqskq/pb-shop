package hooks

import (
	"errors"
	"pb-shop/internal/services"

	"github.com/pocketbase/pocketbase/core"
)

func RegisterOrderHooks(app core.App) {
	//Update
	app.OnRecordUpdateRequest("orders").BindFunc(func(e *core.RecordRequestEvent) error {
		oldStatus := e.Record.Original().GetString("status")
		newStatus := e.Record.GetString("status")

		if newStatus == "cancelled" && oldStatus == "pending" {
			orderService := services.NewOrderService(e.App)
			return orderService.Cancel(e.Record)
		}

		if oldStatus != newStatus {
			return errors.New("недопустимое изменение статуса")
		}

		return e.Next()
	})

	//Delete
}
