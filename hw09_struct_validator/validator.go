package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	errProgram        = fmt.Errorf("program error occurred during validation")
	errValueNotInSet  = fmt.Errorf("value is not in given set")
	errRegexpMismatch = fmt.Errorf("value does not match to regexp")
	errLessThanMin    = fmt.Errorf("value is less than min")
	errBiggerThanMax  = fmt.Errorf("value is bigger than max")
	errLengthNotEqual = fmt.Errorf("value is not equal to length")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var allErrors string
	for _, validationError := range v {
		allErrors += fmt.Sprintf("field: %s, error: %s\n", validationError.Field, validationError.Err.Error())
	}
	return allErrors
}

//nolint:gocognit // Высокая когнитивная сложность обусловлена необходимостью обработки множества валидаторов.
func Validate(v interface{}) error {
	errs := ValidationErrors{}

	reflectValue := reflect.ValueOf(v)
	reflectType := reflect.TypeOf(v)
	if reflectValue.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct type")
	}

	for i := 0; i < reflectType.NumField(); i++ {
		tag := reflectType.Field(i).Tag.Get("validate")
		if tag == "" {
			continue // пропуск полей без тега "validate".
		}
		validators := strings.Split(tag, "|")
		fieldType := reflectType.Field(i).Type
		fieldValue := reflectValue.Field(i)

		for _, validator := range validators {
			validatorDevided := strings.Split(validator, ":")
			validatorName := validatorDevided[0]
			validatorValue := validatorDevided[1]

			switch validatorName {
			case "len":
				if err := validateLen(fieldType, fieldValue, validatorValue); err != nil {
					if errors.Is(err, errProgram) {
						return fmt.Errorf("failed to validateLen: %w", err)
					}
					errs = append(errs, ValidationError{
						Field: reflectType.Field(i).Name,
						Err:   err,
					})
				}
			case "max":
				if err := validateMax(fieldType, fieldValue, validatorValue); err != nil {
					if errors.Is(err, errProgram) {
						return fmt.Errorf("failed to validateMax: %w", err)
					}
					errs = append(errs, ValidationError{
						Field: reflectType.Field(i).Name,
						Err:   err,
					})
				}
			case "min":
				if err := validateMin(fieldType, fieldValue, validatorValue); err != nil {
					if errors.Is(err, errProgram) {
						return fmt.Errorf("failed to validateMin: %w", err)
					}
					errs = append(errs, ValidationError{
						Field: reflectType.Field(i).Name,
						Err:   err,
					})
				}
			case "regexp":
				if err := validateRegexp(fieldType, fieldValue, validatorValue); err != nil {
					if errors.Is(err, errProgram) {
						return fmt.Errorf("failed to validateRegexp: %w", err)
					}
					errs = append(errs, ValidationError{
						Field: reflectType.Field(i).Name,
						Err:   err,
					})
				}
			case "in":
				if err := validateIn(fieldType, fieldValue, validatorValue); err != nil {
					if errors.Is(err, errProgram) {
						return fmt.Errorf("failed to validateIn: %w", err)
					}
					errs = append(errs, ValidationError{
						Field: reflectType.Field(i).Name,
						Err:   err,
					})
				}
			default:
				return fmt.Errorf("unknown validator name: %s", validatorName)
			}
		}
	}

	return errs
}

func validateLen(fieldType reflect.Type, fieldValue reflect.Value, validatorValue string) error {
	switch {
	case fieldType.Kind() == reflect.String:
		if err := validateStringLen(fieldValue.String(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.String:
		strSlice := fieldValue.Interface().([]string)
		for _, str := range strSlice {
			if err := validateStringLen(str, validatorValue); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("wrong type, string expected")
	}

	return nil
}

func validateMax(fieldType reflect.Type, fieldValue reflect.Value, validatorValue string) error {
	switch {
	case fieldType.Kind() == reflect.Int:
		if err := validateIntMax(fieldValue.Int(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Int:
		intSlice := fieldValue.Interface().([]int)
		for _, intVal := range intSlice {
			if err := validateIntMax(int64(intVal), validatorValue); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("wrong type, int expected")
	}

	return nil
}

func validateMin(fieldType reflect.Type, fieldValue reflect.Value, validatorValue string) error {
	switch {
	case fieldType.Kind() == reflect.Int:
		if err := validateIntMin(fieldValue.Int(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Int:
		intSlice := fieldValue.Interface().([]int)
		for _, intVal := range intSlice {
			if err := validateIntMin(int64(intVal), validatorValue); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("wrong type, int expected")
	}

	return nil
}

func validateRegexp(fieldType reflect.Type, fieldValue reflect.Value, validatorValue string) error {
	switch {
	case fieldType.Kind() == reflect.String:
		if err := validateStringRegexp(fieldValue.String(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.String:
		strSlice := fieldValue.Interface().([]string)
		for _, str := range strSlice {
			if err := validateStringRegexp(str, validatorValue); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("wrong type, string expected")
	}

	return nil
}

func validateIn(fieldType reflect.Type, fieldValue reflect.Value, validatorValue string) error {
	switch {
	case fieldType.Kind() == reflect.Int:
		if err := validateIntIn(fieldValue.Int(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Int:
		intSlice := fieldValue.Interface().([]int)
		for _, intVal := range intSlice {
			if err := validateIntIn(int64(intVal), validatorValue); err != nil {
				return err
			}
		}
	case fieldType.Kind() == reflect.String:
		if err := validateStrIn(fieldValue.String(), validatorValue); err != nil {
			return err
		}
	case fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.String:
		strSlice := fieldValue.Interface().([]string)
		for _, strVal := range strSlice {
			if err := validateStrIn(strVal, validatorValue); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("wrong type, string or int expected")
	}

	return nil
}

func validateStringLen(strValue string, validatorValue string) error {
	requiredLen, err := strconv.Atoi(validatorValue)
	if err != nil {
		return err
	}
	runeCount := utf8.RuneCountInString(strValue)

	if runeCount != requiredLen {
		return errLengthNotEqual
	}

	return nil
}

func validateIntMax(intValue int64, validatorValue string) error {
	requiredMaxValue, err := strconv.ParseInt(validatorValue, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert string length to int: %w", errProgram)
	}

	if intValue > requiredMaxValue {
		return errBiggerThanMax
	}

	return nil
}

func validateIntMin(intValue int64, validatorValue string) error {
	requiredMinValue, err := strconv.ParseInt(validatorValue, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert string length to int: %w", errProgram)
	}

	if intValue < requiredMinValue {
		return errLessThanMin
	}

	return nil
}

func validateStringRegexp(strValue string, regularExp string) error {
	matchString, err := regexp.MatchString(regularExp, strValue)
	if err != nil {
		return fmt.Errorf("failed to regexp.MatchString: %w", errProgram)
	}

	if !matchString {
		return errRegexpMismatch
	}

	return nil
}

func validateIntIn(intValue int64, validatorValue string) error {
	inValuesStr := strings.Split(validatorValue, ",")
	for _, v := range inValuesStr {
		inValue, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("failed to convert string to int: %w", errProgram)
		}

		if intValue == int64(inValue) {
			return nil
		}
	}

	return errValueNotInSet
}

func validateStrIn(strValue string, validatorValue string) error {
	inValuesStr := strings.Split(validatorValue, ",")
	for _, v := range inValuesStr {
		if strValue == v {
			return nil
		}
	}

	return errValueNotInSet
}
