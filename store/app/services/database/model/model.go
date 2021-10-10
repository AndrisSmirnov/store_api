package model

import (
	"store/app/api/types"
	"time"
)

type StatusEnum uint16

const (
	STATUS_ACTIVE StatusEnum = iota + 1
	STATUS_DELETED
)

type User struct {
	Id         uint32           `json:"id" gorm:"primary_key"`
	Name       string           `json:"name"`
	Surname    string           `json:"surname"`
	Login      string           `json:"login" gorm:"unique"`
	Email      string           `json:"email"`
	Password   string           `json:"password"`
	Balance    uint32           `json:"balance"`
	Status     StatusEnum       `json:"status" gorm:"type:STATUSENUM;"`
	Permission types.Permission `json:"permission"`
	ParentId   uint32           `json:"parentId"`
	Createtime time.Time        `json:"createtime"`
	Updatetime time.Time        `json:"updatatime"`
}

type Transaction struct {
	Id         uint32     `json:"id" gorm:"primary_key"`
	Balance    uint32     `json:"balance"`
	Amount     uint32     `json:"amount"`
	Types      uint32     `json:"types"`
	Status     StatusEnum `json:"status" gorm:"type:STATUSENUM;"`
	Createtime time.Time  `json:"createtime"`
	Updatetime time.Time  `json:"updatetime"`
	UserId     uint32     `json:"userId"`
	User       User       `gorm:"foreignKey:UserId;references:Id"`
	Basket     string     `json:"basket"`
	Product    []Product  `gorm:"many2many:transaction_products;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Category struct {
	Id     uint32     `json:"id" gorm:"primary_key"`
	Name   string     `json:"name" gorm:"unique"`
	Status StatusEnum `json:"status" gorm:"type:STATUSENUM;"`
}

func (Category) TableName() string { return "categories" }

type Product struct {
	Id         uint32     `json:"id" gorm:"primary_key"`
	Name       string     `json:"name" gorm:"unique"`
	Price      uint32     `json:"price"`
	Status     StatusEnum `json:"status" gorm:"type:STATUSENUM;"`
	Count      uint32     `json:"count"`
	CategoryId uint32     `json:"categoryId"`
	Category   Category   `gorm:"foreignKey:CategoryId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
