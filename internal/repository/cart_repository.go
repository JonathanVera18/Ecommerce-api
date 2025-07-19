package repository

import (
	"context"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID uint) (*models.Cart, error)
	GetCart(ctx context.Context, userID uint) (*models.Cart, error)
	AddItem(ctx context.Context, cartItem *models.CartItem) error
	UpdateItem(ctx context.Context, cartItem *models.CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID uint) error
	GetItem(ctx context.Context, cartID, itemID uint) (*models.CartItem, error)
	GetItemByProduct(ctx context.Context, cartID, productID uint) (*models.CartItem, error)
	ClearCart(ctx context.Context, userID uint) error
	GetCartWithItems(ctx context.Context, userID uint) (*models.Cart, error)
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", userID).
		First(&cart).Error
	
	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{CustomerID: userID}
		err = r.db.WithContext(ctx).Create(&cart).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	
	return &cart, nil
}

func (r *cartRepository) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *cartRepository) UpdateItem(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, itemID uint) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ? AND id = ?", cartID, itemID).
		Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetItem(ctx context.Context, cartID, itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("cart_id = ? AND id = ?", cartID, itemID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) GetItemByProduct(ctx context.Context, cartID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) ClearCart(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("cart_id IN (SELECT id FROM carts WHERE customer_id = ?)", userID).
		Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCartWithItems(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Product.ProductImages").
		Where("customer_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
