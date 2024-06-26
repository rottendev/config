package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTomlProvider_Decode(t *testing.T) {
	type tomlTest struct {
		Name    string
		Age     int
		Weight  float64
		Married bool
		Address struct {
			City    string
			Country string
		}
		Hobbies []string
	}

	t.Run("Success", func(t *testing.T) {
		p := TomlProvider{}
		data := tomlTest{}
		err := p.Decode([]byte(`
				# this is comment 
				Name = "John Doe"
				Age = 30
				Weight = 70.5
				Married = true
				Hobbies = ["reading", "swimming"]
				[Address]
					City = "New-York"
					Country = "USA"
		`), &data)
		require.NoError(t, err)
		require.Equal(t, "John Doe", data.Name)
		require.Equal(t, 30, data.Age)
		require.Equal(t, 70.5, data.Weight)
		require.True(t, data.Married)
		require.Len(t, data.Hobbies, 2)
		require.Equal(t, "New-York", data.Address.City)
		require.Equal(t, "USA", data.Address.Country)
	})

	t.Run("Error-1", func(t *testing.T) {
		p := TomlProvider{}
		data := tomlTest{}
		err := p.Decode([]byte(`
				Name = "John Doe"
				Age = "30"
				Weight = 70.5
				Address = {City = "New-York", Country = "USA"}
				Hobbies = ["reading", "swimming"]
		`), &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "value has type string; destination has type integer")
	})

	t.Run("Error-2", func(t *testing.T) {
		p := TomlProvider{}
		data := tomlTest{}
		err := p.Decode([]byte(`
				Name = "John Doe"
				Age = 30
				Weight = 70.5
				Address = {City = "New-York", Country = "USA"}
				Hobbies = ["reading", "swimming"
		`), &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "but got end of file")
	})
}
