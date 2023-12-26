package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	role struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	roles = map[int]*role{}
	seq   = 1
	lock  = sync.Mutex{}
)

func getAllRoles(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, roles)
}
func createRole(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	r := &role{
		ID: seq,
	}
	err := c.Bind(r)
	if err != nil {
		log.Printf("Failed addUserRequest: %s", err)
		return err
	}
	roles[r.ID] = r
	seq++
	return c.JSON(http.StatusCreated, r)
}
func getRole(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	//Atoi<=>parsint()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, roles[id])
}
func updateRole(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	r := new(role)
	err := c.Bind(r)
	if err != nil {
		log.Printf("Failed updateUserRequest: %s", err)
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	roles[id].Name = r.Name

	return c.JSON(http.StatusOK, r)
}
func deleteRole(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(roles, id)
	return c.NoContent(http.StatusNoContent)

}
func main() {

	e := echo.New()

	e.Use(middleware.Logger()) // 👈 log all requests

	// Routes
	e.GET("/roles", getAllRoles)
	e.POST("/roles", createRole)
	e.GET("/roles/:id", getRole)
	e.PUT("/roles/:id", updateRole)
	e.DELETE("/roles/:id", deleteRole)

	e.Start(":3000")
}
