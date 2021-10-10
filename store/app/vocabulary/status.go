package voc

import "store/app/services/database/model"

const (
	USER_STATUS_ACTIVE         model.StatusEnum = iota + 1
	USER_STATUS_DELETED                         = 2
	CATEGORY_STATUS_DELETED                     = 2
	PRODUCT_STATUS_DELETED                      = 2
	TRANSACTION_STATUS_DELETED                  = 2
)
