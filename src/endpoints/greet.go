package endpoints

import (
	"go-api/src/store"
	"net/http"

	"github.com/labstack/echo/v4"
)



func RegisterGreetingsRoutes(e *echo.Echo){
	prefix := "/greetings"
	e.GET(prefix + "/greet",getGreetHandler)
	e.POST(prefix + "/hello",postHelloHandler)
}


// Request struct with validation tags
type HelloRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

// Response struct
type HelloResponse struct {
	Message string `json:"message"`
}

type GreetQuery struct {
	Name  string `query:"name" validate:"required,min=2"`
	Title string `query:"title" validate:"omitempty,min=2"`
}

func getGreetHandler(c echo.Context) error {
	var params GreetQuery
	if err:= c.Bind(&params);err !=nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err:= c.Validate((&params)); err !=nil{
		return echo.NewHTTPError(http.StatusUnprocessableEntity,err.Error())
	}

	if err:= store.Cache().Set("greet", "Hello "+ params.Title+ " " + params.Name + "!",0); err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError)	
	}

	greeting, err := store.Cache().Get("greet")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	response := HelloResponse{Message:greeting}
	return c.JSON(http.StatusOK,response)
}


func postHelloHandler(c echo.Context) error {
	var req HelloRequest

	if err:=c.Bind(&req); err != nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	} 
	
	if err:=c.Validate(&req); err != nil{
		return echo.NewHTTPError(http.StatusUnprocessableEntity,err.Error())
	} 
	response := HelloResponse{Message: "Hello "+ req.Name + "!"}
	return c.JSON(http.StatusOK,response)
}