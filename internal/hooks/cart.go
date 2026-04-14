package hooks

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterCartHooks(app core.App) {
	//Create
	app.OnRecordCreateRequest("cart_items").BindFunc(func(e *core.RecordRequestEvent) error {

		e.Record.Set("user", e.Auth.Id)

		product, err := e.App.FindRecordById("products", e.Record.GetString("product"))
		if err != nil {
			return errors.New("товар не найден")
		}

		if !product.GetBool("is_active") {
			return errors.New("товар недоступен")
		}

		quantity := e.Record.GetInt("quantity")
		if quantity <= 0 {
			return errors.New("кол-во должно быть больше нуля")
		}
		if product.GetInt("stock") < quantity {
			return errors.New("недостаточно товара на складе")
		}

		existing, err := e.App.FindAllRecords("cart_items",
			dbx.HashExp{
				"user":    e.Auth.Id,
				"product": e.Record.GetString("product"),
			},
		)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			item := existing[0]
			item.Set("quantity", item.GetInt("quantity")+quantity)

			if err := e.App.Save(item); err != nil {
				return err
			}
			return e.JSON(200, item)
		}
		return e.Next()
	})
	//Update
	app.OnRecordUpdateRequest("cart_items").BindFunc(func(e *core.RecordRequestEvent) error {

		quantity := e.Record.GetInt("quantity")

		if quantity < 0 {
			return errors.New("кол-во должно быть больше нуля")
		}

		product, err := e.App.FindRecordById("products", e.Record.GetString("product"))
		if err != nil {
			return err
		}

		if product.GetInt("stock") < quantity {
			return errors.New("недостаточно товара на складе")
		}

		return e.Next()
	})
}
