package envreader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvReader(t *testing.T) {
	t.Run("Fails if passed nil", func(t *testing.T) {
		err := Load(nil)

		assert.EqualError(t, err, "Load expected to be passed a pointer to a struct, recieved a(n) invalid instead")
	})

	t.Run("Fails if passed a nil pointer", func(t *testing.T) {
		type testStruct struct {
			unexported bool
		}

		var test testStruct

		err := Load(test)

		assert.EqualError(t, err, "Load expected to be passed a pointer to a struct, recieved a(n) struct instead")
	})

	t.Run("Fails if passed by value", func(t *testing.T) {
		type testStruct struct {
			unexported bool
		}

		err := Load(testStruct{})

		assert.EqualError(t, err, "Load expected to be passed a pointer to a struct, recieved a(n) struct instead")
	})

	t.Run("Fails if passed a pointer to something random", func(t *testing.T) {
		input := 5
		err := Load(&input)

		assert.EqualError(t, err, "Load expected to be passed a pointer to a struct, recieved a pointer to a(n) int instead")
	})

	t.Run("Fails if passed a struct with no tagged fields", func(t *testing.T) {
		type testStruct struct {
			Param bool
		}

		err := Load(&testStruct{})

		assert.EqualError(t, err, "No tagged fields were found for struct type 'testStruct'")
	})

	t.Run("Fails if environment value not set", func(t *testing.T) {
		type testStruct struct {
			Param interface{} `env:"PARAM_1"`
		}

		test := testStruct{}
		err := Load(&test)

		assert.EqualError(t, err, "Failed to load field 'Param': Environment variable 'PARAM_1' not set")
	})

	t.Run("Fails if unsupported field type used", func(t *testing.T) {
		type testStruct struct {
			Param interface{} `env:"PARAM_1"`
		}

		os.Setenv("PARAM_1", "value")

		test := testStruct{}
		err := Load(&test)

		assert.EqualError(t, err, "Failed to load field 'Param': Field type 'interface' is not supported")
	})

	t.Run("Fails if bool field is given arbitrary data", func(t *testing.T) {
		type testStruct struct {
			BoolParam bool `env:"PARAM_1"`
		}

		os.Setenv("PARAM_1", "value")

		test := testStruct{}
		err := Load(&test)

		assert.EqualError(t, err, "Failed to load field 'BoolParam': Error converting variable 'PARAM_1' ('value') to boolean: strconv.ParseBool: parsing \"value\": invalid syntax")
	})

	t.Run("Fails if int field is given arbitrary data", func(t *testing.T) {
		type testStruct struct {
			IntParam int `env:"PARAM_1"`
		}

		os.Setenv("PARAM_1", "value")

		test := testStruct{}
		err := Load(&test)

		assert.EqualError(t, err, "Failed to load field 'IntParam': Error converting variable 'PARAM_1' ('value') to integer: strconv.Atoi: parsing \"value\": invalid syntax")
	})

	t.Run("Loads data", func(t *testing.T) {
		type testStruct struct {
			StringParam  string `env:"PARAM_1"`
			BoolParam    bool   `env:"PARAM_2"`
			IntParam     int    `env:"PARAM_3"`
			IgnoredParam string
		}

		os.Setenv("PARAM_1", "value")
		os.Setenv("PARAM_2", "true")
		os.Setenv("PARAM_3", "1")

		test := testStruct{}
		err := Load(&test)

		assert.NoError(t, err)
		assert.Equal(
			t,
			testStruct{
				StringParam:  "value",
				BoolParam:    true,
				IntParam:     1,
				IgnoredParam: "",
			},
			test,
		)
	})

	t.Run("Loads nested data", func(t *testing.T) {
		type subStruct struct {
			StringParam string `env:"PARAM_3"`
		}

		type testStruct struct {
			StringParam string `env:"PARAM_1"`
			BoolParam   bool   `env:"PARAM_2"`
			StructParam subStruct
		}

		os.Setenv("PARAM_1", "value")
		os.Setenv("PARAM_2", "true")
		os.Setenv("PARAM_3", "nested")

		test := testStruct{}
		err := Load(&test)

		assert.NoError(t, err)
		assert.Equal(
			t,
			testStruct{
				StringParam: "value",
				BoolParam:   true,
				StructParam: subStruct{StringParam: "nested"},
			},
			test,
		)
	})
}
