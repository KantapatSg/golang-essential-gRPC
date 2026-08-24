package grpcclient

import (
	"context"
	notificationv1 "github.com/KantapatSg/golang-essential-gRPC/gen/notification"
	orderv1 "github.com/KantapatSg/golang-essential-gRPC/gen/order"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	client orderv1.OrderServiceClient
}

func NewOrderClient(target string) (*OrderClient, error) {
	c, e := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if e != nil {
		return nil, e
	}
	return &OrderClient{conn: c, client: orderv1.NewOrderServiceClient(c)}, nil
}
func (c *OrderClient) Close() error { return c.conn.Close() }
func (c *OrderClient) Create(ctx context.Context, r *orderv1.CreateOrderRequest) (*orderv1.Order, error) {
	x, cancel := timeout(ctx)
	defer cancel()
	return c.client.CreateOrder(x, r)
}
func (c *OrderClient) Get(ctx context.Context, r *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	x, cancel := timeout(ctx)
	defer cancel()
	return c.client.GetOrder(x, r)
}
func (c *OrderClient) List(ctx context.Context, r *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	x, cancel := timeout(ctx)
	defer cancel()
	return c.client.ListOrders(x, r)
}
func (c *OrderClient) Update(ctx context.Context, r *orderv1.UpdateOrderRequest) (*orderv1.Order, error) {
	x, cancel := timeout(ctx)
	defer cancel()
	return c.client.UpdateOrder(x, r)
}
func (c *OrderClient) Delete(ctx context.Context, r *orderv1.DeleteOrderRequest) error {
	x, cancel := timeout(ctx)
	defer cancel()
	_, e := c.client.DeleteOrder(x, r)
	return e
}
func (c *OrderClient) Ready(ctx context.Context) error {
	x, cancel := timeout(ctx)
	defer cancel()
	_, e := c.client.ListOrders(x, &orderv1.ListOrdersRequest{Page: 1, PageSize: 1})
	return e
}

type NotificationClient struct {
	conn   *grpc.ClientConn
	client notificationv1.NotificationServiceClient
}

func NewNotificationClient(target string) (*NotificationClient, error) {
	c, e := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if e != nil {
		return nil, e
	}
	return &NotificationClient{conn: c, client: notificationv1.NewNotificationServiceClient(c)}, nil
}
func (c *NotificationClient) Close() error { return c.conn.Close() }
func (c *NotificationClient) Send(ctx context.Context, id uuid.UUID, event, message string) error {
	x, cancel := timeout(ctx)
	defer cancel()
	_, e := c.client.SendNotification(x, &notificationv1.SendNotificationRequest{OrderId: id.String(), EventType: event, Message: message})
	return e
}
func (c *NotificationClient) List(ctx context.Context, r *notificationv1.ListNotificationsRequest) (*notificationv1.ListNotificationsResponse, error) {
	x, cancel := timeout(ctx)
	defer cancel()
	return c.client.ListNotifications(x, r)
}
func (c *NotificationClient) Ready(ctx context.Context) error {
	x, cancel := timeout(ctx)
	defer cancel()
	_, e := c.client.ListNotifications(x, &notificationv1.ListNotificationsRequest{Page: 1, PageSize: 1})
	return e
}
func timeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}
