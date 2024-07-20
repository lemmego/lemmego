package vee

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Errors map[string][]string

func (e Errors) Error() string {
	val, _ := json.Marshal(e)
	return string(val)
}

type Vee struct {
	Errors
}

func New() *Vee {
	return &Vee{
		Errors: make(map[string][]string),
	}
}

func (v *Vee) AddError(field, message string) {
	v.Errors[field] = append(v.Errors[field], message)
}

func (v *Vee) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Vee) ErrorsJSON() map[string][]string {
	return v.Errors
}

// Required checks if the value is not empty
func (v *Vee) Required(field string, value interface{}) bool {
	if value == nil || value == "" {
		v.AddError(field, "This field is required")
		return false
	}
	return true
}

// Min checks if the value is greater than or equal to the minimum
func (v *Vee) Min(field string, value int, min int) bool {
	if value < min {
		v.AddError(field, "This field must be at least "+strconv.Itoa(min))
		return false
	}
	return true
}

// Max checks if the value is less than or equal to the maximum
func (v *Vee) Max(field string, value int, max int) bool {
	if value > max {
		v.AddError(field, "This field must not exceed "+strconv.Itoa(max))
		return false
	}
	return true
}

// Between checks if the value is between min and max (inclusive)
func (v *Vee) Between(field string, value int, min int, max int) bool {
	if value < min || value > max {
		v.AddError(field, "This field must be between "+strconv.Itoa(min)+" and "+strconv.Itoa(max))
		return false
	}
	return true
}

// Email checks if the value is a valid email address
func (v *Vee) Email(field string, value string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(value) {
		v.AddError(field, "This field must be a valid email address")
		return false
	}
	return true
}

// Alpha checks if the value contains only alphabetic characters
func (v *Vee) Alpha(field string, value string) bool {
	for _, char := range value {
		if !unicode.IsLetter(char) {
			v.AddError(field, "This field must contain only alphabetic characters")
			return false
		}
	}
	return true
}

// Numeric checks if the value contains only numeric characters
func (v *Vee) Numeric(field string, value string) bool {
	for _, char := range value {
		if !unicode.IsDigit(char) {
			v.AddError(field, "This field must contain only numeric characters")
			return false
		}
	}
	return true
}

// AlphaNumeric checks if the value contains only alphanumeric characters
func (v *Vee) AlphaNumeric(field string, value string) bool {
	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			v.AddError(field, "This field must contain only alphanumeric characters")
			return false
		}
	}
	return true
}

// Date checks if the value is a valid date in the specified format
func (v *Vee) Date(field string, value string, layout string) bool {
	_, err := time.Parse(layout, value)
	if err != nil {
		v.AddError(field, "This field must be a valid date in the format "+layout)
		return false
	}
	return true
}

// In checks if the value is in the given slice of valid values
func (v *Vee) In(field string, value string, validValues []string) bool {
	for _, validValue := range validValues {
		if value == validValue {
			return true
		}
	}
	v.AddError(field, "This field must be one of the following: "+strings.Join(validValues, ", "))
	return false
}

// Regex checks if the value matches the given regular expression
func (v *Vee) Regex(field string, value string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		v.AddError(field, "Invalid regular expression pattern")
		return false
	}
	if !regex.MatchString(value) {
		v.AddError(field, "This field must match the pattern: "+pattern)
		return false
	}
	return true
}

// URL checks if the value is a valid URL
func (v *Vee) URL(field string, value string) bool {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		v.AddError(field, "This field must be a valid URL")
		return false
	}
	return true
}

// IP checks if the value is a valid IP address (v4 or v6)
func (v *Vee) IP(field string, value string) bool {
	ip := net.ParseIP(value)
	if ip == nil {
		v.AddError(field, "This field must be a valid IP address")
		return false
	}
	return true
}

// UUID checks if the value is a valid UUID
func (v *Vee) UUID(field string, value string) bool {
	_, err := uuid.Parse(value)
	if err != nil {
		v.AddError(field, "This field must be a valid UUID")
		return false
	}
	return true
}

// Boolean checks if the value is a valid boolean
func (v *Vee) Boolean(field string, value interface{}) bool {
	switch value.(type) {
	case bool:
		return true
	case string:
		lowercaseValue := strings.ToLower(value.(string))
		if lowercaseValue == "true" || lowercaseValue == "false" {
			return true
		}
	case int:
		intValue := value.(int)
		if intValue == 0 || intValue == 1 {
			return true
		}
	}
	v.AddError(field, "This field must be a boolean value")
	return false
}

// JSON checks if the value is a valid JSON string
func (v *Vee) JSON(field string, value string) bool {
	var js json.RawMessage
	if json.Unmarshal([]byte(value), &js) != nil {
		v.AddError(field, "This field must be a valid JSON string")
		return false
	}
	return true
}

// AfterDate checks if the date is after the specified date
func (v *Vee) AfterDate(field string, value time.Time, afterDate time.Time) bool {
	if value.After(afterDate) {
		return true
	}
	v.AddError(field, "This field must be a date after "+afterDate.String())
	return false
}

// BeforeDate checks if the date is before the specified date
func (v *Vee) BeforeDate(field string, value time.Time, beforeDate time.Time) bool {
	if value.Before(beforeDate) {
		return true
	}
	v.AddError(field, "This field must be a date before "+beforeDate.String())
	return false
}

// StartsWith checks if the string starts with the specified substring
func (v *Vee) StartsWith(field string, value string, prefix string) bool {
	if strings.HasPrefix(value, prefix) {
		return true
	}
	v.AddError(field, "This field must start with "+prefix)
	return false
}

