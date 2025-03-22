// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gen

import (
	"time"

	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/server"
)

const (
	Api_keyScopes       = "api_key.Scopes"
	Petstore_authScopes = "petstore_auth.Scopes"
)

// Defines values for OrderStatus.
const (
	Approved  OrderStatus = "approved"
	Delivered OrderStatus = "delivered"
	Placed    OrderStatus = "placed"
)

// Defines values for PetStatus.
const (
	PetStatusAvailable PetStatus = "available"
	PetStatusPending   PetStatus = "pending"
	PetStatusSold      PetStatus = "sold"
)

// Defines values for FindPetsByStatusParamsStatus.
const (
	FindPetsByStatusParamsStatusAvailable FindPetsByStatusParamsStatus = "available"
	FindPetsByStatusParamsStatusPending   FindPetsByStatusParamsStatus = "pending"
	FindPetsByStatusParamsStatusSold      FindPetsByStatusParamsStatus = "sold"
)

// ApiResponse defines model for ApiResponse.
type ApiResponse struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// BaseResponse defines model for BaseResponse.
type BaseResponse struct {
	Code      *int    `json:"code,omitempty"`
	Message   *string `json:"message,omitempty"`
	Success   *bool   `json:"success,omitempty"`
	Timestamp *int    `json:"timestamp,omitempty"`
}

// Category defines model for Category.
type Category struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Expires *int    `json:"expires,omitempty"`
	Token   *string `json:"token,omitempty"`
}

// Order defines model for Order.
type Order struct {
	Complete *bool      `json:"complete,omitempty"`
	Id       *int64     `json:"id,omitempty"`
	PetId    *int64     `json:"petId,omitempty"`
	Quantity *int32     `json:"quantity,omitempty"`
	ShipDate *time.Time `json:"shipDate,omitempty"`

	// Status Order Status
	Status *OrderStatus `json:"status,omitempty"`
}

// OrderStatus Order Status
type OrderStatus string

// Pet defines model for Pet.
type Pet struct {
	Category  *Category `json:"category,omitempty"`
	Id        *int64    `json:"id,omitempty"`
	Name      string    `json:"name"`
	PhotoUrls []string  `json:"photoUrls"`

	// Status pet status in the store
	Status *PetStatus `json:"status,omitempty"`
	Tags   *[]Tag     `json:"tags,omitempty"`
}

// PetStatus pet status in the store
type PetStatus string

