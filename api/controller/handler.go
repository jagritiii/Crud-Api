package controller

import (
	"awesomeProject/api/service"
	"awesomeProject/pkg/enums"
	"awesomeProject/pkg/models"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

func createUser(e echo.Context) error {
	var reqbody []models.User
	err := json.NewDecoder(e.Request().Body).Decode(&reqbody)
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusBadRequest, enums.Faileddecode)
	}
	v := validator.New()
	for _, val := range reqbody {
		err = v.Struct(&val)
		if err != nil {
			log.Error(err.Error())
			return e.JSON(http.StatusBadRequest, enums.Validation)
		}
	}
	err = service.CreateUser(reqbody)
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusInternalServerError, enums.ServerIssue)
	}
	return e.JSON(http.StatusOK, enums.Statusok)
}
func getalluser(e echo.Context) error {
	res, err := service.Getalluser()
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusInternalServerError, enums.ServerIssue)
	}
	return e.JSON(http.StatusOK, res)
}
func updateuser(e echo.Context) error {
	var reqbody []models.User
	err := json.NewDecoder(e.Request().Body).Decode(&reqbody)
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusBadRequest, enums.Faileddecode)
	}
	v := validator.New()
	for _, val := range reqbody {
		err = v.Struct(&val)
		if err != nil {
			log.Error(err.Error())
			return e.JSON(http.StatusBadRequest, enums.Validation)
		}
	}
	err = service.UpdateUser(reqbody)
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusInternalServerError, enums.ServerIssue)
	}
	return e.JSON(http.StatusOK, enums.Statusok)
}
func getauser(e echo.Context) error {
	str, err := strconv.Atoi(e.QueryParam("id"))
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusBadRequest, enums.Faileddecode)
	}
	res, err1, err2 := service.Getauser(str)
	if err1 != nil {
		return e.JSON(http.StatusAccepted, err1.Error())
	} else if err2 != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusInternalServerError, enums.ServerIssue)
	}
	return e.JSON(http.StatusOK, res)
}
func deleteuser(e echo.Context) error {
	str, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return e.JSON(http.StatusBadRequest, enums.Faileddecode)
	}
	err1, err2 := service.Deleteuser(str)
	if err1 != nil && err2 == nil {
		return e.JSON(http.StatusInternalServerError, enums.ServerIssue)
	} else if err1 != nil && err2 != nil {
		return e.JSON(http.StatusInternalServerError, "no such key exist")
	}
	return e.JSON(http.StatusOK, enums.Deletesuccess)

}
