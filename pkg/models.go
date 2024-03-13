package models

type Users struct {
	email         string
	password_hash string
	cookie_key    string
}

type Clients struct {
	seller_id   int
	seller_name string
}

type Product struct {
	offer_id  int
	seller_id int
	name      string
	price     int
	quantity  int
}
