package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
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
		SCode    string
		SName    string
		Address  string
		Tel      string
		URL      string
		STime    string
		Capacity int
	}

	type Menu struct {
		MCode  string
		MName  string
		Price  int
		Detail string
	}

	router := gin.Default()

	// GET a store detail
	router.GET("/store/:scode", func(c *gin.Context) {
		var (
			store  Store
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
		url := c.PostForm("url")
		stime := c.PostForm("store_time")
		capa := c.PostForm("capacity")

		stmt, err := db.Prepare("insert into store (scode, sname, address, tel, url, stime, capacity) values(?,?,?,?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(scode, sname, address, tel, url, stime, capa)

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
		buffer.WriteString(capa)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// PUT - update a store details
	router.PUT("/store", func(c *gin.Context) {
		var buffer bytes.Buffer
		scode := c.Query("store_code")
		sname := c.PostForm("store_name")
		address := c.PostForm("address")
		tel := c.PostForm("tel")
		url := c.PostForm("url")
		stime := c.PostForm("store_time")
		capa := c.PostForm("capacity")
		stmt, err := db.Prepare("update store set sname= ?, address= ?, tel= ?, url= ?, stime= ?, capacity= ? where scode= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(sname, address, tel, url, stime, capa, scode)
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
		buffer.WriteString(capa)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources store
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
			"message": fmt.Sprintf("Successfully deleted Store_Code: %s", scode),
		})
	})

	//Get a menu detail
	router.GET("/menu/:mcode", func(c *gin.Context) {
		var (
			menu   Menu
			result gin.H
		)
		mcode := c.Param("mcode")
		row := db.QueryRow("select mcode, mname, price, detail from menu where mcode = ?;", mcode)
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
		price := c.PostForm("price")
		detail := c.PostForm("detail")

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
		buffer.WriteString(price)
		buffer.WriteString(" ")
		buffer.WriteString(detail)
		buffer.WriteString(" ")
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// PUT - update a menu details
	router.PUT("/menu", func(c *gin.Context) {
		var buffer bytes.Buffer
		mcode := c.Query("menu_code")
		mname := c.PostForm("menu_name")
		price := c.PostForm("price")
		detail := c.PostForm("detail")

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
		buffer.WriteString(price)
		buffer.WriteString(" ")
		buffer.WriteString(detail)
		buffer.WriteString(" ")

		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources menu
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
			"message": fmt.Sprintf("Successfully deleted Menu_Code: %s", mcode),
		})
	})

	router.Run(":80")
}
