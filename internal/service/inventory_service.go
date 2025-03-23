package service

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"net/http"
)

type InventoryService struct {
}

func (s *InventoryService) GetInventoryRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(baseResponse("Inventory not found", http.StatusNotFound))
}

func (s *InventoryService) PlaceOrderRequest(requestBody gen.Order, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *InventoryService) DeleteOrderRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(baseResponse("Order deleted", http.StatusOK))
}

func (s *InventoryService) GetOrderByIdRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(baseResponse("Order not found", http.StatusNotFound))
}
