package hooks

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterProductsHooks(app core.App) {
	//Create
	app.OnRecordCreate("products").BindFunc(func(e *core.RecordEvent) error {
		e.Record.Set("is_active", false)
		e.Record.Set("stock", 0)

		return e.Next()
	})
	
	//Delete
	app.OnRecordDelete("products").BindFunc(func(e *core.RecordEvent) error {
		orderItems, err := e.App.FindAllRecords("order_items",
			dbx.HashExp{"products": e.Record.Id},
		)
		if err != nil {
			return err
		}
		if len(orderItems) > 0 {
			return errors.New("нельзя удалить товар - он уже заказан")
		}

		return e.Next()
	})
}
