package v2

// commonRequestBody is used to optionally omit the `State` field on create and update
// by embedding it to the request in the FilterAPIRequestBody hook
type commonRequestBody struct {
	State string `json:"state,omitempty"`
}
