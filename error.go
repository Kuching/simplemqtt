package smqtt 

import (
	"errors"
)

func InternalError(err_str string) error {
	return errors.New(err_str)
}

func ParamRequiredError(param_name string) error {
	return errors.New(param_name + " is required")
}

func ClientError(err_str string) error {
	return errors.New("Fail to init mqtt client: " + err_str)
}