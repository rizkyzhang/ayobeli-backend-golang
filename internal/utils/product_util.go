package utils

import (
	"fmt"

	"github.com/bojanz/currency"
	"github.com/rizkyzhang/ayobeli-backend/domain"
)

type baseProductUtil struct{}

func NewProductUtil() domain.ProductUtil {
	return &baseProductUtil{}
}

func (b *baseProductUtil) FormatRupiah(value uint64) (string, error) {
	amount, err := currency.NewAmount(fmt.Sprint(value), "IDR")
	if err != nil {
		return "", err
	}
	locale := currency.NewLocale("id")
	formatter := currency.NewFormatter(locale)
	formatter.MaxDigits = 0

	return formatter.Format(amount), nil
}

func (b *baseProductUtil) CalculatePrice(baseValue uint64, discount uint8) (*domain.CalculatedPrice, error) {
	base, err := b.FormatRupiah(uint64(baseValue))
	if err != nil {
		return nil, err
	}

	offerValue := baseValue
	if discount > 0 {
		discountValue := float32(discount) / float32(100)
		offerValue = baseValue - uint64(float32(baseValue)*discountValue)
	}
	offer, err := b.FormatRupiah(uint64(offerValue))
	if err != nil {
		return nil, err
	}

	return &domain.CalculatedPrice{
		Base:       base,
		Offer:      offer,
		OfferValue: offerValue,
	}, nil
}

func (b *baseProductUtil) FormatWeight(weightInGram float64) string {
	if weightInGram >= 1000 {
		weightInKG := weightInGram / 1000

		return fmt.Sprintf("%.2fkg", weightInKG)
	}

	return fmt.Sprintf("%.2fgr", weightInGram)
}
