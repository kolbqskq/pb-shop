# PB-SHOP

Тестовая реализация магазина с помощью pocketbase

# Установка

**1. Запустить pocketbase.exe и дождаться установки**

**2.Запустить программу:**
```
go run ./cmd serve
```

**3.Зайти в админпанель и создать суперюзера:**
```
http://127.0.0.1:8090/_/
```

**4.Создать и настроить необходимые коллекции:**

## "categories":
    "name": Plain text, noneEmpty

## "products":
    "name": Plain text, noneEmpty
    "description": Plain text
    "price": Number, nonZero
    "stock": Number
    "is_active": Bool
    "category": Relation(categories), single

    API Rules:
        List/Search: any
        View: any

## "cart_items":
    "user": Relation(users), single, noneEmpty, cascade delete
    "product": Relation(product), single, noneEmpty, cascade delete
    "quantity": Number, nonZero

    API Rules:
        List/Search: @request.auth.id != "" && user = @request.auth.id
        View: @request.auth.id != "" && user = @request.auth.id
        Create: @request.auth.id != ""
        Update: @request.auth.id != "" && user = @request.auth.id
        Delete: @request.auth.id != "" && user = @request.auth.id

## "orders":
    "user" Relation(users), single, noneEmpty, cascade delete
    "status" Plain text, noneEmpty
    "total" Number

    API Rules:
        List/Search: @request.auth.id != "" && user = @request.auth.id
        View: @request.auth.id != "" && user = @request.auth.id
        Create: @request.auth.id != ""
        Update: @request.auth.id != "" && user = @request.auth.id
        Delete: @request.auth.id != "" && user = @request.auth.id

## "order_items":
    "order": Relation(orders), single, noneEmpty, cascade delete
    "product": Relation(products), single, noneEmpty, cascade delete
    "quantity": Number, nonZero
    "price": Number, nonZero