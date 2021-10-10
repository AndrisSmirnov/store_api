package voc

const (
	HTTP_INVALID_REQUEST = "invalid request, wrong data"
	HTTP_NO_HEADER       = "no Authorization header provided"
	HTTP_WRONG_HEADER    = "could not find bearer token in Authorization header"
	HTTP_BAD_TOKEN       = "the token is not correct"
	TOKEN_EXPIRED        = "token expired"
	PERMISSION_ERROR     = "You don't have permission for this operation."

	HTTP_CLIENT_NOT_FOUND       = "client not found"
	HTTP_WROND_PASSWORD         = "wrong password"
	HTTP_CLIENT_DELETED         = "user was deleted"
	HTTP_CLIENT_WITHOUT_FOUNDS  = "client has not enough funds"
	HTTP_CATEGORY_NOT_FOUND     = "category not found"
	HTTP_CATEGORY_DELETED       = "category was deleted"
	HTTP_PRODUCT_NOT_FOUND      = "product not found"
	HTTP_PRODUCT_DELETED        = "product was deleted"
	HTTP_PRODUCTS_OUT_STOCK     = "products out of stock"
	HTTP_PRODUCTS_NOT_AVAILABLE = "some products are not available"
	HTTP_TRANSACTION_NOT_FOUND  = "transaction not found"
	HTTP_TRANSACTION_DELETED    = "transaction was deleted"

	ERROR_CONNECT_REDIS_RECONNECT = "error Connect to Redis. Waiting for reconnection..."
	ERROR_CONNECT_REDIS           = "error Connect to Redis"
	ERROR_REDIS_NIL_OBJ           = "object pointer is empty"
	HTTP_CLIENT_ALREADY_EXIST     = "client already exists"
)
