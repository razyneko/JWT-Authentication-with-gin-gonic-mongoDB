package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// for ensuring that only admins get all users data
func CheckUserType(c *gin.Context, role string)(err error){
	userType := c.GetString("userType")
	err = nil
	if userType != role{
		err = errors.New("unauthorized to access this resource")
	}
	return err
}

func MatchUserTypeToUId(c *gin.Context, userId string) (err error) {
	userType := c.GetString("userType")
	uid := c.GetString("userId")
	err = nil

	// not allowing user to access other user's data
	if userType == "USER" && uid != userId{
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}