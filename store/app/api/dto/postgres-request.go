package dto

import (
	"time"

	types "store/app/api/types"
	model "store/app/services/database/model"
)

type ResponseUser struct {
	Id         uint32           `json:"id"`
	Name       string           `json:"name"`
	Surname    string           `json:"surname"`
	Login      string           `json:"login"`
	Email      string           `json:"email"`
	Password   string           `json:"password"`
	Balance    uint32           `json:"balance"`
	Status     model.StatusEnum `json:"status"`
	Permission types.Permission `json:"permission"`
	ParentId   uint32           `json:"parentId"`
	Createtime time.Time        `json:"createtime"`
	Updatetime time.Time        `json:"updatatime"`
}

type ResponseCategory struct {
	Id     uint32           `json:"id"`
	Name   string           `json:"name"`
	Status model.StatusEnum `json:"status"`
}

type Products struct {
	Product []ResponseProduct `json:"product"`
}

type ResponseProduct struct {
	Id         uint32           `json:"id"`
	Name       string           `json:"name"`
	Price      uint32           `json:"price"`
	Status     model.StatusEnum `json:"status"`
	Count      uint32           `json:"count"`
	CategoryId uint32           `json:"categoryId"`
}

type ResponseTransaction struct {
	Id         uint32           `json:"id"`
	Balance    uint32           `json:"balance"`
	Amount     uint32           `json:"amount"`
	Types      uint32           `json:"types"`
	Status     model.StatusEnum `json:"status"`
	Createtime time.Time        `json:"createtime"`
	Updatetime time.Time        `json:"updatetime"`
	UserId     uint32           `json:"userId"`
	Basket     string           `json:"basket"`
}

type TransactionProdacts struct {
	TransactionId uint32 `json:"transactionId"`
	ProductId     uint32 `json:"productId"`
}

type TranProdTable struct {
	Raw []TransactionProdacts `json:"raw"`
}

type WinnersTable struct {
	Id          uint32
	UserId      uint32
	Amount      uint32
	Basket      string
	PId         uint32
	Price       uint32
	PCategoryId uint32
}
