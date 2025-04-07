package handlers

import (
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	firstName := c.Query("firstName")
	secondName := c.Query("secondName")
	thirdName := c.Query("thirdName")
	fmt.Println(firstName, secondName, thirdName, " ", rand.Int()+500, "$")
}
