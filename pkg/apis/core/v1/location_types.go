package v1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// anxcloud:object

// Location describes a Anexia site where resources can be deployed.
type Location struct {
	Identifier  string `json:"identifier" anxcloud:"identifier"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	CountryCode string `json:"country"`
	CityCode    string `json:"city_code"`

	Latitude  *float64 `json:"lat"`
	Longitude *float64 `json:"lon"`
}

// UnmarshalJSON handles latitude and longitude given as strings instead of floats by the API.
func (l *Location) UnmarshalJSON(body []byte) error {
	loc := reflect.New(generateAPIType())
	if err := json.Unmarshal(body, loc.Interface()); err != nil {
		return err
	}

	return l.fromAPIType(loc)
}

// MarshalJSON handles latitude and longitude being expected as strings instead of floats by the API.
func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.asAPIType())
}

// generateAPIType generates a new struct type as the API wants the JSON to look, matching the Location type
// but instead of having *float64 types on Latitude and Longitude, the API wants *string.
// We cannot simply make a "apiLocation" struct with Location embedded because we then run into a stack
// overflow - as this function is called to unmarshal the embedded type. Instead we build the struct
// type with reflection. This gives this complicated code here and below, but we only have to add new fields
// to a single struct, making it way less error prone.
func generateAPIType() reflect.Type {
	originalType := reflect.TypeOf(Location{})
	fields := make([]reflect.StructField, originalType.NumField())

	for i := range fields {
		field := originalType.Field(i)
		fields[i] = field

		// for the fields "Latitude" and "Longitude" we change the type to *string
		if field.Name == "Latitude" || field.Name == "Longitude" {
			fields[i].Type = reflect.PtrTo(
				reflect.TypeOf(""),
			)
		}
	}

	return reflect.StructOf(fields)
}

// fromAPIType fills the received Location from the given API type (see generateAPIType above).
func (l *Location) fromAPIType(apiLocation reflect.Value) error {
	apiLocationType := apiLocation.Type().Elem()

	// copy the fields over, parsing Latitude and Longitude from string to float64
	destValue := reflect.ValueOf(l)

	for i := 0; i < apiLocation.Elem().NumField(); i++ {
		field := apiLocationType.Field(i)
		srcField := apiLocation.Elem().Field(i)
		destField := destValue.Elem().Field(i)

		if field.Name == "Latitude" || field.Name == "Longitude" {
			if srcField.IsNil() {
				nilptr := reflect.Zero(
					reflect.PtrTo(
						reflect.TypeOf(float64(42)),
					),
				)
				destField.Set(nilptr)
			} else {
				f, err := strconv.ParseFloat(srcField.Elem().Interface().(string), 64)
				if err != nil {
					return fmt.Errorf("error parsing %v: %w", field.Name, err)
				}

				destField.Set(reflect.ValueOf(&f))

			}
		} else {
			destField.Set(srcField)
		}
	}

	return nil
}

// asAPIType creates a new instance of the struct expected by the API (see generateAPIType above) and fills
// it with data from the received Location.
func (l Location) asAPIType() interface{} {
	apiLocationType := generateAPIType()
	apiLocation := reflect.New(apiLocationType)

	// copy the fields over, formatting  Latitude and Longitude from float64 to string
	srcValue := reflect.ValueOf(l)

	for i := 0; i < apiLocation.Elem().NumField(); i++ {
		field := apiLocationType.Field(i)
		srcField := srcValue.Field(i)
		destField := apiLocation.Elem().Field(i)

		if field.Name == "Latitude" || field.Name == "Longitude" {
			if srcField.IsNil() {
				nilptr := reflect.Zero(
					reflect.PtrTo(
						reflect.TypeOf(""),
					),
				)
				destField.Set(nilptr)
			} else {
				s := strconv.FormatFloat(srcField.Elem().Interface().(float64), 'f', -1, 64)
				destField.Set(reflect.ValueOf(&s))

			}
		} else {
			destField.Set(srcField)
		}
	}

	return apiLocation.Elem().Interface()
}
