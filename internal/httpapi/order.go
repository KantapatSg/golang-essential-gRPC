package httpapi

import (
	"time"

	notificationv1 "github.com/KantapatSg/golang-essential-gRPC/gen/notification"
	orderv1 "github.com/KantapatSg/golang-essential-gRPC/gen/order"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcclient"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type Handler struct {
	orders        *grpcclient.OrderClient
	notifications *grpcclient.NotificationClient
}

func New(app *fiber.App, orders *grpcclient.OrderClient, notifications *grpcclient.NotificationClient) {
	h := &Handler{orders: orders, notifications: notifications}
	app.Get("/healthz", h.health)
	app.Get("/readyz", h.ready)
	g := app.Group("/api/v1/orders")
	g.Post("/", h.create)
	g.Get("/", h.list)
	g.Get("/:id", h.get)
	g.Put("/:id", h.update)
	g.Delete("/:id", h.delete)
	g.Get("/:id/notifications", h.notificationsList)
}

type input struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	Items         []item `json:"items"`
	Status        string `json:"status"`
}
type item struct {
	Name      string  `json:"name"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type orderResponse struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	Items         []item    `json:"items"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type orderListResponse struct {
	Page     int32           `json:"page"`
	PageSize int32           `json:"page_size"`
	Total    int64           `json:"total"`
	Orders   []orderResponse `json:"orders"`
}

type notificationResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	EventType string    `json:"event_type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type notificationListResponse struct {
	Page          int32                  `json:"page"`
	PageSize      int32                  `json:"page_size"`
	Total         int64                  `json:"total"`
	Notifications []notificationResponse `json:"notifications"`
}

func (i input) items() []*orderv1.OrderItem {
	o := make([]*orderv1.OrderItem, len(i.Items))
	for x, v := range i.Items {
		o[x] = &orderv1.OrderItem{Name: v.Name, Quantity: v.Quantity, UnitPrice: v.UnitPrice}
	}
	return o
}
func statusEnum(s string) orderv1.OrderStatus {
	switch s {
	case "", "PENDING":
		return orderv1.OrderStatus_PENDING
	case "CONFIRMED":
		return orderv1.OrderStatus_CONFIRMED
	case "CANCELLED":
		return orderv1.OrderStatus_CANCELLED
	default:
		return orderv1.OrderStatus(-1)
	}
}
func (h *Handler) create(c *fiber.Ctx) error {
	var i input
	if err := c.BodyParser(&i); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
	}
	o, e := h.orders.Create(c.Context(), &orderv1.CreateOrderRequest{CustomerName: i.CustomerName, CustomerEmail: i.CustomerEmail, Items: i.items(), Status: statusEnum(i.Status)})
	if e != nil {
		return grpcError(c, e)
	}
	return c.Status(201).JSON(toOrderResponse(o))
}
func (h *Handler) get(c *fiber.Ctx) error {
	o, e := h.orders.Get(c.Context(), &orderv1.GetOrderRequest{Id: c.Params("id")})
	if e != nil {
		return grpcError(c, e)
	}
	return c.JSON(toOrderResponse(o))
}
func (h *Handler) list(c *fiber.Ctx) error {
	o, e := h.orders.List(c.Context(), &orderv1.ListOrdersRequest{Page: int32(query(c, "page", 1)), PageSize: int32(query(c, "page_size", 20))})
	if e != nil {
		return grpcError(c, e)
	}
	response := orderListResponse{Page: o.GetPage(), PageSize: o.GetPageSize(), Total: o.GetTotal(), Orders: make([]orderResponse, 0, len(o.GetOrders()))}
	for _, order := range o.GetOrders() {
		response.Orders = append(response.Orders, toOrderResponse(order))
	}
	return c.JSON(response)
}
func (h *Handler) update(c *fiber.Ctx) error {
	var i input
	if err := c.BodyParser(&i); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
	}
	o, e := h.orders.Update(c.Context(), &orderv1.UpdateOrderRequest{Id: c.Params("id"), CustomerName: i.CustomerName, CustomerEmail: i.CustomerEmail, Items: i.items(), Status: statusEnum(i.Status)})
	if e != nil {
		return grpcError(c, e)
	}
	return c.JSON(toOrderResponse(o))
}
func (h *Handler) delete(c *fiber.Ctx) error {
	if e := h.orders.Delete(c.Context(), &orderv1.DeleteOrderRequest{Id: c.Params("id")}); e != nil {
		return grpcError(c, e)
	}
	return c.SendStatus(204)
}
func (h *Handler) notificationsList(c *fiber.Ctx) error {
	n, e := h.notifications.List(c.Context(), &notificationv1.ListNotificationsRequest{OrderId: c.Params("id"), Page: int32(query(c, "page", 1)), PageSize: int32(query(c, "page_size", 20))})
	if e != nil {
		return grpcError(c, e)
	}
	response := notificationListResponse{Page: n.GetPage(), PageSize: n.GetPageSize(), Total: n.GetTotal(), Notifications: make([]notificationResponse, 0, len(n.GetNotifications()))}
	for _, notification := range n.GetNotifications() {
		response.Notifications = append(response.Notifications, notificationResponse{
			ID: notification.GetId(), OrderID: notification.GetOrderId(), EventType: notification.GetEventType(), Message: notification.GetMessage(), CreatedAt: notification.GetCreatedAt().AsTime(),
		})
	}
	return c.JSON(response)
}

func toOrderResponse(order *orderv1.Order) orderResponse {
	items := make([]item, 0, len(order.GetItems()))
	for _, orderItem := range order.GetItems() {
		items = append(items, item{Name: orderItem.GetName(), Quantity: orderItem.GetQuantity(), UnitPrice: orderItem.GetUnitPrice()})
	}
	return orderResponse{
		ID: order.GetId(), CustomerName: order.GetCustomerName(), CustomerEmail: order.GetCustomerEmail(), Items: items,
		TotalAmount: order.GetTotalAmount(), Status: order.GetStatus().String(), CreatedAt: order.GetCreatedAt().AsTime(), UpdatedAt: order.GetUpdatedAt().AsTime(),
	}
}
func (h *Handler) health(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) }
func (h *Handler) ready(c *fiber.Ctx) error {
	if e := h.orders.Ready(c.Context()); e != nil {
		return c.Status(503).JSON(fiber.Map{"status": "not_ready", "backend": "order"})
	}
	if e := h.notifications.Ready(c.Context()); e != nil {
		return c.Status(503).JSON(fiber.Map{"status": "not_ready", "backend": "notification"})
	}
	return c.JSON(fiber.Map{"status": "ready"})
}
func query(c *fiber.Ctx, k string, d int) int {
	n, e := strconv.Atoi(c.Query(k))
	if e != nil || n < 1 {
		return d
	}
	return n
}
func grpcError(c *fiber.Ctx, e error) error {
	st, ok := status.FromError(e)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return c.Status(400).JSON(fiber.Map{"error": st.Message()})
	case codes.NotFound:
		return c.Status(404).JSON(fiber.Map{"error": st.Message()})
	case codes.Unavailable, codes.DeadlineExceeded:
		return c.Status(503).JSON(fiber.Map{"error": "backend service unavailable"})
	default:
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
}
