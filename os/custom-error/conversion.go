package customerror

import (
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/errors"
)

// Conversion Error
func Conversion(err error, svrName string, method string) error {
	if err == nil {
		return nil
	}
	if err.Error() == "not found" {
		return errors.NotFound(svrName+"."+method, "not found")
	}
	return err
}

// BadRequest Error
func BadRequest(svrName string, method string, str string) error {
	return errors.BadRequest(svrName+"."+method, str)
}

// ValidateShopID validate shop id
func ValidateShopID(shopID string, svrName string, method string) error {
	if len(shopID) < 6 {
		return BadRequest(svrName, method, "invalid Shop ID")
	}

	return nil
}

// ValidateName validate name
func ValidateName(name string, svrName string, method string) error {
	if len(name) < 2 {
		return BadRequest(svrName, method, "invalid name")
	}

	return nil
}

// ValidateShopIDAndName validate shop id and name
func ValidateShopIDAndName(shopID string, name string, svrName string, method string) error {
	if err := ValidateShopID(shopID, svrName, method); err != nil {
		return err
	}
	return ValidateName(name, svrName, method)
}

// ValidateID validate id
func ValidateID(id string, svrName string, method string) error {

	if len(id) < 5 {
		return BadRequest(svrName, method, "invalid id")
	}

	return nil
}

// ValidateShopIDAndID validate shop id and name
func ValidateShopIDAndID(shopID string, id string, svrName string, method string) error {
	if err := ValidateShopID(shopID, svrName, method); err != nil {
		return err
	}
	return ValidateID(id, svrName, method)
}

// WriteError write error to response
func WriteError(err error, rsp *restful.Response) {
	error := errors.Parse(err.Error())
	if error.Code == 0 {
		rsp.WriteError(500, err)
	} else {
		rsp.WriteError(int(error.Code), error)
	}
}