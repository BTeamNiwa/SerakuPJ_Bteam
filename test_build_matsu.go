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

    router := gin.Default()

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
            "message": fmt.Sprintf("Successfully deleted user: %s", mcode),
        })
    })

    router.Run(":80")
}
