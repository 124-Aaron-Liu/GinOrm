package main

import (
	"fmt"
	"gosql/Model"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Must Bind跟Should Bind二種
// 1.Must Bind：binding拋出err時會回傳http 400錯誤碼
// 2.Should Bind：binding拋出err時就沒東西而己

// BindQuery :取得query string 並解析
// BindJson  :取得body 並解析

// time_format tag is only used in form binding, not json;
type PostUserInfo struct {
	Username   string          `json:"username" binding:"required"`
	Department string          `json:"department" binding:"required"`
	Created    Model.LocalTime `json:"created" binding:"required"`
}

type GetRequest struct {
	Id int `form:"id" binding:"required"`
}

func main() {
	route := gin.Default()

	route.GET("/userInfo/:id", getUserInfo)
	// /userInfo/byQueryStr?id=3
	route.GET("/userInfo/byQueryStr", getUserInfoByQueryStr)

	route.POST("/userInfo", postUserInfo)
	route.Run(":8000")
}

func getUserInfoByQueryStr(c *gin.Context) {

	var request GetRequest
	if c.Bind(&request) == nil {
		log.Println("====== Bind By Query String ======")
		log.Println(request.Id)
	}
	id := request.Id
	userInfo, err := Model.GetUserInfo(id)

	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"param": userInfo})
}

func getUserInfo(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		panic(err)
	}

	userInfo, err := Model.GetUserInfo(id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"param": userInfo})

}

func postUserInfo(c *gin.Context) {
	var userInfo PostUserInfo

	if err := c.ShouldBindJSON(&userInfo); err != nil {
		fmt.Println("BindJson fault", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	Model.Create(userInfo.Username, userInfo.Department, userInfo.Created)
	fmt.Printf("customer:%+v", userInfo)
	c.JSON(200, gin.H{"param": userInfo})
}
