package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
)

// The shopping cart.
//
// Hand-written. `grit generate resource Cart` and `grit generate resource
// CartItem` produced the models, the migrations, the admin screens and the
// staff-only CRUD, which is most of the work and none of the shop. What is
// here is the part a shop has to decide for itself:
//
//   - a cart belongs to a browser rather than to a user, because a shopper
//     fills one before deciding whether to have an account
//   - the price of a line comes from the database at the moment it is added,
//     never from the request, because a price in a request is a price the
//     customer chose
//   - adding the same variant twice raises the quantity rather than making a
//     second line, which is what every shopper expects and nobody asks for
//
// Grit's --public generator emits read-only endpoints (list, get, related) on
// purpose: the audience is the internet. A cart is the other thing a storefront
// needs, a public endpoint that writes, and there is no generator for it. This
// file is what one looks like.

// ErrNotInCart is returned when a line id does not belong to the cart that
// asked, which is the same answer as a line that never existed: whether
// somebody else's cart holds it is not a question this API answers.
var ErrNotInCart = errors.New("no such line in this cart")

// ErrNotForSale covers a product that is archived, switched off or sold out,
// and a variant that is out of stock. One error, because the shopper is told
// the same thing either way and the difference is not theirs to act on.
var ErrNotForSale = errors.New("that is not for sale")

// ShopCartService holds the storefront's cart operations.
//
// Not CartService: `grit generate resource Cart` already generated one of
// those, with the staff CRUD the admin screens call. This is the shopper's
// side of the same tables, and the two are deliberately separate services
// rather than one with two audiences.
type ShopCartService struct {
	DB *gorm.DB
}

func NewShopCartService(db *gorm.DB) *ShopCartService {
	return &ShopCartService{DB: db}
}

func (s *ShopCartService) shopDB(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// NewToken mints the opaque string that identifies a cart to a browser.
//
// 32 bytes from crypto/rand, so it cannot be guessed: the token IS the
// authorisation to read and change that cart, since there is no user to check
// it against. A sequential id or a hash of something known would let anyone
// read the cart next door.
func NewToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating a cart token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// Cart is a cart and its lines, priced.
type Cart struct {
	Token    string     `json:"token"`
	Currency string     `json:"currency"`
	Lines    []CartLine `json:"lines"`
	// Count is the number of items, not the number of lines: three of one
	// thing is three, which is what the badge on the header shows.
	Count    int         `json:"count"`
	Subtotal money.Money `json:"subtotal"`
}

// CartLine is one line, with the total it contributes.
type CartLine struct {
	ID           string      `json:"id"`
	ProductID    string      `json:"product_id"`
	ProductTitle string      `json:"product_title"`
	ProductPath  string      `json:"product_path"`
	VariantID    string      `json:"variant_id"`
	VariantLabel string      `json:"variant_label"`
	ImageURL     string      `json:"image_url"`
	Quantity     int         `json:"quantity"`
	UnitPrice    money.Money `json:"unit_price"`
	LineTotal    money.Money `json:"line_total"`
}

// FindOrCreate returns the cart for a token, creating one when the token is
// empty or names a cart that no longer exists.
//
// A token that does not resolve gets a new cart rather than an error. A cart is
// not an account: it expires, it gets cleaned up, and a shopper returning to a
// stale one should be able to start shopping rather than be told something
// broke.
func (s *ShopCartService) FindOrCreate(ctx context.Context, token string) (*models.Cart, error) {
	if token != "" {
		var cart models.Cart
		err := s.shopDB(ctx).Where("token = ?", token).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("loading the cart: %w", err)
		}
	}

	fresh, err := NewToken()
	if err != nil {
		return nil, err
	}
	cart := models.Cart{Token: fresh, Currency: "USD"}
	if err := s.shopDB(ctx).Create(&cart).Error; err != nil {
		return nil, fmt.Errorf("creating a cart: %w", err)
	}
	return &cart, nil
}

