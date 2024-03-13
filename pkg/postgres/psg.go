package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type ProductMod struct {
	DB *sql.DB
}

type Auth struct {
	us int
}

//type Sign struct {
//	user int
//}

type ListProd struct {
	Products []Prod
}

type Answer struct {
	CountAdd    int
	CountUpdate int
	CountDelete int
	CountErrors int
}
type Prod struct {
	OffId int
	Id    int
	Name  string
	Price int
	Count int
}

func (s *ProductMod) CheckPassLogUser(user, pass string) (int, error) {
	//sign := Sign{}
	_, err := s.DB.Query("SELECT * FROM users WHERE email = $1 AND password_hash = $2", user, pass)
	if err != nil {
		return 1, err
	}
	return 0, err
}

func (s *ProductMod) InputInfo(email, password_hash string) (int, error) {
	auth := Auth{}
	err := s.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&auth.us)
	if err != nil {
		return auth.us, err
	}
	if auth.us > 0 {
		return auth.us, nil
	} else {
		_, err := s.DB.Exec("INSERT INTO users (email,  password_hash) VALUES ($1, $2)", email, password_hash)
		if err != nil {
			return auth.us, err
		}
	}
	return auth.us, nil
}

func (s *ProductMod) DownloadInfo(tmp, seller_id string) (Answer, error) {
	ans := Answer{}
	var q int
	seller_name := "Seller"
	for _, kkk := range strings.Split(tmp, "\r\n") {
		gen := strings.Split(kkk, ",")
		_, err := s.DB.Exec("INSERT INTO sellers (seller_id, seller_name) VALUES ($1, $2) ON CONFLICT (seller_id) DO UPDATE SET seller_name = EXCLUDED.seller_name", seller_id, seller_name)
		if err != nil {
			//log.Fatal(err)
			ans.CountErrors++
			return ans, err
		}
		offer_id := gen[0]
		name := gen[1]
		price := gen[2]
		quantity := gen[3]
		available := gen[4]
		if available == "TRUE" {
			err := s.DB.QueryRow("SELECT COUNT(*) FROM products WHERE offer_id = $1 AND seller_id = $2", offer_id, seller_id).Scan(&q)
			if err != nil {
				//log.Fatal(err)
				ans.CountErrors++
				return ans, err
			}
			if q > 0 {
				_, err := s.DB.Exec("UPDATE products SET name = $1, price = $2, quantity = $3 WHERE offer_id = $4 AND seller_id = $5", name, price, quantity, offer_id, seller_id)
				if err != nil {
					//log.Fatal(err)
					ans.CountErrors++
					return ans, err
				}
				ans.CountUpdate++
			} else {
				_, err := s.DB.Exec("INSERT INTO products (offer_id, seller_id, name, price, quantity) VALUES ($1, $2, $3, $4, $5)", offer_id, seller_id, name, price, quantity)
				if err != nil {
					//log.Fatal(err)
					ans.CountErrors++
					return ans, err
				}
				ans.CountAdd++
			}
		} else {
			_, err := s.DB.Exec("DELETE FROM products WHERE offer_id = $1 AND seller_id = $2", offer_id, seller_id)
			if err != nil {
				ans.CountErrors++
				return ans, err
			}
			ans.CountDelete++

		}
	}
	return ans, nil
}

func (s *ProductMod) SearchInfo(offer_id, seller_id, name string) ([]Prod, error) {
	var query string
	ql := 0
	var queryParams []interface{}
	query = "SELECT * FROM products WHERE name ILIKE $1"
	queryParams = append(queryParams, "%"+name+"%")
	if offer_id != " " {
		ql++
		query += " AND offer_id = $2"
		queryParams = append(queryParams, offer_id)
	}

	if seller_id != " " {
		if ql == 0 {
			query += " AND seller_id = $2"
		} else {
			query += " AND seller_id = $3"

		}
		queryParams = append(queryParams, seller_id)
	}

	PutInf, err := s.DB.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer PutInf.Close()
	var minProd []Prod
	for PutInf.Next() {

		var prodMin Prod
		var priceStr string
		err := PutInf.Scan(&prodMin.OffId, &prodMin.Id, &prodMin.Name, &priceStr, &prodMin.Count)
		if err != nil {
			return nil, err
		}
		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, err
		}

		prodMin.Price = int(priceFloat)
		minProd = append(minProd, prodMin)
		fmt.Println(minProd)
	}

	return minProd, nil
}
