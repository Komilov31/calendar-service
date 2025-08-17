package validator

import (
	"fmt"
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("date_after_now", dateAfterNow)
}

// функция подготовит сообщение ошибки в случае ошибки валидации поля
func CreateValidationErrorResponse(err error) string {
	var msg string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fe.Field())
			case "date_after_now":
				msg = fmt.Sprintf("%s must be a date in the future", fe.Field())
			default:
				msg = fmt.Sprintf("%s is not valid due to %s", fe.Field(), fe.Tag())
			}
		}
	}
	return msg
}

func dateAfterNow(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(model.Date)
	if !ok {
		return false
	}

	t := time.Time(date)
	return t.After(time.Now())
}
