package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONProvider_Decode(t *testing.T) {
	type jsonTest struct {
		Name    string
		Age     int
		Weight  float64
		Address struct {
			City    string
			Country string
		} `json:"address,omitempty"`
	}

	t.Run("Success", func(t *testing.T) {
		p := JSONProvider{}
		data := jsonTest{}
		err := p.Decode([]byte(`{"Name":"John Doe","Age":30,"Weight":70.5,"address":{"City":"New-York","Country":"USA"}}`), &data)
		require.NoError(t, err)
		require.Equal(t, "John Doe", data.Name)
		require.Equal(t, 30, data.Age)
		require.Equal(t, 70.5, data.Weight)
		require.Equal(t, "New-York", data.Address.City)
		require.Equal(t, "USA", data.Address.Country)
	})
	t.Run("Error-1", func(t *testing.T) {
		p := JSONProvider{}
		data := jsonTest{}
		err := p.Decode([]byte(`{"Name":"John Doe","Age":30,"Weight":70.5,"Address":{"City":"New-York","Country":"USA"}`), &data)
		require.Error(t, err)
		require.Equal(t, "unexpected end of JSON input", err.Error())
	})
	t.Run("Error-2", func(t *testing.T) {
		p := JSONProvider{}
		data := jsonTest{}
		err := p.Decode([]byte(`{"Name":"John Doe","Age":"30","Weight":70.5,"Address":{"City":"New-York","Country":"USA"}}`), &data)
		require.Error(t, err)
		require.Equal(t, "json: cannot unmarshal string into Go struct field jsonTest.Age of type int", err.Error())
	})
}

func TestJSONProvider_Encode(t *testing.T) {
	type jsonTest struct {
		Name    string
		Age     int
		Weight  float64
		Address struct {
			City    string
			Country string
		}
	}

	p := JSONProvider{}
	data := jsonTest{
		Name:   "John Doe",
		Age:    30,
		Weight: 70.5,
		Address: struct {
			City    string
			Country string
		}{City: "New-York", Country: "USA"},
	}
	b, err := p.Encode(data)
	require.NoError(t, err)
	require.Equal(t, `{"Name":"John Doe","Age":30,"Weight":70.5,"Address":{"City":"New-York","Country":"USA"}}`, string(b))
}