// Get returns the cart's lines with their totals.
func (s *ShopCartService) Get(ctx context.Context, cart *models.Cart) (*Cart, error) {
	var items []models.CartItem
	if err := s.shopDB(ctx).
		Where("cart_id = ?", cart.ID).
		Order("created_at asc").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("loading the cart lines: %w", err)
	}

	// The handles, in one query rather than one per line.
	handles := map[string]string{}
	if len(items) > 0 {
		ids := make([]string, 0, len(items))
		for _, item := range items {
			ids = append(ids, item.ProductID)
		}
		var products []models.Product
		if err := s.shopDB(ctx).Select("id", "handle").Where("id IN ?", ids).Find(&products).Error; err != nil {
			return nil, fmt.Errorf("loading the products in the cart: %w", err)
		}
		for _, p := range products {
			handles[p.ID] = p.Handle
		}
	}

	out := &Cart{
		Token:    cart.Token,
		Currency: cart.Currency,
		Lines:    make([]CartLine, 0, len(items)),
		Subtotal: money.New(0, cart.Currency),
	}
	for _, item := range items {
		lineTotal := item.UnitPrice.MulInt(int64(item.Quantity))
		out.Lines = append(out.Lines, CartLine{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductTitle: item.Title,
			ProductPath:  "/product/" + handles[item.ProductID],
			VariantID:    item.VariantID,
			VariantLabel: item.VariantLabel,
			ImageURL:     item.ImageURL,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			LineTotal:    lineTotal,
		})
		out.Count += item.Quantity
		// Add returns an error rather than panicking when two currencies are
		// mixed. A cart holds one currency, so this cannot fire today; it is
		// checked anyway, because the day the shop sells in two it would
		// otherwise produce a silently wrong subtotal.
		subtotal, err := out.Subtotal.Add(lineTotal)
		if err != nil {
			return nil, fmt.Errorf("totalling the cart: %w", err)
		}
		out.Subtotal = subtotal
	}
	return out, nil
}

