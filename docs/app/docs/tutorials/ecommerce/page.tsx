import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/ecommerce')

export default function TutorialEcommercePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Build an E-Commerce Store
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a fully functional e-commerce platform with products,
                categories, orders, order items, image uploads, stock
                validation, order confirmation emails, a revenue dashboard, and
                Redis caching for popular products. This is the most
                comprehensive Grit tutorial and covers nearly every feature of
                the framework.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  "Go 1.21+ installed",
                  "Node.js 18+ and pnpm installed",
                  "Docker and Docker Compose installed",
                  "Grit CLI installed — see /docs/getting-started/installation (one-line install script)",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* ============================================================ */}
            {/* STEP 1 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create the project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Scaffold a new Grit project called <code>shopgrit</code>.
                </p>

                <CodeBlock terminal code={`grit new shopgrit\ncd shopgrit`} className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 2 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start Docker services
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start PostgreSQL, Redis, MinIO, and Mailhog. MinIO is
                  especially important for this project because we will use it
                  for product image uploads.
                </p>

                <CodeBlock terminal code="docker compose up -d" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 3 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Product resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Products are the core of the store. Generate a resource with
                  name, description, price, SKU (unique stock-keeping unit),
                  stock count, and a published flag.
                </p>

                <CodeBlock terminal code={`grit generate resource Product --fields "name:string,description:text,price:float,sku:string:unique,stock:int,published:bool"`} className="glow-purple-sm mb-4" />

                <CodeBlock filename="apps/api/internal/models/product.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Price       float64        \`gorm:"not null" json:"price" binding:"required"\`
    SKU         string         \`gorm:"size:100;uniqueIndex;not null" json:"sku" binding:"required"\`
    Stock       int            \`gorm:"default:0" json:"stock"\`
    Published   bool           \`gorm:"default:false" json:"published"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 4 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Category resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Products are organized into categories. Generate a Category
                  resource and add the relationship to Product.
                </p>

                <CodeBlock terminal code={`grit generate resource Category --fields "name:string:unique,slug:string:unique,description:text"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now add the category relationship to the Product model. Also
                  add an
                  <code> ImageURL</code> field for product images (we will
                  handle uploads in Step 8):
                </p>

                <CodeBlock filename="apps/api/internal/models/product.go &mdash; updated" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Price       float64        \`gorm:"not null" json:"price" binding:"required"\`
    SKU         string         \`gorm:"size:100;uniqueIndex;not null" json:"sku" binding:"required"\`
    Stock       int            \`gorm:"default:0" json:"stock"\`
    Published   bool           \`gorm:"default:false" json:"published"\`
    ImageURL    string         \`gorm:"size:500" json:"image_url"\`

    // Belongs to Category
    CategoryID  uint           \`gorm:"index;not null" json:"category_id" binding:"required"\`
    Category    Category       \`gorm:"foreignKey:CategoryID" json:"category,omitempty"\`

    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 5 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                5
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Order resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Orders track customer purchases. Generate an Order resource
                  with status, total, and notes fields. An order belongs to the
                  authenticated user.
                </p>

                <CodeBlock terminal code={`grit generate resource Order --fields "status:string,total:float,notes:text"`} className="glow-purple-sm mb-4" />

                <CodeBlock filename="apps/api/internal/models/order.go &mdash; with User
                      relationship" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

// Order status constants
const (
    OrderStatusPending    = "pending"
    OrderStatusProcessing = "processing"
    OrderStatusShipped    = "shipped"
    OrderStatusDelivered  = "delivered"
    OrderStatusCancelled  = "cancelled"
)

type Order struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Status    string         \`gorm:"size:50;default:pending" json:"status"\`
    Total     float64        \`gorm:"not null" json:"total"\`
    Notes     string         \`gorm:"type:text" json:"notes"\`

    // Belongs to User (the customer)
    UserID    uint           \`gorm:"index;not null" json:"user_id"\`
    User      User           \`gorm:"foreignKey:UserID" json:"user,omitempty"\`

    // Has many order items
    Items     []OrderItem    \`gorm:"foreignKey:OrderID" json:"items,omitempty"\`

    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 6 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                6
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the OrderItem resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Each order contains one or more items. An OrderItem records
                  the quantity and the price at the time of purchase (so price
                  changes do not affect past orders).
                </p>

                <CodeBlock terminal code={`grit generate resource OrderItem --fields "quantity:int,price:float"`} className="glow-purple-sm mb-4" />

                <CodeBlock filename="apps/api/internal/models/order_item.go &mdash; with
                      relationships" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type OrderItem struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Quantity  int            \`gorm:"not null" json:"quantity" binding:"required,min=1"\`
    Price     float64        \`gorm:"not null" json:"price"\`

    // Belongs to Order
    OrderID   uint           \`gorm:"index;not null" json:"order_id"\`
    Order     Order          \`gorm:"foreignKey:OrderID" json:"order,omitempty"\`

    // References a Product
    ProductID uint           \`gorm:"index;not null" json:"product_id" binding:"required"\`
    Product   Product        \`gorm:"foreignKey:ProductID" json:"product,omitempty"\`

    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Run <code>grit sync</code> to regenerate all TypeScript types
                  with the new relationships across all four resources.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 7 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                7
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Set up all relationships
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Let&apos;s review the complete relationship diagram. Add the
                  inverse relationships to the Category model:
                </p>

                <CodeBlock filename="relationship diagram" code={`Category
  └── has many Products

Product
  └── belongs to Category

Order
  ├── belongs to User (customer)
  └── has many OrderItems

OrderItem
  ├── belongs to Order
  └── belongs to Product

User (built-in)
  └── has many Orders`} />

                <CodeBlock filename="apps/api/internal/models/category.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Category struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;uniqueIndex;not null" json:"name" binding:"required"\`
    Slug        string         \`gorm:"size:255;uniqueIndex;not null" json:"slug" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Products    []Product      \`gorm:"foreignKey:CategoryID" json:"products,omitempty"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 8 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                8
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  File uploads for product images
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Use Grit&apos;s built-in storage service to upload product
                  images to MinIO (S3-compatible). The storage service is
                  already configured &mdash; you just need to add an upload
                  handler that saves the file URL to the product.
                </p>

                <CodeBlock filename="apps/api/internal/handlers/product.go &mdash; add upload
                      method" code={`// UploadImage handles product image upload.
func (h *ProductHandler) UploadImage(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
        })
        return
    }

    // Get the file from the multipart form
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{"code": "NO_FILE", "message": "No image file provided"},
        })
        return
    }
    defer file.Close()

    // Validate file type
    contentType := header.Header.Get("Content-Type")
    if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "INVALID_TYPE",
                "message": "Only JPEG, PNG, and WebP images are allowed",
            },
        })
        return
    }

    // Upload to storage (MinIO / S3)
    path := fmt.Sprintf("products/%d/%s", id, header.Filename)
    url, err := h.storage.Upload(c, path, file, contentType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "UPLOAD_FAILED", "message": "Failed to upload image"},
        })
        return
    }

    // Update the product record with the image URL
    result := h.db.Model(&models.Product{}).Where("id = ?", id).Update("image_url", url)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "UPDATE_FAILED", "message": "Failed to update product"},
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":    gin.H{"image_url": url},
        "message": "Image uploaded successfully",
    })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Register the upload route:
                </p>

                <CodeBlock filename="apps/api/internal/routes/routes.go" code={`// Inside the authenticated products group
products.POST("/:id/upload-image", productHandler.UploadImage)`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Update the admin resource definition to include the file
                  upload field and display the product image in the table:
                </p>

                <CodeBlock language="typescript" filename="apps/admin/resources/products.ts" code={`import { defineResource } from '@grit/admin'

export default defineResource({
  name: 'Product',
  endpoint: '/api/products',
  icon: 'Package',

  table: {
    columns: [
      { key: 'id', label: 'ID', sortable: true },
      { key: 'image_url', label: 'Image', format: 'image' },
      { key: 'name', label: 'Name', sortable: true, searchable: true },
      { key: 'sku', label: 'SKU', sortable: true },
      { key: 'category.name', label: 'Category', relation: 'category' },
      { key: 'price', label: 'Price', sortable: true, format: 'currency' },
      { key: 'stock', label: 'Stock', sortable: true },
      { key: 'published', label: 'Status', badge: {
        true: { color: 'green', label: 'Published' },
        false: { color: 'yellow', label: 'Draft' },
      }},
    ],
    filters: [
      { key: 'published', type: 'select', options: [
        { label: 'Published', value: 'true' },
        { label: 'Draft', value: 'false' },
      ]},
      { key: 'category_id', type: 'select', resource: 'categories',
        displayKey: 'name', label: 'Category' },
      { key: 'price', type: 'number-range' },
    ],
    actions: ['create', 'edit', 'delete', 'export'],
    bulkActions: ['delete', 'export'],
  },

  form: {
    fields: [
      { key: 'name', label: 'Name', type: 'text', required: true },
      { key: 'sku', label: 'SKU', type: 'text', required: true },
      { key: 'description', label: 'Description', type: 'textarea' },
      { key: 'category_id', label: 'Category', type: 'relation',
        resource: 'categories', displayKey: 'name', required: true },
      { key: 'price', label: 'Price', type: 'number', prefix: '$', required: true },
      { key: 'stock', label: 'Stock', type: 'number', default: 0 },
      { key: 'image', label: 'Product Image', type: 'file',
        accept: 'image/*', uploadEndpoint: '/api/products/{id}/upload-image' },
      { key: 'published', label: 'Published', type: 'toggle', default: false },
    ],
  },
})`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 9 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                9
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Custom order creation with stock validation
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The default CRUD handler is not sufficient for order creation.
                  We need to validate stock availability, calculate totals,
                  decrement stock, and create order items &mdash; all within a
                  database transaction.
                </p>

                <CodeBlock filename="apps/api/internal/handlers/order.go &mdash; CreateOrder" code={`type CreateOrderInput struct {
    Notes string \`json:"notes"\`
    Items []struct {
        ProductID uint \`json:"product_id" binding:"required"\`
        Quantity  int  \`json:"quantity" binding:"required,min=1"\`
    } \`json:"items" binding:"required,min=1"\`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    user, _ := c.Get("user")
    currentUser := user.(*models.User)

    var input CreateOrderInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }

    // Use a transaction to ensure atomicity
    tx := h.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    var total float64
    var orderItems []models.OrderItem

    for _, item := range input.Items {
        // Lock the product row for update
        var product models.Product
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            First(&product, item.ProductID).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "error": gin.H{
                    "code":    "PRODUCT_NOT_FOUND",
                    "message": fmt.Sprintf("Product %d not found", item.ProductID),
                },
            })
            return
        }

        // Validate stock
        if product.Stock < item.Quantity {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "error": gin.H{
                    "code": "INSUFFICIENT_STOCK",
                    "message": fmt.Sprintf(
                        "Insufficient stock for %s: requested %d, available %d",
                        product.Name, item.Quantity, product.Stock,
                    ),
                },
            })
            return
        }

        // Decrement stock
        tx.Model(&product).Update("stock", product.Stock-item.Quantity)

        // Calculate line total
        lineTotal := product.Price * float64(item.Quantity)
        total += lineTotal

        orderItems = append(orderItems, models.OrderItem{
            ProductID: product.ID,
            Quantity:  item.Quantity,
            Price:     product.Price, // snapshot the price at time of purchase
        })
    }

    // Create the order
    order := models.Order{
        UserID: currentUser.ID,
        Status: models.OrderStatusPending,
        Total:  total,
        Notes:  input.Notes,
    }

    if err := tx.Create(&order).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "CREATE_FAILED", "message": "Failed to create order"},
        })
        return
    }

    // Create order items
    for i := range orderItems {
        orderItems[i].OrderID = order.ID
    }
    if err := tx.Create(&orderItems).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "CREATE_FAILED", "message": "Failed to create order items"},
        })
        return
    }

    tx.Commit()

    // Enqueue confirmation email (non-blocking)
    _ = h.jobClient.EnqueueOrderConfirmation(jobs.OrderConfirmationPayload{
        OrderID:    order.ID,
        UserEmail:  currentUser.Email,
        UserName:   currentUser.Name,
        Total:      total,
        ItemCount:  len(orderItems),
    })

    // Reload with relationships
    h.db.Preload("Items.Product").Preload("User").First(&order, order.ID)

    c.JSON(http.StatusCreated, gin.H{
        "data":    order,
        "message": "Order created successfully",
    })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Register the custom order creation route:
                </p>

                <CodeBlock filename="apps/api/internal/routes/routes.go" code={`// Replace the default POST /api/orders with the custom handler
orders := auth.Group("/orders")
{
    orders.POST("", orderHandler.CreateOrder)    // custom creation
    orders.GET("", orderHandler.GetAll)           // generated
    orders.GET("/:id", orderHandler.GetByID)      // generated
    // No PUT or DELETE — orders are immutable once created
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 10 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                10
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Background job for order confirmation emails
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  When an order is placed, queue a background job that sends a
                  confirmation email. This keeps the API response fast.
                </p>

                <CodeBlock filename="apps/api/internal/jobs/order_confirmation.go" code={`package jobs

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hibiken/asynq"
    "shopgrit/apps/api/internal/mail"
)

const TypeOrderConfirmation = "order:confirmation"

type OrderConfirmationPayload struct {
    OrderID   uint    \`json:"order_id"\`
    UserEmail string  \`json:"user_email"\`
    UserName  string  \`json:"user_name"\`
    Total     float64 \`json:"total"\`
    ItemCount int     \`json:"item_count"\`
}

// EnqueueOrderConfirmation creates a new order confirmation email job.
func (c *Client) EnqueueOrderConfirmation(payload OrderConfirmationPayload) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    task := asynq.NewTask(TypeOrderConfirmation, data)
    _, err = c.client.Enqueue(task, asynq.MaxRetry(5), asynq.Queue("critical"))
    return err
}

// HandleOrderConfirmation processes the order confirmation email.
func HandleOrderConfirmation(mailer *mail.Mailer) asynq.HandlerFunc {
    return func(ctx context.Context, t *asynq.Task) error {
        var payload OrderConfirmationPayload
        if err := json.Unmarshal(t.Payload(), &payload); err != nil {
            return fmt.Errorf("failed to unmarshal payload: %w", err)
        }

        return mailer.Send(ctx, mail.SendOptions{
            To:       payload.UserEmail,
            Subject:  fmt.Sprintf("Order #%d Confirmed", payload.OrderID),
            Template: "order-confirmation",
            Data: map[string]interface{}{
                "Name":      payload.UserName,
                "OrderID":   payload.OrderID,
                "Total":     fmt.Sprintf("$%.2f", payload.Total),
                "ItemCount": payload.ItemCount,
            },
        })
    }
}`} />

                <CodeBlock language="markup" filename="apps/api/internal/mail/templates/order-confirmation.html" code={`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: 'Onest', sans-serif; background: #0a0a0f; color: #e8e8f0; padding: 40px;">
  <div style="max-width: 500px; margin: 0 auto; background: #111118; border-radius: 12px; padding: 32px; border: 1px solid #2a2a3a;">
    <h2 style="color: #00b894; margin-top: 0;">Order Confirmed</h2>
    <p>Hi {{.Name}},</p>
    <p>Your order has been confirmed and is being processed.</p>
    <div style="background: #1a1a24; border-radius: 8px; padding: 16px; margin: 16px 0; border: 1px solid #2a2a3a;">
      <p style="margin: 0;">
        <strong>Order #{{.OrderID}}</strong>
      </p>
      <p style="margin: 4px 0 0; color: #9090a8; font-size: 14px;">
        {{.ItemCount}} item(s) &bull; Total: {{.Total}}
      </p>
    </div>
    <p style="color: #9090a8; font-size: 14px;">
      We'll send you another email when your order ships.
    </p>
  </div>
</body>
</html>`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 11 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                11
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Admin dashboard with revenue stats
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Create a stats endpoint that returns revenue totals, order
                  counts, and product inventory stats. Then build a dashboard
                  widget for the admin panel.
                </p>

                <CodeBlock filename="apps/api/internal/handlers/order.go &mdash; revenue stats
                      endpoint" code={`// GetRevenueStats returns store-wide revenue and order statistics.
func (h *OrderHandler) GetRevenueStats(c *gin.Context) {
    var totalRevenue float64
    var totalOrders int64
    var pendingOrders int64
    var totalProducts int64
    var lowStockProducts int64

    // Revenue from non-cancelled orders
    h.db.Model(&models.Order{}).
        Where("status != ?", models.OrderStatusCancelled).
        Select("COALESCE(SUM(total), 0)").
        Scan(&totalRevenue)

    h.db.Model(&models.Order{}).Count(&totalOrders)

    h.db.Model(&models.Order{}).
        Where("status = ?", models.OrderStatusPending).
        Count(&pendingOrders)

    h.db.Model(&models.Product{}).
        Where("published = ?", true).
        Count(&totalProducts)

    // Products with stock below 10
    h.db.Model(&models.Product{}).
        Where("stock < ? AND published = ?", 10, true).
        Count(&lowStockProducts)

    // Revenue by day for the last 7 days
    type DailyRevenue struct {
        Date    string  \`json:"date"\`
        Revenue float64 \`json:"revenue"\`
    }
    var dailyRevenue []DailyRevenue
    h.db.Model(&models.Order{}).
        Select("DATE(created_at) as date, SUM(total) as revenue").
        Where("created_at >= NOW() - INTERVAL '7 days' AND status != ?",
            models.OrderStatusCancelled).
        Group("DATE(created_at)").
        Order("date ASC").
        Scan(&dailyRevenue)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "total_revenue":     totalRevenue,
            "total_orders":      totalOrders,
            "pending_orders":    pendingOrders,
            "total_products":    totalProducts,
            "low_stock_products": lowStockProducts,
            "daily_revenue":     dailyRevenue,
        },
    })
}`} />

                <CodeBlock language="tsx" filename="apps/admin/components/widgets/revenue-stats.tsx" code={`'use client'

import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import { DollarSign, ShoppingCart, Package, AlertTriangle } from 'lucide-react'

interface RevenueStats {
  total_revenue: number
  total_orders: number
  pending_orders: number
  total_products: number
  low_stock_products: number
  daily_revenue: { date: string; revenue: number }[]
}

export function RevenueStatsWidget() {
  const { data } = useQuery<{ data: RevenueStats }>({
    queryKey: ['revenue-stats'],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/orders/stats')
      return data
    },
    refetchInterval: 60000,
  })

  const stats = data?.data

  const cards = [
    {
      label: 'Total Revenue',
      value: stats ? \`$\${stats.total_revenue.toLocaleString('en-US', {
        minimumFractionDigits: 2 })}\` : '$0.00',
      icon: DollarSign,
      color: 'text-emerald-400',
    },
    {
      label: 'Total Orders',
      value: stats?.total_orders ?? 0,
      icon: ShoppingCart,
      color: 'text-primary',
    },
    {
      label: 'Products Listed',
      value: stats?.total_products ?? 0,
      icon: Package,
      color: 'text-blue-400',
    },
    {
      label: 'Low Stock',
      value: stats?.low_stock_products ?? 0,
      icon: AlertTriangle,
      color: 'text-red-400',
    },
  ]

  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {cards.map((card) => (
          <div
            key={card.label}
            className="rounded-xl border border-border/40 bg-card/50 p-5"
          >
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs font-medium text-muted-foreground/60
                              uppercase tracking-wider">
                {card.label}
              </span>
              <card.icon className={\`h-4 w-4 \${card.color}\`} />
            </div>
            <p className="text-2xl font-bold tracking-tight">{card.value}</p>
          </div>
        ))}
      </div>

      {/* Pending orders callout */}
      {stats && stats.pending_orders > 0 && (
        <div className="rounded-xl border border-yellow-500/20 bg-yellow-500/5 p-4
                        flex items-center gap-3">
          <ShoppingCart className="h-5 w-5 text-yellow-500" />
          <p className="text-sm text-yellow-200/80">
            You have <strong>{stats.pending_orders}</strong> pending order(s) waiting
            to be processed.
          </p>
        </div>
      )}
    </div>
  )
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 12 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                12
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Cache popular products with Redis
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  High-traffic product pages should not hit the database on
                  every request. Use Grit&apos;s built-in Redis cache service to
                  cache popular products and invalidate the cache when a product
                  is updated.
                </p>

                <CodeBlock filename="apps/api/internal/services/product.go &mdash; cached
                      GetByID" code={`package services

