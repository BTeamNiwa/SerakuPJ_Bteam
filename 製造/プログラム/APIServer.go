package main

import (
    "bytes"
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    gin.SetMode(gin.ReleaseMode)
    db, err := sql.Open("mysql", "user1:Seraku1!@tcp(localhost:3306)/pj")
    if err != nil {
        fmt.Print(err.Error())
    }
    defer db.Close()
    // make sure connection is available
    err = db.Ping()
    if err != nil {
        fmt.Print(err.Error())
    }


    type Store struct {
        SCode         string
        SName         string
        Address       string
	Tel           string
	URL           string
	STime         string
	Capacity      int
    }

    type Menu struct{
	MCode        string
	MName        string
	Price        int
	Detail       string
    }


type Order struct {
	OCode  string
	Day    string
	SCode  string
	Tabnum int
	Flg    int
}

type List struct{
	LCode    string
	MCode    string
	Quantity int
	OCode    string
}

    router := gin.Default()


// GET a store detail
    router.GET("/store/:scode", func(c *gin.Context) {
        var (
            store Store
            result gin.H
        )
        scode := c.Param("scode")
        row := db.QueryRow("select scode, sname, address, tel, url, stime, capacity from store where scode = ?;", scode)
        err = row.Scan(&store.SCode, &store.SName, &store.Address, &store.Tel, &store.URL, &store.STime, &store.Capacity)
        if err != nil {
            // If no results send null
            result = gin.H{
                "result": nil,
                "count":  0,
            }
        } else {
            result = gin.H{
                "result": store,
                "count":  1,
            }
        }
        c.JSON(http.StatusOK, result)
    })


    // GET all stores
    router.GET("/stores", func(c *gin.Context) {
        var (
            store  Store
            stores []Store
        )
        rows, err := db.Query("select scode, sname, address, tel, url, stime, capacity from store;")
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&store.SCode, &store.SName, &store.Address, &store.Tel, &store.URL, &store.STime, &store.Capacity)
            stores = append(stores, store)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.JSON(http.StatusOK, gin.H{
            "result": stores,
            "count":  len(stores),
        })
    })

    // POST new store details
    router.POST("/store", func(c *gin.Context) {
        var buffer bytes.Buffer

        scode := c.PostForm("store_code")
        sname := c.PostForm("store_name")
        address := c.PostForm("address")
        tel := c.PostForm("tel")
        url:= c.PostForm("url")
        stime:= c.PostForm("store_time")
        capacity,_:= strconv.Atoi(c.PostForm("capacity"))

	if capacity < 0{
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf("input(capacity data %d < 0) is not valid.", capacity),
        	})
	}else{

	        stmt, err := db.Prepare("insert into store (scode, sname, address, tel, url, stime, capacity) values(?,?,?,?,?,?,?);")
        	if err != nil {
	            fmt.Print(err.Error())
        	}
	        _, err = stmt.Exec(scode, sname, address, tel, url, stime, capacity)

        	if err != nil {
	            fmt.Print(err.Error())
        	}

	        // Fastest way to append strings
        	buffer.WriteString(scode)
	        buffer.WriteString(" ")
        	buffer.WriteString(sname)
	        buffer.WriteString(" ")
        	buffer.WriteString(address)
	        buffer.WriteString(" ")
        	buffer.WriteString(tel)
	        buffer.WriteString(" ")
        	buffer.WriteString(url)
	        buffer.WriteString(" ")
        	buffer.WriteString(stime)
	        buffer.WriteString(" ")
        	buffer.WriteString(strconv.Itoa(capacity))
	        buffer.WriteString(" ")

	        defer stmt.Close()
        	name := buffer.String()
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf(" %s successfully created", name),
	        })
	}
    })

    // PUT - update a store details
    router.PUT("/store", func(c *gin.Context) {
        var buffer bytes.Buffer
        scode := c.Query("store_code")
        sname := c.PostForm("store_name")
        address := c.PostForm("address")
	tel := c.PostForm("tel")
        url:= c.PostForm("url")
        stime:= c.PostForm("store_time")
        capacity,_ := strconv.Atoi(c.PostForm("capacity"))

	if capacity < 0{
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf("input(capacity data %d < 0) is not valid.", capacity),
        	})
	}else{

	        stmt, err := db.Prepare("update store set sname= ?, address= ?, tel= ?, url=?, stime=?, capacity=? where scode= ?;")
        	if err != nil {
	            fmt.Print(err.Error())
        	}
	        _, err = stmt.Exec(sname, address, tel, url, stime, capacity, scode)
        	if err != nil {
	            fmt.Print(err.Error())
        	}

	        // Fastest way to append strings
        	buffer.WriteString(scode)
	        buffer.WriteString(" ")
        	buffer.WriteString(sname)
	        buffer.WriteString(" ")
        	buffer.WriteString(address)
	        buffer.WriteString(" ")
        	buffer.WriteString(tel)
	        buffer.WriteString(" ")
		buffer.WriteString(url)
	        buffer.WriteString(" ")
		buffer.WriteString(stime)
	        buffer.WriteString(" ")
        	buffer.WriteString(strconv.Itoa(capacity))

	        defer stmt.Close()
	        name := buffer.String()
        	c.JSON(http.StatusOK, gin.H{
	            "message": fmt.Sprintf("Successfully updated to %s", name),
	        })
	}
    })

    // Delete resources
    router.DELETE("/store", func(c *gin.Context) {
        scode := c.Query("store_code")
        stmt, err := db.Prepare("delete from store where scode= ?;")
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(scode)
        if err != nil {
            fmt.Print(err.Error())
        }
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully deleted user: %s", scode),
        })
    })


    // GET a menu detail
    router.GET("/menu/:mcode", func(c *gin.Context) {
        var (
            menu Menu
            result gin.H
        )
        mcode := c.Param("mcode")
        row := db.QueryRow("select mcode, mname,price,detail from menu where mcode = ?;", mcode)
        err = row.Scan(&menu.MCode, &menu.MName, &menu.Price, &menu.Detail)
        if err != nil {
            // If no results send null
            result = gin.H{
                "result": nil,
                "count":  0,
            }
        } else {
            result = gin.H{
                "result": menu,
                "count":  1,
            }
        }
        c.JSON(http.StatusOK, result)
    })


    // GET all menus
    router.GET("/menus", func(c *gin.Context) {
        var (
            menu  Menu
            menus []Menu
        )
        rows, err := db.Query("select mcode, mname, price, detail from menu;")
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&menu.MCode, &menu.MName, &menu.Price, &menu.Detail)
            menus = append(menus, menu)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.JSON(http.StatusOK, gin.H{
            "result": menus,
            "count":  len(menus),
        })
    })

    // POST new menu details
    router.POST("/menu", func(c *gin.Context) {
        var buffer bytes.Buffer

        mcode := c.PostForm("menu_code")
        mname := c.PostForm("menu_name")
        price,_ := strconv.Atoi(c.PostForm("price"))
        detail := c.PostForm("detail")

	if price < 0{
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf("input(price data %d < 0) is not valid.", price),
	        })
	}else{

	        stmt, err := db.Prepare("insert into menu (mcode, mname, price, detail) values(?,?,?,?);")
        	if err != nil {
	            fmt.Print(err.Error())
        	}
	        _, err = stmt.Exec(mcode, mname, price, detail)

        	if err != nil {
	            fmt.Print(err.Error())
        	}

	        // Fastest way to append strings
        	buffer.WriteString(mcode)
	        buffer.WriteString(" ")
        	buffer.WriteString(mname)
	        buffer.WriteString(" ")
        	buffer.WriteString(strconv.Itoa(price))
	        buffer.WriteString(" ")
        	buffer.WriteString(detail)

	        defer stmt.Close()
        	name := buffer.String()
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf(" %s successfully created", name),
	        })
	}
    })

    // PUT - update a menu details
    router.PUT("/menu", func(c *gin.Context) {
        var buffer bytes.Buffer
        mcode := c.Query("menu_code")
        mname := c.PostForm("menu_name")
        price,_ := strconv.Atoi(c.PostForm("price"))
	detail := c.PostForm("detail")

	if price < 0{
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf("input(price data %d < 0) is not valid.", price),
	        })
	}else{

	        stmt, err := db.Prepare("update menu set mname= ?, price= ?, detail= ? where mcode= ?;")
        	if err != nil {
	            fmt.Print(err.Error())
        	}
	        _, err = stmt.Exec(mname, price, detail, mcode)
        	if err != nil {
	            fmt.Print(err.Error())
        	}

	        // Fastest way to append strings
        	buffer.WriteString(mcode)
	        buffer.WriteString(" ")
        	buffer.WriteString(mname)
	        buffer.WriteString(" ")
        	buffer.WriteString(strconv.Itoa(price))
        	buffer.WriteString(" ")
	        buffer.WriteString(detail)

	        defer stmt.Close()
        	name := buffer.String()
	        c.JSON(http.StatusOK, gin.H{
        	    "message": fmt.Sprintf("Successfully updated to %s", name),
	        })
	}
    })

    // Delete resources
    router.DELETE("/menu", func(c *gin.Context) {
        mcode := c.Query("menu_code")
        stmt, err := db.Prepare("delete from menu where mcode= ?;")
        if err != nil {
            fmt.Print(err.Error())
    }
        _, err = stmt.Exec(mcode)
        if err != nil {
            fmt.Print(err.Error())
        }
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully deleted menu: %s", mcode),
        })
    })




