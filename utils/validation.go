package utils

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/go-playground/validator/v10"
	"log"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	err := validate.RegisterValidation("role", validateRole)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := models.GetValidRoles()
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
