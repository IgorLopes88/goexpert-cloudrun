package main

import (
	// "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertZipcode(t *testing.T) {
	zipcode, err := convertZipcode("11666-000")
	assert.Nil(t, err)
	assert.Equal(t, "11666000", zipcode)

	zipcode, err = convertZipcode("123ABc")
	assert.Error(t, err, "invalid zipcode")
	assert.Equal(t, "", zipcode)
}

func TestConvertName(t *testing.T) {
	city := convertName("Ilhabela")
	assert.Equal(t, "Ilhabela", city)

	city = convertName("São Paulo")
	assert.Equal(t, "Sao%20Paulo", city)
}

func TestSearchLocation(t *testing.T) {
	city, err := SearchLocation("11600-300")
	assert.Nil(t, err)
	assert.Equal(t, "São Sebastião", city)

	city, err = SearchLocation("11600")
	assert.Error(t, err, "can not find zipcode")
	assert.Equal(t, "", city)
}

func TestGetTemperature(t *testing.T) {
	temperature, err := GetTemperature("Ilhabela")
	assert.Nil(t, err)
	assert.NotEqual(t, 0.00, temperature)

	temperature, _ = GetTemperature("xyz")
	assert.Equal(t, 0.00, temperature)
}