//以下・追加統合分（オーダー・リスト）

//Get order,list table
	router.GET("/orders", func (c *gin.Context){
		var(
			order  Order
			orders []Order
			list   List
			lists  []List
		)
		rows1, err1 := db.Query("select ocode, day, scode, tabnum, flg from `order`; ")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnum, &order.Flg)
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

		rows1, err1 := db.Query("select ocode, day, scode, tabnum, flg from `order`;")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnum, &order.Flg)
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
		tabnum, _ := strconv.Atoi(c.PostForm("table_no"))

		rows1, err1 := db.Query("select ocode, day, scode, tabnum, flg from `order`;")
		if err1 != nil{
			fmt.Print(err1.Error())
		}
		for rows1.Next(){
			err1 = rows1.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnum, &order.Flg)
			if err1 != nil{
				fmt.Print(err1.Error())
			}

			if order.Tabnum == tabnum && order.SCode == scode{
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



	// POST a order details
	router.POST("/order", func(c *gin.Context) {
		var order  Order

		ocode := strconv.FormatInt(time.Now().Unix(), 16)
		day := time.Now().Format("2006-01-02")
		scode := c.PostForm("store_code")
		tabnum, _ := strconv.Atoi(c.PostForm("table_no"))
		mcode := c.PostForm("menu_code")
		quantity, _ := strconv.Atoi(c.PostForm("quantity"))
		flg := 0
		Order_flg := 0

		if tabnum < 0 || quantity < 0{
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Input(tabnum %d or quantity %d) is not valid.",tabnum, quantity),
			})
		}else{
			//Search Order table
			rows, err := db.Query("select ocode, day, scode, tabnum, flg from `order`;")
			if err != nil {
				fmt.Print(err.Error())
			}

			for rows.Next() {
				err = rows.Scan(&order.OCode, &order.Day, &order.SCode, &order.Tabnum, &order.Flg)
				if err != nil {
					fmt.Print(err.Error())
				}
				if order.Flg == 0 && order.SCode == scode && order.Tabnum == tabnum {
					ocode = order.OCode
					Order_flg = 1
					break
				}
			}
			defer rows.Close()

			if Order_flg == 0 {
				stmt, err1 := db.Prepare("insert into `order` (ocode, day, scode, tabnum, flg) values(?,?,?,?,?);")
				if err1 != nil {
					fmt.Print(err1.Error())
				}
				_, err1 = stmt.Exec(ocode, day, scode, tabnum, flg)
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
				"message": fmt.Sprintf("Successfully Created:Order OCode %s, day %s, SCode %s, Tabnum %d, Flg %d",ocode, day, scode, tabnum, flg),
			})
		}
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

	router.DELETE("/order", func(c *gin.Context) {
		ocode := c.Query("order_code")

		stmt1, err1 := db.Prepare("delete list from list left join `order` on list.ocode = `order`.ocode where list.ocode=? and`order`.flg=0;")
		if err1 != nil {
			fmt.Print(err1.Error())
		}
		_, err1 = stmt1.Exec(ocode)
		if err1 != nil {
			fmt.Print(err1.Error())
		}
		
		stmt2, err2 := db.Prepare("delete from `order` where ocode = ? and flg = 0;")
		if err2 != nil {
			fmt.Print(err2.Error())
		}
		_, err2 = stmt2.Exec(ocode)
		if err2 != nil {
			fmt.Print(err2.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted Order_Code: %s", ocode),
		})
	})

	router.DELETE("/all_order", func(c *gin.Context){
		_, err1 := db.Query("delete from list;")
		if err1 != nil {
			fmt.Print(err1.Error())
		}
		
		_, err2 := db.Query("delete from `order`;")
		if err2 != nil {
			fmt.Print(err2.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted All Order."),
		})

	})



    router.Run(":80")
}

