package usecase

import (
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseProductUsecase struct {
	productRepository domain.ProductRepository
	aesEncryptUtil    domain.AesEncryptUtil
	productUtil       domain.ProductUtil
}

func NewProductUsecase(productRepository domain.ProductRepository, aesEncryptUtil domain.AesEncryptUtil, productUtil domain.ProductUtil) domain.ProductUsecase {
	return &baseProductUsecase{productRepository: productRepository, aesEncryptUtil: aesEncryptUtil, productUtil: productUtil}
}

func (b *baseProductUsecase) Create(payload *domain.ProductUsecasePayloadCreateProduct) (string, error) {
	metadata := utils.GenerateMetadata()
	computedPrice, err := b.productUtil.CalculatePrice(payload.BasePriceValue, payload.Discount)
	if err != nil {
		return "", err
	}
	formattedWeight := b.productUtil.FormatWeight(payload.WeightValue)

	productPayload := &domain.ProductRepositoryPayloadCreateProduct{
		UID:             metadata.UID(),
		Name:            payload.Name,
		Slug:            metadata.Slug(payload.Name),
		SKU:             payload.SKU,
		Description:     payload.Description,
		Images:          payload.Images,
		Weight:          formattedWeight,
		WeightValue:     payload.WeightValue,
		BasePrice:       computedPrice.Base,
		BasePriceValue:  payload.BasePriceValue,
		OfferPrice:      computedPrice.Offer,
		OfferPriceValue: computedPrice.OfferValue,
		Discount:        payload.Discount,
		Stock:           payload.Stock,
		Status:          payload.Status,
		CreatedAt:       metadata.CreatedAt,
		UpdatedAt:       metadata.UpdatedAt,
	}

	UID, err := b.productRepository.Create(productPayload)
	if err != nil {
		return "", err
	}

	return UID, nil
}

func (b *baseProductUsecase) List(limit int, encryptedCursor, direction string) (*domain.ProductControllerResponseListProducts, error) {
	var paginationRes domain.ProductControllerResponseListProducts

	var cursor int
	if direction != "" {
		_cursor, err := b.aesEncryptUtil.Decrypt(encryptedCursor)
		if err != nil {
			return nil, err
		}
		cursor, err = strconv.Atoi(_cursor)
		if err != nil {
			return nil, err
		}
	}

	_products, err := b.productRepository.List(limit, cursor, direction)
	if err != nil {
		return nil, err
	}
	if _products == nil {
		return nil, nil
	}

	var products []*domain.ProductControllerResponseGetProductByUID
	err = copier.Copy(&products, &_products)
	if err != nil {
		return nil, err
	}

	if cursor == 0 {
		paginationRes.IsFirstPage = true
	}
	paginationRes.Products = products
	paginationRes.Limit = limit

	if direction == "" {
		paginationRes.PrevCursor = ""

		nextCursor, err := b.aesEncryptUtil.Encrypt(strconv.Itoa(int(_products[len(_products)-1].ID)))
		if err != nil {
			return nil, err
		}
		paginationRes.NextCursor = nextCursor
	} else if direction == "prev" {
		prevCursor, err := b.aesEncryptUtil.Encrypt(strconv.Itoa(int(_products[0].ID)))
		if err != nil {
			return nil, err
		}
		paginationRes.PrevCursor = prevCursor

		paginationRes.NextCursor = encryptedCursor
	} else {
		paginationRes.PrevCursor = encryptedCursor

		nextCursor, err := b.aesEncryptUtil.Encrypt(strconv.Itoa(int(_products[len(products)-1].ID)))
		if err != nil {
			return nil, err
		}
		paginationRes.NextCursor = nextCursor
	}

	return &paginationRes, nil
}

func (b *baseProductUsecase) GetByUID(UID string) (*domain.ProductControllerResponseGetProductByUID, error) {
	product, err := b.productRepository.GetByUID(UID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	var res domain.ProductControllerResponseGetProductByUID
	err = copier.Copy(&res, &product)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (b *baseProductUsecase) UpdateByUID(UID string, payload *domain.ProductUsecasePayloadUpdateProduct) error {
	metadata := utils.GenerateMetadata()
	computedPrice, err := b.productUtil.CalculatePrice(payload.BasePriceValue, payload.Discount)
	if err != nil {
		return err
	}
	formattedWeight := b.productUtil.FormatWeight(payload.WeightValue)

	productPayload := &domain.ProductRepositoryPayloadUpdateProduct{
		UID:             UID,
		Name:            payload.Name,
		Slug:            metadata.Slug(payload.Name),
		SKU:             payload.SKU,
		Description:     payload.Description,
		Images:          payload.Images,
		Weight:          formattedWeight,
		WeightValue:     payload.WeightValue,
		BasePrice:       computedPrice.Base,
		BasePriceValue:  payload.BasePriceValue,
		OfferPrice:      computedPrice.Offer,
		OfferPriceValue: computedPrice.OfferValue,
		Discount:        payload.Discount,
		Stock:           payload.Stock,
		Status:          payload.Status,
		UpdatedAt:       metadata.UpdatedAt,
	}

	err = b.productRepository.UpdateByUID(productPayload)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseProductUsecase) DeleteByUID(UID string) error {
	err := b.productRepository.DeleteByUID(UID)
	if err != nil {
		return err
	}

	return nil
}
