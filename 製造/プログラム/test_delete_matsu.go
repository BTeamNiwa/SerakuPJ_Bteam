  // Delete resources
    router.DELETE("/list", func(c *gin.Context) {
        ocode := c.Query("order_code")
        stmt, err := db.Prepare("delete from list where lcode= ?;")
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(lcode)
        if err != nil {
            fmt.Print(err.Error())
        }
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully deleted user: %s", ocode, lcode),
        })
    })