import (
    "encoding/json"
    "fmt"
    "time"

    "gorm.io/gorm"
    "shopgrit/apps/api/internal/cache"
    "shopgrit/apps/api/internal/models"
)

type ProductService struct {
    db    *gorm.DB
    cache *cache.Cache
}

func NewProductService(db *gorm.DB, c *cache.Cache) *ProductService {
    return &ProductService{db: db, cache: c}
}

// GetByID retrieves a product by ID, using the Redis cache.
func (s *ProductService) GetByID(id uint) (*models.Product, error) {
    cacheKey := fmt.Sprintf("product:%d", id)

    // Try cache first
    cached, err := s.cache.Get(cacheKey)
    if err == nil && cached != "" {
        var product models.Product
        if err := json.Unmarshal([]byte(cached), &product); err == nil {
            return &product, nil
        }
    }

    // Cache miss — query the database
    var product models.Product
    result := s.db.Preload("Category").First(&product, id)
    if result.Error != nil {
        return nil, fmt.Errorf("product not found: %w", result.Error)
    }

    // Store in cache for 5 minutes
    data, _ := json.Marshal(product)
    _ = s.cache.Set(cacheKey, string(data), 5*time.Minute)

    return &product, nil
}

// Update updates a product and invalidates the cache.
func (s *ProductService) Update(id uint, input map[string]interface{}) (*models.Product, error) {
    result := s.db.Model(&models.Product{}).Where("id = ?", id).Updates(input)
    if result.Error != nil {
        return nil, fmt.Errorf("failed to update product: %w", result.Error)
    }

    // Invalidate the cache for this product
    cacheKey := fmt.Sprintf("product:%d", id)
    _ = s.cache.Delete(cacheKey)

    // Also invalidate the published products list cache
    _ = s.cache.DeletePattern("products:published:*")

    return s.GetByID(id)
}

