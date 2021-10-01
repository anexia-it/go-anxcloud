package param

import "net/url"

type Parameter func(values url.Values)

func ParameterBuilder(filterKey string) func(string) Parameter {
	return func(filterValue string) Parameter {
		return func(values url.Values) {
			values.Set(filterKey, filterValue)
		}
	}
}
