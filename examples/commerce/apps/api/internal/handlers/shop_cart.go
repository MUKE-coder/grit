package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// The storefront's cart endpoints.
//
// Hand-written, because Grit's --public generator emits read-only endpoints on
// purpose: an allowlisted list, a get-by-slug and a related strip, for an
// audience that is the internet. A cart is the other thing a storefront needs,
// a public endpoint that writes, and there is no generator for one yet.
//
// Three decisions worth knowing about, because they are the ones that make a
// public write surface safe rather than merely working:
//
//  1. The cart is identified by a token the client sends and never by anything
//     the client could guess. The token is the authorisation, since there is no
//     user to check it against.
//
//  2. Every response says Cache-Control: private. These endpoints are mounted
//     outside Grit's cached public group for the same reason, but saying it in
//     the response means the answer stays correct if somebody moves the route:
//     the cache middleware refuses to store a response that declares itself
//     per-caller.
//
//  3. Nothing about money comes from the request. The client sends a product
//     id, maybe a variant id and a quantity. Every price is read from the
//     database, which is the difference between a shop and a suggestion.

// ShopCartHandler serves the cart.
type ShopCartHandler struct {
	DB    *gorm.DB
	Carts *services.ShopCartService
}

func NewShopCartHandler(db *gorm.DB) *ShopCartHandler {
	return &ShopCartHandler{DB: db, Carts: services.NewShopCartService(db)}
}

// cartTokenHeader is where the client puts the token.
//
// A header rather than a cookie, because the caller is the storefront's own
// server: the browser talks to Next.js, Next.js keeps the cookie and talks to
// this API. That keeps the token out of reach of any script on the shop's
// pages, and means this endpoint needs no CORS allowance and no CSRF token,
// because no browser ever calls it directly.
const cartTokenHeader = "X-Cart-Token"

// private marks a response as one caller's own.
func private(c *gin.Context) {
	c.Header("Cache-Control", "private, no-store")
}

// cart resolves the request's cart, creating one when there is no usable token.
func (h *ShopCartHandler) cart(c *gin.Context) (*models.Cart, bool) {
	cart, err := h.Carts.FindOrCreate(c.Request.Context(), c.GetHeader(cartTokenHeader))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not open a cart")
		return nil, false
	}
	return cart, true
}

// respondWithCart writes the whole cart, which every endpoint here returns.
//
// All four return the same shape, so the client has one thing to parse and one
// place to put it: a mutation that answered only "ok" would need a second
// request to redraw, on the interaction users repeat most.
func (h *ShopCartHandler) respondWithCart(c *gin.Context, cart *models.Cart, status int) {
	view, err := h.Carts.Get(c.Request.Context(), cart)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not read the cart")
		return
	}
	private(c)
	c.JSON(status, gin.H{"data": view})
}

// Get handles GET /api/v1/shop/cart.
func (h *ShopCartHandler) Get(c *gin.Context) {
	cart, ok := h.cart(c)
	if !ok {
		return
	}
	h.respondWithCart(c, cart, http.StatusOK)
}

// Add handles POST /api/v1/shop/cart/items.
func (h *ShopCartHandler) Add(c *gin.Context) {
	var body struct {
		ProductID string `json:"product_id" binding:"required"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, "A product id is required")
		return
	}

	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.Add(c.Request.Context(), cart, body.ProductID, body.VariantID, body.Quantity)
	switch {
	case errors.Is(err, services.ErrNotForSale):
		// 409 rather than 404: the thing exists, it just cannot be bought
		// right now, and the storefront shows those differently.
		private(c)
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code": "NOT_FOR_SALE", "message": "That is not available at the moment",
		}})
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not add that to the cart")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// SetQuantity handles PATCH /api/v1/shop/cart/items/:id.
func (h *ShopCartHandler) SetQuantity(c *gin.Context) {
	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, "A quantity is required")
		return
	}

	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.SetQuantity(c.Request.Context(), cart, c.Param("id"), body.Quantity)
	switch {
	case errors.Is(err, services.ErrNotInCart):
		h.notInCart(c)
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not change the quantity")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// Remove handles DELETE /api/v1/shop/cart/items/:id.
func (h *ShopCartHandler) Remove(c *gin.Context) {
	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.Remove(c.Request.Context(), cart, c.Param("id"))
	switch {
	case errors.Is(err, services.ErrNotInCart):
		h.notInCart(c)
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not remove that line")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// notInCart answers a line id this cart does not hold.
//
// The same 404 whether the line never existed or belongs to somebody else's
// cart. Which of the two it is would tell a caller holding a guessed id that
// they had guessed correctly.
func (h *ShopCartHandler) notInCart(c *gin.Context) {
	private(c)
	c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
		"code": "NOT_FOUND", "message": "No such line in this cart",
	}})
}
