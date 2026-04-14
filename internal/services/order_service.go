package services

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type OrderService struct {
	app core.App
}

func NewOrderService(app core.App) *OrderService {
	return &OrderService{
		app: app,
	}
}

func (s *OrderService) CreateFromCart(userId string) (*core.Record, error) {
	cartItems, err := s.app.FindAllRecords("cart_items",
		dbx.HashExp{"user": userId},
	)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("корзина пуста")
	}

	orderCollection, err := s.app.FindCollectionByNameOrId("orders")
	if err != nil {
		return nil, err
	}

	var resultOrder *core.Record

	err = s.app.RunInTransaction(func(txApp core.App) error {
		var total float64

		type cartLine struct {
			product  *core.Record
			quantity int
			price    float64
		}

		lines := make([]cartLine, 0, len(cartItems))

		for _, cartItem := range cartItems {
			product, err := txApp.FindRecordById("products",
				cartItem.GetString("product"),
			)
			if err != nil {
				return err
			}

			quantity := cartItem.GetInt("quantity")
			price := product.GetFloat("price")

			if product.GetInt("stock") < quantity {
				return errors.New("товар <" + product.GetString("name") + "> закончился на складе")
			}

			total += price * float64(quantity)
			lines = append(lines, cartLine{product, quantity, price})
		}

		order := core.NewRecord(orderCollection)
		order.Set("user", userId)
		order.Set("status", "pending")
		order.Set("total", total)

		if err := txApp.Save(order); err != nil {
			return err
		}

		for _, line := range lines {
			itemCollection, err := txApp.FindCollectionByNameOrId("order_items")
			if err != nil {
				return err
			}

			orderItem := core.NewRecord(itemCollection)
			orderItem.Set("order", order.Id)
			orderItem.Set("product", line.product.Id)
			orderItem.Set("quantity", line.quantity)
			orderItem.Set("price", line.price)

			if err := txApp.Save(orderItem); err != nil {
				return err
			}

			line.product.Set("stock", line.product.GetInt("stock")-line.quantity)
			if err := txApp.Save(line.product); err != nil {
				return err
			}
		}

		for _, cartItem := range cartItems {
			if err := txApp.Delete(cartItem); err != nil {
				return err
			}
		}

		resultOrder = order
		return nil
	})

	return resultOrder, err
}

func (s *OrderService) Cancel(order *core.Record) error {
	return s.app.RunInTransaction(func(txApp core.App) error {
		orderItems, err := txApp.FindAllRecords("order_items",
			dbx.HashExp{"order": order.Id},
		)
		if err != nil {
			return err
		}

		for _, item := range orderItems {
			product, err := txApp.FindRecordById("products",
				item.GetString("product"),
			)
			if err != nil {
				return err
			}

			product.Set("stock",
				product.GetInt("stock")+item.GetInt("quantity"),
			)

			if err := txApp.Save(product); err != nil {
				return err
			}
		}

		order.Set("status", "cancelled")
		return txApp.Save(order)
	})
}