// Add puts a quantity of one variant, or of a plain product, into the cart.
//
// variantID may be empty, for a product with no options: most of a real shop
// is mugs and tote bags, and a storefront should not have to invent a variant
// to buy one.
//
// The price is read here and stored on the line. That is a snapshot on purpose:
// a shopper who put something in the basket at $28 is charged $28 at checkout
// even if the shop raises the price while they browse, and the alternative is a
// total that changes between the cart and the card.
func (s *ShopCartService) Add(ctx context.Context, cart *models.Cart, productID, variantID string, quantity int) error {
	if quantity < 1 {
		quantity = 1
	}
	// A cap, because the quantity comes from a request. Without one, a single
	// POST can ask for two billion of something and the subtotal overflows.
	if quantity > 99 {
		quantity = 99
	}

	var product models.Product
	if err := s.shopDB(ctx).
		Where("id = ?", productID).
		Where("archived_at IS NULL").
		First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotForSale
		}
		return fmt.Errorf("loading the product: %w", err)
	}
	if !product.Available {
		return ErrNotForSale
	}

	price := product.Price
	label := ""

	if variantID != "" {
		var variant models.ProductVariant
		if err := s.shopDB(ctx).
			Where("id = ?", variantID).
			Where("product_id = ?", product.ID).
			First(&variant).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotForSale
			}
			return fmt.Errorf("loading the variant: %w", err)
		}
		if !variant.InStock() {
			return ErrNotForSale
		}
		// The variant's resolved price, which the variant service computes from
		// the base price, the per-value deltas and any override. Asking it
		// rather than recomputing keeps one answer to "what does this cost".
		//
		// It needs the variant's own option values loaded and the options they
		// belong to, because only an option marked affects_price contributes
		// its delta: a colour does not change the price and a size does.
		variants := NewProductVariantService(s.DB)
		if err := s.shopDB(ctx).Preload("OptionValues").First(&variant, "id = ?", variant.ID).Error; err != nil {
			return fmt.Errorf("loading the variant's option values: %w", err)
		}
		options, err := variants.OptionsFor(product.ID)
		if err != nil {
			return fmt.Errorf("loading the product's options: %w", err)
		}
		byID := make(map[string]models.Option, len(options))
		for _, o := range options {
			byID[o.ID] = o
		}
		price = variants.ResolvePrice(product.Price, variant, byID)

		label = s.variantLabel(variant, options)
	}

	image := ""
	if product.FeaturedImage != nil {
		image = product.FeaturedImage.URL
	} else if len(product.Images) > 0 {
		image = product.Images[0].URL
	}

	// The same variant twice is one line with a higher quantity. Matching on
	// cart, product AND variant: the same product in two sizes is two lines,
	// which is the whole point of variants.
	var existing models.CartItem
	err := s.shopDB(ctx).
		Where("cart_id = ? AND product_id = ? AND variant_id = ?", cart.ID, product.ID, variantID).
		First(&existing).Error
	switch {
	case err == nil:
		want := existing.Quantity + quantity
		if want > 99 {
			want = 99
		}
		if err := s.shopDB(ctx).Model(&existing).Update("quantity", want).Error; err != nil {
			return fmt.Errorf("raising the quantity: %w", err)
		}
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		line := models.CartItem{
			CartID:       cart.ID,
			ProductID:    product.ID,
			VariantID:    variantID,
			Quantity:     quantity,
			UnitPrice:    price,
			Title:        product.Title,
			VariantLabel: label,
			ImageURL:     image,
		}
		if err := s.shopDB(ctx).Create(&line).Error; err != nil {
			return fmt.Errorf("adding the line: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("looking for an existing line: %w", err)
	}
}

// SetQuantity changes a line's quantity, removing the line at zero.
func (s *ShopCartService) SetQuantity(ctx context.Context, cart *models.Cart, lineID string, quantity int) error {
	if quantity < 1 {
		return s.Remove(ctx, cart, lineID)
	}
	if quantity > 99 {
		quantity = 99
	}
	// Scoped to this cart, so a guessed line id from somebody else's cart
	// changes nothing rather than changing theirs.
	res := s.shopDB(ctx).Model(&models.CartItem{}).
		Where("id = ? AND cart_id = ?", lineID, cart.ID).
		Update("quantity", quantity)
	if res.Error != nil {
		return fmt.Errorf("setting the quantity: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotInCart
	}
	return nil
}

// Remove takes a line out of the cart.
func (s *ShopCartService) Remove(ctx context.Context, cart *models.Cart, lineID string) error {
	res := s.shopDB(ctx).
		Where("id = ? AND cart_id = ?", lineID, cart.ID).
		Delete(&models.CartItem{})
	if res.Error != nil {
		return fmt.Errorf("removing the line: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotInCart
	}
	return nil
}

// variantLabel builds "Black / XL" for a variant, for the cart line.
//
// From the values already preloaded onto the variant and the options already
// loaded for the product, so it costs no query. The first version of this ran
// a three-table join and named the join column product_variant_id, which the
// table does not have: it is variant_id. Reaching into a generated schema with
// raw SQL is exactly the thing that breaks when the generator changes, and the
// association was right there.
//
// Ordered by the option's position, because "Black / XL" and "XL / Black" are
// the same variant and only one of them reads like a label. The options slice
// is already in that order, which is what OptionsFor promises.
//
// Stored on the line rather than computed when the cart is read: an option
// value renamed next month must not retitle something somebody already bought,
// and after checkout the line is the only record of what was chosen.
func (s *ShopCartService) variantLabel(variant models.ProductVariant, options []models.Option) string {
	chosen := make(map[string]string, len(variant.OptionValues))
	for _, value := range variant.OptionValues {
		chosen[value.OptionID] = value.Label
	}
	parts := make([]string, 0, len(options))
	for _, option := range options {
		if label, ok := chosen[option.ID]; ok {
			parts = append(parts, label)
		}
	}
	return strings.Join(parts, " / ")
}