// GetPublished returns published products with caching.
func (s *ProductService) GetPublished(page, pageSize int) ([]models.Product, int64, error) {
    cacheKey := fmt.Sprintf("products:published:%d:%d", page, pageSize)

    // Try cache first
    type cachedResult struct {
        Products []models.Product \`json:"products"\`
        Total    int64            \`json:"total"\`
    }

    cached, err := s.cache.Get(cacheKey)
    if err == nil && cached != "" {
        var result cachedResult
        if err := json.Unmarshal([]byte(cached), &result); err == nil {
            return result.Products, result.Total, nil
        }
    }

    // Cache miss — query the database
    var products []models.Product
    var total int64

    query := s.db.Model(&models.Product{}).Where("published = ?", true)
    query.Count(&total)

    offset := (page - 1) * pageSize
    query.Preload("Category").
        Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&products)

    // Cache for 2 minutes
    data, _ := json.Marshal(cachedResult{Products: products, Total: total})
    _ = s.cache.Set(cacheKey, string(data), 2*time.Minute)

    return products, total, nil
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  The cache is automatically invalidated when a product is
                  updated through the admin panel or API. Published product
                  listings are cached for 2 minutes, and individual product
                  pages are cached for 5 minutes. This dramatically reduces
                  database load under high traffic.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 13 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                13
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Run and test everything
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start all services and test the complete e-commerce workflow.
                </p>

                <CodeBlock terminal code="grit dev" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Walk through this testing checklist:
                </p>

                <ol className="space-y-2 text-[13px] text-muted-foreground/70 list-decimal pl-5">
                  <li>
                    Open the admin panel at <code>http://localhost:3001</code>{" "}
                    and register an admin account.
                  </li>
                  <li>
                    Create a category: &quot;Electronics&quot; with slug
                    &quot;electronics&quot;.
                  </li>
                  <li>
                    Create a product: &quot;Wireless Keyboard&quot;, SKU
                    &quot;WK-001&quot;, price $49.99, stock 50, category
                    &quot;Electronics&quot;, published.
                  </li>
                  <li>
                    Upload an image for the product using the image upload
                    field.
                  </li>
                  <li>
                    Verify the image appears in the product table and in MinIO
                    at <code>http://localhost:9001</code>.
                  </li>
                  <li>
                    Place an order via the API:
                    <CodeBlock language="bash" code={`curl -X POST http://localhost:8080/api/orders \\
  -H "Authorization: Bearer <your-token>" \\
  -H "Content-Type: application/json" \\
  -d '{
    "items": [
      { "product_id": 1, "quantity": 2 }
    ],
    "notes": "Please gift wrap"
  }'`} />
                  </li>
                  <li>
                    Check that the product stock decreased from 50 to 48 in the
                    admin panel.
                  </li>
                  <li>
                    Check Mailhog at <code>http://localhost:8025</code> for the
                    order confirmation email.
                  </li>
                  <li>
                    Try ordering more than the available stock &mdash; the API
                    should return an <code>INSUFFICIENT_STOCK</code> error.
                  </li>
                  <li>
                    View the revenue dashboard in the admin panel &mdash; it
                    should show $99.98 in revenue.
                  </li>
                  <li>
                    Request the same product twice via API and verify the second
                    request is served from the Redis cache (check response
                    times).
                  </li>
                  <li>
                    Browse all tables in GORM Studio at{" "}
                    <code>http://localhost:8080/studio</code>.
                  </li>
                </ol>
              </div>
            </div>

            {/* ============================================================ */}
            {/* SUMMARY */}
            {/* ============================================================ */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                What you&apos;ve built
              </h2>
              <ul className="space-y-2.5">
                {[
                  "A complete e-commerce store with Go API and Next.js frontend",
                  "Product catalog with categories, pricing, SKU tracking, and stock management",
                  "Product image uploads to S3-compatible storage (MinIO in development)",
                  "Transactional order creation with stock validation and atomic database operations",
                  "Order items that snapshot product prices at time of purchase",
                  "Background job queue for order confirmation emails",
                  "HTML email templates styled to match the Grit dark theme",
                  "Revenue dashboard with daily revenue chart, order counts, and low-stock alerts",
                  "Redis caching for product pages and published product listings",
                  "Cache invalidation on product updates to keep data fresh",
                  "Admin panel with product image previews, currency formatting, and status badges",
                  "Five resources total: Product, Category, Order, OrderItem, and the built-in User",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="prose-grit mb-8">
              <h2>Next steps</h2>
              <p>
                You have a solid e-commerce foundation. Here are ideas to take
                it further:
              </p>
              <ul>
                <li>
                  <strong>Stripe integration</strong> &mdash; add a payment
                  processing step to the order creation flow using the Stripe Go
                  SDK.
                </li>
                <li>
                  <strong>Product variants</strong> &mdash; generate a{" "}
                  <code>Variant</code> resource (size, color) linked to products
                  with separate stock tracking.
                </li>
                <li>
                  <strong>Wishlist</strong> &mdash; generate a{" "}
                  <code>Wishlist</code> resource linked to User and Product for
                  saved items.
                </li>
                <li>
                  <strong>Search</strong> &mdash; add full-text search on
                  product names and descriptions using PostgreSQL&apos;s{" "}
                  <code>tsvector</code>.
                </li>
                <li>
                  <strong>Reviews</strong> &mdash; generate a{" "}
                  <code>Review</code> resource with a rating field and link it
                  to Product and User.
                </li>
                <li>
                  <strong>Discount codes</strong> &mdash; add a{" "}
                  <code>Coupon</code> resource and apply discounts during order
                  creation.
                </li>
                <li>
                  <strong>Shipping tracking</strong> &mdash; update order status
                  through the workflow and send shipping notification emails.
                </li>
                <li>
                  <strong>Analytics</strong> &mdash; use the Grit admin chart
                  widgets to show sales by category, top-selling products, and
                  customer retention.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/saas" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Build a SaaS
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/product-catalog" className="gap-1.5">
                  Build a Product Catalog
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
