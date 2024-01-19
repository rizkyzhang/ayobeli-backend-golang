package utils

import (
	"math"

	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

type baseCartUtil struct {
	productUtil domain.ProductUtil
}

func NewCartUtil(productUtil domain.ProductUtil) domain.CartUtil {
	return &baseCartUtil{productUtil: productUtil}
}

func (b *baseCartUtil) CalculateCreateCartItem(payload *domain.CartUsecasePayloadCreateCartItem) (*domain.CalculatedCart, error) {
	cartItemTotalPriceValue := payload.Quantity * payload.Product.OfferPriceValue
	cartItemTotalPrice, err := b.productUtil.FormatRupiah(cartItemTotalPriceValue)
	if err != nil {
		return nil, err
	}
	cartItemTotalWeightValue := float64(payload.Quantity) * payload.Product.WeightValue
	cartItemTotalWeightValue = math.Round(cartItemTotalWeightValue*100) / 100
	cartItemTotalWeight := b.productUtil.FormatWeight(float64(payload.Quantity) * payload.Product.WeightValue)

	cartQuantity := payload.Cart.Quantity + payload.Quantity
	cartTotalPriceValue := payload.Cart.TotalPriceValue + cartItemTotalPriceValue
	cartTotalPrice, err := b.productUtil.FormatRupiah(cartTotalPriceValue)
	if err != nil {
		return nil, err
	}
	cartTotalWeightValue := payload.Cart.TotalWeightValue + cartItemTotalWeightValue
	cartTotalWeightValue = math.Round(cartTotalWeightValue*100) / 100
	cartTotalWeight := b.productUtil.FormatWeight(cartTotalWeightValue)

	return &domain.CalculatedCart{
		CartQuantity:             cartQuantity,
		CartTotalPriceValue:      cartTotalPriceValue,
		CartTotalPrice:           cartTotalPrice,
		CartTotalWeightValue:     cartTotalWeightValue,
		CartTotalWeight:          cartTotalWeight,
		CartItemTotalPriceValue:  cartItemTotalPriceValue,
		CartItemTotalPrice:       cartItemTotalPrice,
		CartItemTotalWeightValue: cartItemTotalWeightValue,
		CartItemTotalWeight:      cartItemTotalWeight,
	}, nil
}

func (b *baseCartUtil) CalculateUpdateCartItem(payload *domain.CartUsecasePayloadUpdateCartItem) (*domain.CalculatedCart, error) {
	cartItemTotalPriceValue := payload.Quantity * payload.CartItem.OfferPriceValue
	cartItemTotalWeightValue := float64(payload.Quantity) * payload.CartItem.ProductWeightValue
	cartItemTotalWeightValue = math.Round(cartItemTotalWeightValue*100) / 100
	cartItemTotalPrice, err := b.productUtil.FormatRupiah(cartItemTotalPriceValue)
	if err != nil {
		return nil, err
	}
	cartItemTotalWeight := b.productUtil.FormatWeight(cartItemTotalWeightValue)

	quantityDiff := payload.Quantity - payload.CartItem.Quantity
	quantityDiffAbs := int(math.Abs(float64(quantityDiff)))
	isQuantityIncreased := quantityDiff > 0
	var cartQuantity int
	var cartTotalPriceValue int
	var cartTotalWeightValue float64

	if isQuantityIncreased {
		cartQuantity = payload.Cart.Quantity + quantityDiffAbs
		cartTotalPriceValue = payload.Cart.TotalPriceValue + (quantityDiffAbs * payload.CartItem.OfferPriceValue)
		cartTotalWeightValue = payload.Cart.TotalWeightValue + (float64(quantityDiffAbs) * payload.CartItem.ProductWeightValue)
	} else {
		cartQuantity = payload.Cart.Quantity - quantityDiffAbs
		cartTotalPriceValue = payload.Cart.TotalPriceValue - (quantityDiffAbs * payload.CartItem.OfferPriceValue)
		cartTotalWeightValue = payload.Cart.TotalWeightValue - (float64(quantityDiffAbs) * payload.CartItem.ProductWeightValue)
	}
	cartTotalWeightValue = math.Round(cartTotalWeightValue*100) / 100

	cartTotalPrice, err := b.productUtil.FormatRupiah(cartTotalPriceValue)
	if err != nil {
		return nil, err
	}
	cartTotalWeight := b.productUtil.FormatWeight(cartTotalWeightValue)

	return &domain.CalculatedCart{
		CartQuantity:             cartQuantity,
		CartTotalPriceValue:      cartTotalPriceValue,
		CartTotalPrice:           cartTotalPrice,
		CartTotalWeightValue:     cartTotalWeightValue,
		CartTotalWeight:          cartTotalWeight,
		CartItemTotalPriceValue:  cartItemTotalPriceValue,
		CartItemTotalPrice:       cartItemTotalPrice,
		CartItemTotalWeightValue: cartItemTotalWeightValue,
		CartItemTotalWeight:      cartItemTotalWeight,
	}, nil
}

func (b *baseCartUtil) CalculateDeleteCartItem(payload *domain.CartUsecasePayloadDeleteCartItem) (*domain.CalculatedCart, error) {
	cartQuantity := payload.Cart.Quantity - payload.CartItem.Quantity
	cartTotalPriceValue := payload.Cart.TotalPriceValue - payload.CartItem.TotalPriceValue
	cartTotalPrice, err := b.productUtil.FormatRupiah(cartTotalPriceValue)
	if err != nil {
		return nil, err
	}
	cartTotalWeightValue := payload.Cart.TotalWeightValue - payload.CartItem.TotalWeightValue
	cartTotalWeightValue = math.Round(cartTotalWeightValue*100) / 100
	cartTotalWeight := b.productUtil.FormatWeight(payload.Cart.TotalWeightValue - payload.CartItem.TotalWeightValue)

	return &domain.CalculatedCart{
		CartQuantity:         cartQuantity,
		CartTotalPriceValue:  cartTotalPriceValue,
		CartTotalPrice:       cartTotalPrice,
		CartTotalWeightValue: cartTotalWeightValue,
		CartTotalWeight:      cartTotalWeight,
	}, nil
}
