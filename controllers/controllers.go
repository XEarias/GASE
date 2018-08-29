package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/go-sql-driver/mysql"
)

// BaseController operations for Activities
type BaseController struct {
	beego.Controller
}

//MessageResponse =
type MessageResponse struct {
	Message       string              `json:"message,omitempty"`
	Code          uint16              `json:"code,omitempty"`
	PrettyMessage string              `json:"pretty_message,omitempty"`
	Errors        []map[string]string `json:"errors,omitempty"`
}

func init() {
	validation.SetDefaultMessage(map[string]string{
		"Required":     "This field is required",
		"Min":          "The min length requred is %d",
		"Max":          "The max length requred is %d",
		"Range":        "The range of the values is %d до %d",
		"MinSize":      "Longitud mínima permitida %d",
		"MaxSize":      "Minimum length allowed %d",
		"Length":       "The length should be equal to %d",
		"Alpha":        "Must consist of letters",
		"Numeric":      "Must consist of numbers",
		"AlphaNumeric": "Must consist of letters or numbers",
		"Match":        "Must coincide with %s",
		"NoMatch":      "It should not coincide with %s",
		"AlphaDash":    "Must consist of letters, numbers or symbols (-_)",
		"Email":        "Must be in the correct email format",
		"IP":           "Must be a valid IP address",
		"Base64":       "Must be presented in the correct format base64",
		"Mobile":       "Must be the correct mobile number",
		"Tel":          "Must be the correct phone number",
		"Phone":        "Must be the correct phone or mobile number",
		"ZipCode":      "Must be the correct zip code",
	})
}

//ServeErrorJSON : Serve Json error
func (c *BaseController) ServeErrorJSON(err error) {

	if driverErr, ok := err.(*mysql.MySQLError); ok {

		switch driverErr.Number {
		case 1062:
			c.Ctx.Output.SetStatus(409)
			c.Data["json"] = MessageResponse{
				Message:       "The element already exists",
				Code:          driverErr.Number,
				PrettyMessage: "El elemento de la base de datos ya existe",
			}
		default:
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = MessageResponse{
				Message:       "An error has ocurred",
				Code:          driverErr.Number,
				PrettyMessage: "Un error ha ocurrido",
			}
		}
	}

	c.ServeJSON()
}

//BadRequest =
func (c *BaseController) BadRequest() {
	c.Ctx.Output.SetStatus(400)
	c.Data["json"] = MessageResponse{
		Message:       "Bad request body",
		PrettyMessage: "Peticion mal formada",
	}
	c.ServeJSON()
}

//BadRequestErrors =
func (c *BaseController) BadRequestErrors(errors []*validation.Error) {

	var errorsMessages []map[string]string

	for _, err := range errors {

		var errorMessage = map[string]string{
			"message":    err.Message,
			"field":      err.Field,
			"validation": err.Name,
		}

		errorsMessages = append(errorsMessages, errorMessage)
	}

	c.Ctx.Output.SetStatus(400)
	c.Data["json"] = MessageResponse{
		Errors: errorsMessages,
	}

	c.ServeJSON()
}