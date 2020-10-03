package envreader

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Load takes a struct pointer and reads the values for any properties marked with an `env:"VAR_NAME"` tag from the environment
func Load(output interface{}) error {
	outputValue := reflect.ValueOf(output)

	if kind := outputValue.Kind().String(); kind != "ptr" {
		return fmt.Errorf("Load expected to be passed a pointer to a struct, recieved a(n) %s instead", kind)
	}

	if outputValue.IsNil() {
		return errors.New("Failed to load from environment: output struct is invalid")
	}

	structValue := reflect.Indirect(outputValue)
	if kind := structValue.Kind().String(); kind != "struct" {
		return fmt.Errorf("Load expected to be passed a pointer to a struct, recieved a pointer to a(n) %s instead", kind)
	}

	numFields := structValue.NumField()
	if numFields == 0 {
		return fmt.Errorf("Cannot load type %s from environment: must have at least one exported field", outputValue.Kind().String())
	}

	numLoaded, err := loadStruct(structValue)
	if err != nil {
		return err
	}

	if numLoaded == 0 {
		return fmt.Errorf("No tagged fields were found for struct type '%s'", structValue.Type().Name())
	}

	return nil
}

func loadStruct(structValue reflect.Value) (int, error) {
	numLoaded := 0
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		loaded, err := loadField(field, fieldType)
		if err != nil {
			return numLoaded, fmt.Errorf("Failed to load field '%s': %+v", fieldType.Name, err)
		}

		if loaded {
			numLoaded++
		}
	}

	return numLoaded, nil
}

func loadField(field reflect.Value, fieldType reflect.StructField) (bool, error) {
	if field.Kind() == reflect.Struct {
		numLoaded, err := loadStruct(field)
		return numLoaded > 0, err
	}

	env := fieldType.Tag.Get("env")
	if env == "" {
		return false, nil
	}

	value, ok := os.LookupEnv(env)
	if !ok {
		return false, fmt.Errorf("Environment variable '%s' not set", env)
	}

	switch kind := field.Kind(); kind {
	case reflect.String:
		field.Set(reflect.ValueOf(value))
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("Error converting variable '%s' ('%s') to boolean: %+v", env, value, err)
		}

		field.Set(reflect.ValueOf(boolValue))
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return false, fmt.Errorf("Error converting variable '%s' ('%s') to integer: %+v", env, value, err)
		}

		field.Set(reflect.ValueOf(intValue))
	default:
		return false, fmt.Errorf("Field type '%s' is not supported", kind.String())
	}

	return true, nil
}