// EndsWith checks if the string ends with the specified substring
func (v *Vee) EndsWith(field string, value string, suffix string) bool {
	if strings.HasSuffix(value, suffix) {
		return true
	}
	v.AddError(field, "This field must end with "+suffix)
	return false
}

// Contains checks if the string contains the specified substring
func (v *Vee) Contains(field string, value string, substring string) bool {
	if strings.Contains(value, substring) {
		return true
	}
	v.AddError(field, "This field must contain "+substring)
	return false
}

// Dimensions checks if the image file has the specified dimensions
func (v *Vee) Dimensions(field string, filepath string, width, height int) bool {
	file, err := os.Open(filepath)
	if err != nil {
		v.AddError(field, "Unable to open the file")
		return false
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		v.AddError(field, "Unable to decode the image")
		return false
	}

	if img.Width != width || img.Height != height {
		v.AddError(field, "Image dimensions must be "+strconv.Itoa(width)+"x"+strconv.Itoa(height))
		return false
	}
	return true
}

// MimeTypes checks if the file has one of the specified MIME types
func (v *Vee) MimeTypes(field string, filepath string, allowedTypes []string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		v.AddError(field, "Unable to open the file")
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		v.AddError(field, "Unable to read the file")
		return false
	}

	mimeType := http.DetectContentType(buffer)

	for _, allowedType := range allowedTypes {
		if mimeType == allowedType {
			return true
		}
	}

	v.AddError(field, "File type must be one of: "+strings.Join(allowedTypes, ", "))
	return false
}

// Timezone checks if the value is a valid timezone
func (v *Vee) Timezone(field string, value string) bool {
	_, err := time.LoadLocation(value)
	if err != nil {
		v.AddError(field, "Invalid timezone")
		return false
	}
	return true
}

// ActiveURL checks if the URL is active and reachable
func (v *Vee) ActiveURL(field string, value string) bool {
	resp, err := http.Get(value)
	if err != nil {
		v.AddError(field, "The URL is not active or reachable")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		v.AddError(field, "The URL returned a non-OK status")
		return false
	}
	return true
}

// AlphaDash checks if the string contains only alpha-numeric characters, dashes, or underscores
func (v *Vee) AlphaDash(field string, value string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	if !re.MatchString(value) {
		v.AddError(field, "This field may only contain alpha-numeric characters, dashes, and underscores")
		return false
	}
	return true
}

// Ascii checks if the string contains only ASCII characters
func (v *Vee) Ascii(field string, value string) bool {
	for _, char := range value {
		if char > unicode.MaxASCII {
			v.AddError(field, "This field may only contain ASCII characters")
			return false
		}
	}
	return true
}

// MacAddress checks if the string is a valid MAC address
func (v *Vee) MacAddress(field string, value string) bool {
	_, err := net.ParseMAC(value)
	if err != nil {
		v.AddError(field, "This field must be a valid MAC address")
		return false
	}
	return true
}

// ULID checks if the string is a valid ULID
func (v *Vee) ULID(field string, value string) bool {
	re := regexp.MustCompile("^[0-9A-HJKMNP-TV-Z]{26}$")
	if !re.MatchString(value) {
		v.AddError(field, "This field must be a valid ULID")
		return false
	}
	return true
}

// Distinct checks if all elements in a slice are unique
func (v *Vee) Distinct(field string, values []interface{}) bool {
	seen := make(map[interface{}]bool)
	for _, value := range values {
		if seen[value] {
			v.AddError(field, "This field must contain only unique values")
			return false
		}
		seen[value] = true
	}
	return true
}

// Filled checks if the value is not empty (for strings, slices, maps, and pointers)
func (v *Vee) Filled(field string, value interface{}) bool {
	switch val := value.(type) {
	case string:
		if val == "" {
			v.AddError(field, "This field must be filled")
			return false
		}
	case []interface{}:
		if len(val) == 0 {
			v.AddError(field, "This field must be filled")
			return false
		}
	case map[string]interface{}:
		if len(val) == 0 {
			v.AddError(field, "This field must be filled")
			return false
		}
	case nil:
		v.AddError(field, "This field must be filled")
		return false
	}
	return true
}

// HexColor checks if the string is a valid hexadecimal color code
func (v *Vee) HexColor(field string, value string) bool {
	re := regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")
	if !re.MatchString(value) {
		v.AddError(field, "This field must be a valid hexadecimal color code")
		return false
	}
	return true
}

// ValidationRule is a function type that defines a validation rule
type ValidationRule func(field string, value interface{}, index int) bool

// ForEach applies validation rules to each item in an array
func (v *Vee) ForEach(field string, array interface{}, rules ...ValidationRule) bool {
	// Get the value of the array
	arrayValue := reflect.ValueOf(array)

	// Check if it's a pointer, and if so, get the element it points to
	if arrayValue.Kind() == reflect.Ptr {
		arrayValue = arrayValue.Elem()
	}

	// Ensure we're dealing with a slice or array
	if arrayValue.Kind() != reflect.Slice && arrayValue.Kind() != reflect.Array {
		v.AddError(field, "This field must be an array or slice")
		return false
	}

	// Flag to track if all validations passed
	allValid := true

	// Iterate over each item in the array
	for i := 0; i < arrayValue.Len(); i++ {
		item := arrayValue.Index(i).Interface()
		itemField := fmt.Sprintf("%s.%d", field, i)

		// Apply each validation rule to the item
		for _, rule := range rules {
			if !rule(itemField, item, i) {
				allValid = false
				// Note: We don't return false here so that we can collect all errors
			}
		}
	}

	return allValid
}
