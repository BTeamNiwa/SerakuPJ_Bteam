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

	router.Run(":80")
}
