package util

import (
	"math"
	"math/big"
	"strings"
)

// Add adds two decimal strings with optional scale
func Add(a, b string, scale *int) string {
	s := MaxScale(a, b)
	if scale != nil {
		s = *scale
	}

	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)

	result := new(big.Float).Add(af, bf)
	return FormatFloat(result, s)
}

// Sub subtracts two decimal strings with optional scale
func Sub(a, b string, scale *int) string {
	s := MaxScale(a, b)
	if scale != nil {
		s = *scale
	}

	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)

	result := new(big.Float).Sub(af, bf)
	return FormatFloat(result, s)
}

// Mul multiplies two decimal strings with optional scale
func Mul(a, b string, scale *int) string {
	s := MaxScale(a, b)
	if scale != nil {
		s = *scale
	}

	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)

	result := new(big.Float).Mul(af, bf)
	return FormatFloat(result, s)
}

// Div divides two decimal strings with optional scale
func Div(a, b string, scale *int) *string {
	s := MaxScale(a, b)
	if scale != nil {
		s = *scale
	}

	if Equal(b, "0") {
		return nil
	}

	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)

	result := new(big.Float).Quo(af, bf)
	str := FormatFloat(result, s)
	return &str
}

// Pow raises a to the power of b
func Pow(a, b string) string {
	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)

	afloat, _ := af.Float64()
	bfloat, _ := bf.Float64()
	result := new(big.Float).SetFloat64(math.Pow(afloat, bfloat))
	return FormatFloat(result, MaxScale(a, b))
}

// Mod returns modulo of a divided by b
func Mod(a, b string) *string {
	af := new(big.Int)
	bf := new(big.Int)
	af.SetString(a, 10)
	bf.SetString(b, 10)

	result := new(big.Int).Mod(af, bf)
	str := result.String()
	return &str
}

// Equal returns true if a equals b
func Equal(a, b string) bool {
	return Compare(a, b, MaxScale(a, b)) == 0
}

// Eq is an alias for Equal
func Eq(a, b string) bool {
	return Equal(a, b)
}

// Ne returns true if a does not equal b
func Ne(a, b string) bool {
	return !Equal(a, b)
}

// Gt returns true if a is greater than b
func Gt(a, b string) bool {
	return Compare(a, b, MaxScale(a, b)) == 1
}

// Lt returns true if a is less than b
func Lt(a, b string) bool {
	return Compare(a, b, MaxScale(a, b)) == -1
}

// Gte returns true if a is greater than or equal to b
func Gte(a, b string) bool {
	return Compare(a, b, MaxScale(a, b)) >= 0
}

// Lte returns true if a is less than or equal to b
func Lte(a, b string) bool {
	return Compare(a, b, MaxScale(a, b)) <= 0
}

// Floor returns the largest integer less than or equal to a
func Floor(a string, scale int) string {
	if !IsNumeric(a) {
		return "0"
	}

	af := new(big.Float)
	af.SetString(a)

	if Gte(a, "0") || scale >= GetScale(a) {
		return FormatFloat(af, scale)
	}

	pow := new(big.Float).SetFloat64(math.Pow(0.1, float64(scale)))
	result := new(big.Float).Sub(af, pow)
	return FormatFloat(result, scale)
}

// Ceil returns the smallest integer greater than or equal to a
func Ceil(a string, scale int) string {
	if !IsNumeric(a) {
		return "0"
	}

	af := new(big.Float)
	af.SetString(a)

	if Lte(a, "0") || scale >= GetScale(a) {
		return FormatFloat(af, scale)
	}

	pow := new(big.Float).SetFloat64(math.Pow(0.1, float64(scale)))
	result := new(big.Float).Add(af, pow)
	return FormatFloat(result, scale)
}

// Round returns the nearest integer to a
func Round(a string, scale int) string {
	if !IsNumeric(a) {
		return "0"
	}

	af := new(big.Float)
	af.SetString(a)

	pow := new(big.Float).SetFloat64(math.Pow(0.1, float64(scale)))
	half := new(big.Float).Quo(pow, new(big.Float).SetFloat64(2))

	if Gte(a, "0") {
		result := new(big.Float).Add(af, half)
		return FormatFloat(result, scale)
	}

	result := new(big.Float).Sub(af, half)
	return FormatFloat(result, scale)
}

// Sqrt returns the square root of a with optional scale
func Sqrt(a string, scale *int) *string {
	if !IsNumeric(a) || Lt(a, "0") {
		return nil
	}

	s := GetScale(a)
	if scale != nil {
		s = *scale
	}

	af := new(big.Float)
	af.SetString(a)

	afloat, _ := af.Float64()
	result := new(big.Float).SetFloat64(math.Sqrt(afloat))
	str := FormatFloat(result, s)
	return &str
}

// Helper functions

func GetScale(a string) int {
	if !IsNumeric(a) || !strings.Contains(a, ".") {
		return 0
	}
	parts := strings.Split(a, ".")
	return len(parts[1])
}

func MaxScale(a, b string) int {
	return int(math.Max(float64(GetScale(a)), float64(GetScale(b))))
}

func IsNumeric(s string) bool {
	_, ok := new(big.Float).SetString(s)
	return ok
}

func Compare(a, b string, scale int) int {
	af := new(big.Float)
	bf := new(big.Float)
	af.SetString(a)
	bf.SetString(b)
	return af.Cmp(bf)
}

func FormatFloat(f *big.Float, scale int) string {
	return f.Text('f', scale)
}
