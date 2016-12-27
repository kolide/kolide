package kolide

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionMarshaller(t *testing.T) {
	tests := []struct {
		value         string
		typ           OptionType
		expected      interface{}
		expectSuccess bool
	}{
		{"23", OptionTypeInt, float64(23), true},
		{"abc", OptionTypeInt, nil, false},
		{"true", OptionTypeFlag, true, true},
		{"false", OptionTypeFlag, false, true},
		{"something", OptionTypeFlag, nil, false},
		{"foobar", OptionTypeString, "foobar", true},
	}

	for _, test := range tests {
		optIn := &Option{1, "foo", test.typ, &test.value, true}
		buff, err := json.Marshal(optIn)
		if !test.expectSuccess {
			assert.NotNil(t, err)
			continue
		}
		require.Nil(t, err)
		optOut := &Option{}
		err = json.Unmarshal(buff, optOut)
		require.Nil(t, err)
		assert.Equal(t, optIn.ID, optOut.ID)
		assert.Equal(t, optIn.Name, optOut.Name)
		assert.Equal(t, optIn.ReadOnly, optOut.ReadOnly)
		assert.Equal(t, optIn.Type, optOut.Type)
		assert.Equal(t, test.expected, optOut.Value)

	}

	// test nil
	optIn := &Option{1, "bar", OptionTypeString, nil, true}
	buff, err := json.Marshal(optIn)
	require.Nil(t, err)
	optOut := &Option{}
	err = json.Unmarshal(buff, optOut)
	require.Nil(t, err)
	assert.True(t, reflect.DeepEqual(optIn, optOut))

}

func TestOptionUnmarshaller(t *testing.T) {
	t.Skip("test option unmarshaller")
	errTypeMismatch := fmt.Errorf("option value type mismatch")

	tests := []struct {
		data string
		err  error
	}{
		{`{"id":1,"name":"foo","type":"string","value":"foobar","read_only":true}`, nil},
		{`{"id":1,"name":"foo","type":"float","value":"foobar","read_only":true}`, fmt.Errorf("option type 'float' invalid")},
		{`{"id":1,"name":"foo","type":"int","value":"foobar","read_only":true}`, errTypeMismatch},
		{`{"id":1,"name":"foo","type":"flag","value":"foobar","read_only":true}`, errTypeMismatch},
	}

	for _, test := range tests {
		buff := []byte(test.data)
		opt := &Option{}
		err := json.Unmarshal(buff, opt)
		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
	}

}
