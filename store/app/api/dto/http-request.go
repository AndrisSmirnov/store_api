package dto

import "store/app/api/types"

type RequestFromId struct {
	Id uint32 `json:"id" validate:"required,gt=0"`
}

type RequesCreatetUser struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	ParentId uint32 `json:"parentId" validate:"omitempty,min=1"`
}

type RequestParents struct {
	UserId uint32 `json:"userId" validate:"required,gt=0"`
}

type RequestChildrens struct {
	UserId uint32 `json:"userId" validate:"required,gt=0"`
	Type   uint16 `json:"type" validate:"required,oneof=1 2 3"`
	Deep   uint16 `json:"deep" validate:"required,gt=0"`
}

type RequestCategory struct {
	Name string `json:"name" validate:"required"`
}

type RequestProduct struct {
	Name       string `json:"name" validate:"required"`
	Price      uint32 `json:"price" validate:"required,gt=0"`
	Count      uint32 `json:"count" validate:"required,gt=0"`
	CategoryId uint32 `json:"categoryId" validate:"required,gt=0"`
}

type RequestTransaction struct {
	RequestOrder []RequestOrder `json:"requestorder" validate:"required,dive"`
}

type RequestOrder struct {
	ProductId uint32 `json:"productId" validate:"required,min=1,max=100000"`
	Count     uint32 `json:"count" validate:"required,min=1,max=100000"`
}

type RequestShowTransaction struct {
	UserId   uint32                `json:"userId" validate:"required,min=1"`
	Type     types.TransactionType `json:"type" validate:"required,oneof=1 2 3"`
	UnixTime uint64                `json:"unixtime" validate:"omitempty,min=1633046400"`
	Limit    int                   `json:"limit"  mod:"default=5" validate:"omitempty,gt=0"`
	Offset   int                   `json:"offset" mod:"default=1" validate:"omitempty,gt=0"`
}
