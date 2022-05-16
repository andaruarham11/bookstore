package configs

const MongoUri = "mongodb://admin:admin123@127.0.0.1:27000/?&w=majority"

const DefaultDBName = "bookstore"

// Order configurations
const (
	OrderDBName   = DefaultDBName
	OrderCollName = "orders"
)

// Payment configurations
const (
	PaymentDBName   = DefaultDBName
	PaymentCollName = "payments"
)

// Book configurations
const (
	BookDBName   = DefaultDBName
	BookCollName = "books"
)

// User configurations
const (
	UserDBName   = DefaultDBName
	UserCollName = "users"
)
