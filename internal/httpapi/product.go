package httpapi

import (
	"strconv"

	productv1 "github.com/KantapatSg/golang-essential-gRPC/gen/product"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcclient"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(app *fiber.App, client *grpcclient.ProductClient) {
	h := &handler{client: client}
	app.Get("/healthz", h.health)
	app.Get("/readyz", h.ready)
	group := app.Group("/api/v1/products")
	group.Post("/", h.create)
	group.Get("/", h.list)
	group.Get("/:id", h.get)
	group.Put("/:id", h.update)
	group.Delete("/:id", h.delete)
}

type handler struct{ client *grpcclient.ProductClient }

type productInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
}

func (h *handler) create(c *fiber.Ctx) error {
	var input productInput
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.ErrBadRequest)
	}
	p, err := h.client.Create(c.Context(), &productv1.CreateProductRequest{Name: input.Name, Description: input.Description, Price: input.Price, Stock: input.Stock})
	if err != nil {
		return grpcError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}
func (h *handler) get(c *fiber.Ctx) error {
	p, err := h.client.Get(c.Context(), &productv1.GetProductRequest{Id: c.Params("id")})
	if err != nil {
		return grpcError(c, err)
	}
	return c.JSON(p)
}
func (h *handler) list(c *fiber.Ctx) error {
	page := parseQueryInt(c.Query("page"), 1)
	size := parseQueryInt(c.Query("page_size"), 20)
	result, err := h.client.List(c.Context(), &productv1.ListProductsRequest{Page: int32(page), PageSize: int32(size)})
	if err != nil {
		return grpcError(c, err)
	}
	return c.JSON(result)
}
func (h *handler) update(c *fiber.Ctx) error {
	var input productInput
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.ErrBadRequest)
	}
	p, err := h.client.Update(c.Context(), &productv1.UpdateProductRequest{Id: c.Params("id"), Name: input.Name, Description: input.Description, Price: input.Price, Stock: input.Stock})
	if err != nil {
		return grpcError(c, err)
	}
	return c.JSON(p)
}
func (h *handler) delete(c *fiber.Ctx) error {
	if err := h.client.Delete(c.Context(), &productv1.DeleteProductRequest{Id: c.Params("id")}); err != nil {
		return grpcError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *handler) health(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) }
func (h *handler) ready(c *fiber.Ctx) error {
	if err := h.client.Ready(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "not_ready"})
	}
	return c.JSON(fiber.Map{"status": "ready"})
}

func parseQueryInt(value string, fallback int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}
func writeError(c *fiber.Ctx, err *fiber.Error) error {
	return c.Status(err.Code).JSON(fiber.Map{"error": err.Message})
}
func grpcError(c *fiber.Ctx, err error) error {
	code, ok := status.FromError(err)
	if !ok {
		return writeError(c, fiber.ErrInternalServerError)
	}
	switch code.Code() {
	case codes.InvalidArgument:
		return c.Status(400).JSON(fiber.Map{"error": code.Message()})
	case codes.NotFound:
		return c.Status(404).JSON(fiber.Map{"error": code.Message()})
	case codes.Unavailable, codes.DeadlineExceeded:
		return c.Status(503).JSON(fiber.Map{"error": "product service unavailable"})
	default:
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
}
