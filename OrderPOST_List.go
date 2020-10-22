package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	OCode  string
	Day    string
	SCode  string
	Tabnam int
	Flg    int
}

type List struct{
	LCode    string
	MCode    string
	Quantity int
	OCode    string
}

func main() {
	db, err := sql.Open("mysql", "user1:seraku1!@tcp(localhost:3306)/pj")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	router := gin.Default()

	// POST a order details
	router.POST("/order", func(c *gin.Context) {
		var order  Order

		ocode := strconv.FormatInt(time.Now().Unix(), 16)
		day := time.Now().Format("2006-01-02")
		scode := c.PostForm("store_code")
		tabnam, _ := strconv.Atoi(c.PostForm("table_no"))
		mcode := c.PostForm("menu_code")
		quantity, _ := strconv.Atoi(c.PostForm("quantity"))
		flg := 0
		Order_flg := 0

		//Search Order table
		rows, err := db.Query("select ocode, day, scode, tabnam, flg from `order`;")
		if err != nil {
			fmt.Print(err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnam, &order.Flg)
			if err != nil {
				fmt.Print(err.Error())
			}
			if order.Flg == 0 && order.SCode == scode && order.Tabnam == tabnam {
				ocode = order.OCode
				Order_flg = 1
				break
			}
		}
		defer rows.Close()

		if Order_flg == 0 {
			stmt, err1 := db.Prepare("insert into `order` (ocode, day, scode, tabnam, flg) values(?,?,?,?,?);")
			if err1 != nil {
				fmt.Print(err1.Error())
			}
			_, err1 = stmt.Exec(ocode, day, scode, tabnam, flg)
			if err1 != nil {
				fmt.Print(err1.Error())
			}
		}

		stmt, err2 := db.Prepare("insert into list (mcode, quantity, ocode) values(?,?,?);")
		if err2 != nil {
			fmt.Print(err2.Error())
		}
		_, err2 = stmt.Exec(mcode, quantity, ocode)
		if err2 != nil {
			fmt.Print(err2.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully Created:Order OCode %s, day %s, SCode %s, Tabnam %d, Flg %d",ocode, day, scode, tabnam, flg),
		})

	Order_flg = 0
	})

	//Payment
	router.PUT("/payment", func(c *gin.Context){
		ocode := c.PostForm("order_code")

		stmt, err := db.Prepare("update `order` set flg = 1 where ocode = ?")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(ocode)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully payment : %s", ocode),
		})
	})

	//Get order,list table
	router.GET("/orders", func (c *gin.Context){
		var(
			order  Order
			orders []Order
			list   List
			lists  []List
		)
		rows1, err1 := db.Query("select ocode, day, scode, tabnam, flg from `order`; ")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnam, &order.Flg)
			orders = append(orders, order)
			if err != nil{
				fmt.Print(err1.Error())
			}
		}
		defer rows1.Close()

		rows2, err2 := db.Query("select lcode, mcode, quantity, ocode from list; ")
		if err2 != nil{
			fmt.Print(err2.Error())
		}
		for rows2.Next(){
			err2 = rows2.Scan(&list.LCode, &list.MCode, &list.Quantity, &list.OCode)
			lists = append(lists, list)
			if err != nil{
				fmt.Print(err.Error())
			}
		}
		defer rows2.Close()
		c.JSON(http.StatusOK, gin.H{
			"list_table":lists,
			"order_table":orders,
		})
	})

	//Get order,list table by a store
	router.GET("/order/:scode", func (c *gin.Context){
		var(
			order  Order
			orders []Order
			list   List
			lists  []List
		)
		scode := c.Param("scode")

		rows1, err1 := db.Query("select ocode, day, scode, tabnam, flg from `order`;")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnam, &order.Flg)
			if err1 != nil{
				fmt.Print(err1.Error())
			}

			if order.SCode == scode{
				orders = append(orders, order)

				rows2, err2 := db.Query("select lcode, mcode, quantity, ocode from list;")
				if err2 != nil{
					fmt.Print(err2.Error())
				}
				for rows2.Next(){
					err2 = rows2.Scan(&list.LCode, &list.MCode, &list.Quantity, &list.OCode)
					if err2 != nil{
						fmt.Print(err2.Error())
					}
					if list.OCode == order.OCode{
						lists = append(lists, list)
					}
				}
				defer rows2.Close()
			}
		}
		defer rows1.Close()

		c.JSON(http.StatusOK, gin.H{
			"Order":orders,
			"List":lists,
		})
	})

	//Get order,list table at a table by a store 
	router.GET("/order", func (c *gin.Context){
		var(
			order  Order
			orders []Order
			list   List
			lists  []List
		)
		scode := c.PostForm("store_code")
		tabnam, _ := strconv.Atoi(c.PostForm("table_no"))

		rows1, err1 := db.Query("select ocode, day, scode, tabnam, flg from `order`;")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnam, &order.Flg)
			if err1 != nil{
				fmt.Print(err1.Error())
			}

			if order.Tabnam == tabnam && order.SCode == scode{
				orders = append(orders, order)

				rows2, err2 := db.Query("select lcode, mcode, quantity, ocode from list;")
				if err2 != nil{
					fmt.Print(err2.Error())
				}
				for rows2.Next(){
					err2 = rows2.Scan(&list.LCode, &list.MCode, &list.Quantity, &list.OCode)
					if err2 != nil{
						fmt.Print(err2.Error())
					}
					if list.OCode == order.OCode{
						lists = append(lists, list)
					}
				}
				defer rows2.Close()
			}
		}
		defer rows1.Close()

		c.JSON(http.StatusOK, gin.H{
			"Order":orders,
			"List":lists,
		})
	})

	router.Run(":80")
}