// Tag defines model for Tag.
type Tag struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// User defines model for User.
type User struct {
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	Id        *int64  `json:"id,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Password  *string `json:"password,omitempty"`
	Phone     *string `json:"phone,omitempty"`

	// UserStatus User Status
	UserStatus *int32  `json:"userStatus,omitempty"`
	Username   *string `json:"username,omitempty"`
}

// Authorization defines model for authorization.
type Authorization = string

// AddPetParams defines parameters for AddPet.
type AddPetParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// UpdatePetParams defines parameters for UpdatePet.
type UpdatePetParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// FindPetsByStatusParams defines parameters for FindPetsByStatus.
type FindPetsByStatusParams struct {
	// Status Status values that need to be considered for filter
	Status *FindPetsByStatusParamsStatus `form:"status,omitempty" json:"status,omitempty"`

	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// FindPetsByStatusParamsStatus defines parameters for FindPetsByStatus.
type FindPetsByStatusParamsStatus string

// FindPetsByTagsParams defines parameters for FindPetsByTags.
type FindPetsByTagsParams struct {
	// Tags Tags to filter by
	Tags *[]string `form:"tags,omitempty" json:"tags,omitempty"`

	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// DeletePetParams defines parameters for DeletePet.
type DeletePetParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
	ApiKey        *string       `json:"api_key,omitempty"`
}

// GetPetByIdParams defines parameters for GetPetById.
type GetPetByIdParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// UpdatePetWithFormParams defines parameters for UpdatePetWithForm.
type UpdatePetWithFormParams struct {
	// Name Name of pet that needs to be updated
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Status Status of pet that needs to be updated
	Status *string `form:"status,omitempty" json:"status,omitempty"`

	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// UploadFileParams defines parameters for UploadFile.
type UploadFileParams struct {
	// AdditionalMetadata Additional Metadata
	AdditionalMetadata *string `form:"additionalMetadata,omitempty" json:"additionalMetadata,omitempty"`

	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// PlaceOrderParams defines parameters for PlaceOrder.
type PlaceOrderParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// DeleteOrderParams defines parameters for DeleteOrder.
type DeleteOrderParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// GetOrderByIdParams defines parameters for GetOrderById.
type GetOrderByIdParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// CreateUsersWithListInputJSONBody defines parameters for CreateUsersWithListInput.
type CreateUsersWithListInputJSONBody = []User

// CreateUsersWithListInputParams defines parameters for CreateUsersWithListInput.
type CreateUsersWithListInputParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// LogoutUserParams defines parameters for LogoutUser.
type LogoutUserParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// DeleteUserParams defines parameters for DeleteUser.
type DeleteUserParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// GetUserByNameParams defines parameters for GetUserByName.
type GetUserByNameParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

// UpdateUserParams defines parameters for UpdateUser.
type UpdateUserParams struct {
	// Authorization jwt
	Authorization Authorization `json:"Authorization"`
}

func (rh *requestHandlerImpl[T]) PostLogin(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := LoginRequest{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body LoginRequest) {
		rh.Service.PostLoginRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) AddPet(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := Pet{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body Pet) {
		rh.Service.AddPetRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) UpdatePet(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := Pet{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body Pet) {
		rh.Service.UpdatePetRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) FindPetsByStatus(ctx *apicontext.ApiRequestContext[T]) {

	//request := FindPetsByStatusRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.FindPetsByStatusRequest(ctx)

}

func (rh *requestHandlerImpl[T]) FindPetsByTags(ctx *apicontext.ApiRequestContext[T]) {

	//request := FindPetsByTagsRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.FindPetsByTagsRequest(ctx)

}

func (rh *requestHandlerImpl[T]) DeletePet(ctx *apicontext.ApiRequestContext[T]) {

	//request := DeletePetRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.DeletePetRequest(ctx)

}

func (rh *requestHandlerImpl[T]) GetPetById(ctx *apicontext.ApiRequestContext[T]) {

	//request := GetPetByIdRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.GetPetByIdRequest(ctx)

}

func (rh *requestHandlerImpl[T]) UpdatePetWithForm(ctx *apicontext.ApiRequestContext[T]) {

	//request := UpdatePetWithFormRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.UpdatePetWithFormRequest(ctx)

}

func (rh *requestHandlerImpl[T]) UploadFile(ctx *apicontext.ApiRequestContext[T]) {

	//request := UploadFileRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.UploadFileRequest(ctx)

}

func (rh *requestHandlerImpl[T]) GetInventory(ctx *apicontext.ApiRequestContext[T]) {

	//request := GetInventoryRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.GetInventoryRequest(ctx)

}

func (rh *requestHandlerImpl[T]) PlaceOrder(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := Order{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body Order) {
		rh.Service.PlaceOrderRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) DeleteOrder(ctx *apicontext.ApiRequestContext[T]) {

	//request := DeleteOrderRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.DeleteOrderRequest(ctx)

}

func (rh *requestHandlerImpl[T]) GetOrderById(ctx *apicontext.ApiRequestContext[T]) {

	//request := GetOrderByIdRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.GetOrderByIdRequest(ctx)

}

func (rh *requestHandlerImpl[T]) CreateUser(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := User{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body User) {
		rh.Service.CreateUserRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) CreateUsersWithListInput(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := CreateUsersWithListInputJSONBody{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body CreateUsersWithListInputJSONBody) {
		rh.Service.CreateUsersWithListInputRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

func (rh *requestHandlerImpl[T]) LogoutUser(ctx *apicontext.ApiRequestContext[T]) {

	//request := LogoutUserRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.LogoutUserRequest(ctx)

}

func (rh *requestHandlerImpl[T]) DeleteUser(ctx *apicontext.ApiRequestContext[T]) {

	//request := DeleteUserRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.DeleteUserRequest(ctx)

}

func (rh *requestHandlerImpl[T]) GetUserByName(ctx *apicontext.ApiRequestContext[T]) {

	//request := GetUserByNameRequestParams{}
	// server.PopulateFieldsFromRequest(ctx, &request)
	rh.Service.GetUserByNameRequest(ctx)

}

func (rh *requestHandlerImpl[T]) UpdateUser(ctx *apicontext.ApiRequestContext[T]) {

	requestBody := User{}
	server.GetRequestBody(ctx, requestBody, func(ctx *apicontext.ApiRequestContext[T], body User) {
		rh.Service.UpdateUserRequest(body, ctx)
	}, func(ctx *apicontext.ApiRequestContext[T], err error) {
		ctx.InternalServerError("Internal server error")
	})

}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginRequest = LoginRequest

// AddPetJSONRequestBody defines body for AddPet for application/json ContentType.
type AddPetRequest = Pet

// UpdatePetJSONRequestBody defines body for UpdatePet for application/json ContentType.
type UpdatePetRequest = Pet

// PlaceOrderJSONRequestBody defines body for PlaceOrder for application/json ContentType.
type PlaceOrderRequest = Order

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserRequest = User

// CreateUsersWithListInputJSONRequestBody defines body for CreateUsersWithListInput for application/json ContentType.
type CreateUsersWithListInputRequest = CreateUsersWithListInputJSONBody

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserRequest = User

// RequestHandler represents all server handlers.
type RequestHandler[T apicontext.ApiPrincipalContext] interface {
	// Authentication endpoint
	// (POST /login)
	PostLogin(ctx *apicontext.ApiRequestContext[T])
	// Add a new pet to the store
	// (POST /pet)
	AddPet(ctx *apicontext.ApiRequestContext[T])
	// Update an existing pet
	// (PUT /pet)
	UpdatePet(ctx *apicontext.ApiRequestContext[T])
	// Finds Pets by status
	// (GET /pet/findByStatus)
	FindPetsByStatus(ctx *apicontext.ApiRequestContext[T])
	// Finds Pets by tags
	// (GET /pet/findByTags)
	FindPetsByTags(ctx *apicontext.ApiRequestContext[T])
	// Deletes a pet
	// (DELETE /pet/{petId})
	DeletePet(ctx *apicontext.ApiRequestContext[T])
	// Find pet by ID
	// (GET /pet/{petId})
	GetPetById(ctx *apicontext.ApiRequestContext[T])
	// Updates a pet in the store with form data
	// (POST /pet/{petId})
	UpdatePetWithForm(ctx *apicontext.ApiRequestContext[T])
	// uploads an image
	// (POST /pet/{petId}/uploadImage)
	UploadFile(ctx *apicontext.ApiRequestContext[T])
	// Returns pet inventories by status
	// (GET /store/inventory)
	GetInventory(ctx *apicontext.ApiRequestContext[T])
	// Place an order for a pet
	// (POST /store/order)
	PlaceOrder(ctx *apicontext.ApiRequestContext[T])
	// Delete purchase order by ID
	// (DELETE /store/order/{orderId})
	DeleteOrder(ctx *apicontext.ApiRequestContext[T])
	// Find purchase order by ID
	// (GET /store/order/{orderId})
	GetOrderById(ctx *apicontext.ApiRequestContext[T])
	// Create user
	// (POST /user)
	CreateUser(ctx *apicontext.ApiRequestContext[T])
	// Creates list of users with given input array
	// (POST /user/createWithList)
	CreateUsersWithListInput(ctx *apicontext.ApiRequestContext[T])
	// Logs out current logged in user session
	// (GET /user/logout)
	LogoutUser(ctx *apicontext.ApiRequestContext[T])
	// Delete user
	// (DELETE /user/{username})
	DeleteUser(ctx *apicontext.ApiRequestContext[T])
	// Get user by user name
	// (GET /user/{username})
	GetUserByName(ctx *apicontext.ApiRequestContext[T])
	// Update user
	// (PUT /user/{username})
	UpdateUser(ctx *apicontext.ApiRequestContext[T])
}

type ServiceRequestHandler[T apicontext.ApiPrincipalContext] interface {

	// PostLoginRequest(requestBody LoginRequest, requestParams PostLoginRequestParams, ctx *context.ApiRequestContext[T])
	PostLoginRequest(requestBody LoginRequest, ctx *apicontext.ApiRequestContext[T])

	// AddPetRequest(requestBody Pet, requestParams AddPetRequestParams, ctx *context.ApiRequestContext[T])
	AddPetRequest(requestBody Pet, ctx *apicontext.ApiRequestContext[T])

	// UpdatePetRequest(requestBody Pet, requestParams UpdatePetRequestParams, ctx *context.ApiRequestContext[T])
	UpdatePetRequest(requestBody Pet, ctx *apicontext.ApiRequestContext[T])

	FindPetsByStatusRequest(ctx *apicontext.ApiRequestContext[T])

	FindPetsByTagsRequest(ctx *apicontext.ApiRequestContext[T])

	DeletePetRequest(ctx *apicontext.ApiRequestContext[T])

	GetPetByIdRequest(ctx *apicontext.ApiRequestContext[T])

	UpdatePetWithFormRequest(ctx *apicontext.ApiRequestContext[T])

	UploadFileRequest(ctx *apicontext.ApiRequestContext[T])

	GetInventoryRequest(ctx *apicontext.ApiRequestContext[T])

	// PlaceOrderRequest(requestBody Order, requestParams PlaceOrderRequestParams, ctx *context.ApiRequestContext[T])
	PlaceOrderRequest(requestBody Order, ctx *apicontext.ApiRequestContext[T])

	DeleteOrderRequest(ctx *apicontext.ApiRequestContext[T])

	GetOrderByIdRequest(ctx *apicontext.ApiRequestContext[T])

	// CreateUserRequest(requestBody User, requestParams CreateUserRequestParams, ctx *context.ApiRequestContext[T])
	CreateUserRequest(requestBody User, ctx *apicontext.ApiRequestContext[T])

	// CreateUsersWithListInputRequest(requestBody CreateUsersWithListInputJSONBody, requestParams CreateUsersWithListInputRequestParams, ctx *context.ApiRequestContext[T])
	CreateUsersWithListInputRequest(requestBody CreateUsersWithListInputJSONBody, ctx *apicontext.ApiRequestContext[T])

	LogoutUserRequest(ctx *apicontext.ApiRequestContext[T])

	DeleteUserRequest(ctx *apicontext.ApiRequestContext[T])

	GetUserByNameRequest(ctx *apicontext.ApiRequestContext[T])

	// UpdateUserRequest(requestBody User, requestParams UpdateUserRequestParams, ctx *context.ApiRequestContext[T])
	UpdateUserRequest(requestBody User, ctx *apicontext.ApiRequestContext[T])
}

type requestHandlerImpl[T apicontext.ApiPrincipalContext] struct {
	Service ServiceRequestHandler[T]
}

// ResourcesHandler registers API endpoints from generated code.

// - RequestHandler.PostLogin
// - RequestHandler.AddPet
// - RequestHandler.UpdatePet
// - RequestHandler.FindPetsByStatus
// - RequestHandler.FindPetsByTags
// - RequestHandler.DeletePet
// - RequestHandler.GetPetById
// - RequestHandler.UpdatePetWithForm
// - RequestHandler.UploadFile
// - RequestHandler.GetInventory
// - RequestHandler.PlaceOrder
// - RequestHandler.DeleteOrder
// - RequestHandler.GetOrderById
// - RequestHandler.CreateUser
// - RequestHandler.CreateUsersWithListInput
// - RequestHandler.LogoutUser
// - RequestHandler.DeleteUser
// - RequestHandler.GetUserByName
// - RequestHandler.UpdateUser
// Parameters:
//   - apiServer: The API router handler used for setting up routes and middleware.
//   - server: The server interface implementation containing the endpoint handlers.
//
// Generics:
//   - T: A type that satisfies the context.ApiPrincipalContext interface, representing the principal/context
//     involved in the API operations.
//
// This function will use the RequestHandler implementation
// that has already been generated to bind specific API routes
// dynamically at runtime, based on the provided security definitions
// and endpoint configurations.
func ResourcesHandler[T apicontext.ApiPrincipalContext](apiServer server.ApiRouterHandler[T], service ServiceRequestHandler[T]) {
	handler := &requestHandlerImpl[T]{
		Service: service,
	}
	ApiResourceRegister(apiServer, handler)
}

// ApiResourceRegister is a customizable resource handler that registers API endpoints from generated code.
// This method binds the custom `RequestHandler` implementation to specific API routes,
// allowing dynamic configuration of handlers.
//
// Parameters:
//   - apiServer: The API router handler used for setting up routes and middleware.
//   - handler: The `RequestHandler` interface implementation containing the actual endpoint handlers.
//
// Generics:
//   - T: A type that satisfies the context.ApiPrincipalContext interface, representing the principal/context
//     involved in the API operations.
func ApiResourceRegister[T apicontext.ApiPrincipalContext](apiServer server.ApiRouterHandler[T], handler RequestHandler[T]) {
	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.PostLogin, "/login", "POST")

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.AddPet, "/pet", "POST", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.UpdatePet, "/pet", "PUT", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.FindPetsByStatus, "/pet/findByStatus", "GET", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.FindPetsByTags, "/pet/findByTags", "GET", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.DeletePet, "/pet/{petId}", "DELETE", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.GetPetById, "/pet/{petId}", "GET", []string{"api_key", "write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.UpdatePetWithForm, "/pet/{petId}", "POST", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.// Initialize $scopes if it's empty.// Append with a comma if $scopes is not empty.
	apiServer.Add(handler.UploadFile, "/pet/{petId}/uploadImage", "POST", []string{"write:pets", "read:pets"}...)

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.GetInventory, "/store/inventory", "GET")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.PlaceOrder, "/store/order", "POST")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.DeleteOrder, "/store/order/{orderId}", "DELETE")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.GetOrderById, "/store/order/{orderId}", "GET")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.CreateUser, "/user", "POST")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.CreateUsersWithListInput, "/user/createWithList", "POST")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.LogoutUser, "/user/logout", "GET")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.DeleteUser, "/user/{username}", "DELETE")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.GetUserByName, "/user/{username}", "GET")

	// Initialize an empty string for the merged scopes.
	apiServer.PublicRouter(handler.UpdateUser, "/user/{username}", "PUT")

}

func ApiResourceHandler[T apicontext.ApiPrincipalContext](service ServiceRequestHandler[T]) func(handler server.ApiRouterHandler[T]) {
	return func(handler server.ApiRouterHandler[T]) {
		ResourcesHandler(handler, service)
	}
}
