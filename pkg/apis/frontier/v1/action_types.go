package v1

type ActionType string

const (
	// ActionTypeURLRewrite configures a rewrite to another URL
	ActionTypeURLRewrite ActionType = "url_rewrite"
	// ActionTypeMockResponse configures the action to return mock data
	ActionTypeMockResponse ActionType = "mock_response"
	// ActionTypeE5EFunction configures the action to synchronously call an e5e function
	ActionTypeE5EFunction ActionType = "e5e_function"
	// ActionTypeE5EAsyncFunction configures the action to asynchronously call an e5e function
	ActionTypeE5EAsyncFunction ActionType = "e5e_async_function"
	// ActionTypeE5EAsyncResult configures the action to fetch an asynchronously e5e function result
	ActionTypeE5EAsyncResult ActionType = "e5e_async_result"
)

// anxcloud:object

// Action represents the lowest entity within Frontier's hierarchy and maps HTTP methods for an endpoint to action handlers.
// Those action handlers may be e5e functions, other HTTP-based APIs or mock responses.
type Action struct {
	omitResponseDecodeOnDestroy
	Identifier         string      `json:"identifier,omitempty" anxcloud:"identifier"`
	EndpointIdentifier string      `json:"endpoint_identifier,omitempty"`
	HTTPRequestMethod  string      `json:"http_request_method,omitempty"`
	Type               ActionType  `json:"type,omitempty"`
	Meta               *ActionMeta `json:"meta,omitempty"`
}

// ActionMeta is used to configure an Action based on its Type
type ActionMeta struct {
	*ActionMetaURLRewrite
	*ActionMetaMockResponse
	*ActionMetaE5EFunction
	*ActionMetaE5EAsyncFunction
	*ActionMetaE5EAsyncResult
}

// ActionMetaURLRewrite is used to configure a resource of type "url_rewrite"
type ActionMetaURLRewrite struct {
	URL string `json:"url_rewrite_url,omitempty"`
}

// ActionMetaMockResponse is used to configure a resource of type "mock_response"
type ActionMetaMockResponse struct {
	Body     string `json:"mock_response_body,omitempty"`
	Language string `json:"mock_response_language,omitempty"`
}

// ActionMetaE5EFunction is used to configure a resource of type "e5e_function"
type ActionMetaE5EFunction struct {
	FunctionIdentifier string `json:"e5e_function_function,omitempty"`
}

// ActionMetaE5EAsyncFunction is used to configure a resource of type "e5e_async_function"
type ActionMetaE5EAsyncFunction struct {
	FunctionIdentifier string `json:"e5e_async_function_function,omitempty"`
}

// ActionMetaE5EAsyncResult is used to configure a resource of type "e5e_async_result"
type ActionMetaE5EAsyncResult struct {
	FunctionIdentifier string `json:"e5e_async_result_function,omitempty"`
}
