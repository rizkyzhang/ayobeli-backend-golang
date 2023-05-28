package usecase

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type baseCartUsecase struct {
	cartRepository domain.CartRepository
	cartUtil       domain.CartUtil
}

func NewCartUsecase(cartRepository domain.CartRepository, cartUtil domain.CartUtil) domain.CartUsecase {
	return &baseCartUsecase{cartRepository: cartRepository, cartUtil: cartUtil}
}

func (b *baseCartUsecase) GetCartByUserID(userID uint64) (*domain.CartControllerResponseGetCart, error) {
	var res domain.CartControllerResponseGetCart
	cart, err := b.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, nil
	}

	err = copier.Copy(&res, &cart)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (b *baseCartUsecase) GetCartByUserIDMiddleware(userID uint64) (*domain.CartModel, error) {
	cart, err := b.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (b *baseCartUsecase) CreateCartItem(payload *domain.CartUsecasePayloadCreateCartItem) (string, error) {
	metadata := utils.GenerateMetadata()
	calculatedCart, err := b.cartUtil.CalculateCreateCartItem(payload)
	if err != nil {
		return "", err
	}

	log.Println(calculatedCart.CartTotalWeightValue)
	log.Println(calculatedCart.CartItemTotalWeightValue)

	cartItemPayload := domain.CartRepositoryPayloadCreateCartItem{
		UID:                metadata.UID(),
		Quantity:           payload.Quantity,
		TotalPriceValue:    calculatedCart.CartItemTotalPriceValue,
		TotalPrice:         calculatedCart.CartItemTotalPrice,
		TotalWeightValue:   calculatedCart.CartItemTotalWeightValue,
		TotalWeight:        calculatedCart.CartItemTotalWeight,
		ProductName:        payload.Product.Name,
		ProductSlug:        payload.Product.Slug,
		ProductImage:       payload.Product.Images[0],
		ProductWeight:      payload.Product.Weight,
		ProductWeightValue: payload.Product.WeightValue,
		BasePrice:          payload.Product.BasePrice,
		BasePriceValue:     payload.Product.BasePriceValue,
		OfferPrice:         payload.Product.OfferPrice,
		OfferPriceValue:    payload.Product.OfferPriceValue,
		Discount:           payload.Product.Discount,
		CartID:             payload.Cart.ID,
		ProductID:          payload.Product.ID,
		CreatedAt:          metadata.CreatedAt,
		UpdatedAt:          metadata.UpdatedAt,
	}

	cartPayload := domain.CartRepositoryPayloadUpdateCart{
		UID:              payload.Cart.UID,
		Quantity:         calculatedCart.CartQuantity,
		TotalPriceValue:  calculatedCart.CartTotalPriceValue,
		TotalPrice:       calculatedCart.CartTotalPrice,
		TotalWeightValue: calculatedCart.CartTotalWeightValue,
		TotalWeight:      calculatedCart.CartTotalWeight,
		UpdatedAt:        metadata.UpdatedAt,
	}

	UID, err := b.cartRepository.CreateCartItem(cartItemPayload, cartPayload)
	if err != nil {
		return "", err
	}

	return UID, nil
}

func (b *baseCartUsecase) GetCartItemByUID(UID string) (*domain.CartItemModel, error) {
	cartItem, err := b.cartRepository.GetCartItemByUID(UID)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (b *baseCartUsecase) GetCartItemByProductID(productID uint64) (*domain.CartItemModel, error) {
	cartItem, err := b.cartRepository.GetCartItemByProductID(productID)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (b *baseCartUsecase) UpdateCartItem(payload *domain.CartUsecasePayloadUpdateCartItem) error {
	metadata := utils.GenerateMetadata()
	calculatedCart, err := b.cartUtil.CalculateUpdateCartItem(payload)
	if err != nil {
		return err
	}

	cartItemPayload := domain.CartRepositoryPayloadUpdateCartItem{
		UID:              payload.UID,
		Quantity:         payload.Quantity,
		TotalPriceValue:  calculatedCart.CartItemTotalPriceValue,
		TotalPrice:       calculatedCart.CartItemTotalPrice,
		TotalWeightValue: calculatedCart.CartItemTotalWeightValue,
		TotalWeight:      calculatedCart.CartItemTotalWeight,
		UpdatedAt:        metadata.UpdatedAt,
	}

	cartPayload := domain.CartRepositoryPayloadUpdateCart{
		UID:              payload.Cart.UID,
		Quantity:         calculatedCart.CartQuantity,
		TotalPriceValue:  calculatedCart.CartTotalPriceValue,
		TotalPrice:       calculatedCart.CartTotalPrice,
		TotalWeightValue: calculatedCart.CartTotalWeightValue,
		TotalWeight:      calculatedCart.CartTotalWeight,
		UpdatedAt:        metadata.UpdatedAt,
	}

	err = b.cartRepository.UpdateCartItem(cartItemPayload, cartPayload)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCartUsecase) DeleteCartItemByUID(payload *domain.CartUsecasePayloadDeleteCartItem) error {
	metadata := utils.GenerateMetadata()
	calculatedCart, err := b.cartUtil.CalculateDeleteCartItem(payload)
	if err != nil {
		return err
	}

	cartPayload := domain.CartRepositoryPayloadUpdateCart{
		UID:              payload.Cart.UID,
		Quantity:         payload.Cart.Quantity - payload.CartItem.Quantity,
		TotalPriceValue:  calculatedCart.CartTotalPriceValue,
		TotalPrice:       calculatedCart.CartTotalPrice,
		TotalWeightValue: calculatedCart.CartTotalWeightValue,
		TotalWeight:      calculatedCart.CartTotalWeight,
		UpdatedAt:        metadata.UpdatedAt,
	}

	err = b.cartRepository.DeleteCartItemByUID(payload.UID, cartPayload)
	if err != nil {
		return err
	}

	return nil
}
