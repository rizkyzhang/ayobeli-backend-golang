package usecase_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
	"github.com/rizkyzhang/ayobeli-backend-golang/repository"
	"github.com/rizkyzhang/ayobeli-backend-golang/usecase"
	"github.com/stretchr/testify/suite"
)

type CartUsecaseSuite struct {
	suite.Suite
	db             *sqlx.DB
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	ctx            context.Context
	now            time.Time
	nowUTC         time.Time
	authRepo       domain.AuthRepository
	cartRepo       domain.CartRepository
	productRepo    domain.ProductRepository
	aesEncryptUtil domain.AesEncryptUtil
	cartUtil       domain.CartUtil
	productUtil    domain.ProductUtil
	productUIDS    []string
	cartItemUID    string
	userID         int
}

func (s *CartUsecaseSuite) BeforeTest(suiteName, testName string) {
	metadata := utils.GenerateMetadata()
	ID, err := s.authRepo.CreateUser(&domain.AuthRepositoryPayloadCreateUser{
		UID:          metadata.UID(),
		Email:        gofakeit.Email(),
		Name:         gofakeit.Name(),
		Phone:        gofakeit.Phone(),
		ProfileImage: gofakeit.ImageURL(100, 100),
		CreatedAt:    metadata.CreatedAt,
		UpdatedAt:    metadata.UpdatedAt,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.userID = ID

	for i := 1; i <= 10; i++ {
		metadata := utils.GenerateMetadata()
		name := fmt.Sprintf("Product Test %d", i)
		weightValue := gofakeit.Float64Range(100.0, 100_000.0)
		weight := s.productUtil.FormatWeight(weightValue)
		basePriceValue := gofakeit.IntRange(5000, 1_000_000)
		discount := gofakeit.IntRange(0, 100)
		stock := gofakeit.IntRange(0, 100)
		computedPrice, _ := s.productUtil.CalculatePrice(basePriceValue, discount)
		sku := gofakeit.LoremIpsumWord() + fmt.Sprint(i)

		payload := &domain.ProductRepositoryPayloadCreateProduct{
			UID:             metadata.UID(),
			Name:            name,
			Slug:            metadata.Slug(name),
			SKU:             sku,
			Description:     gofakeit.Sentence(100),
			Images:          domain.StringSlice{"test.jpg"},
			Weight:          weight,
			WeightValue:     weightValue,
			BasePrice:       computedPrice.Base,
			BasePriceValue:  basePriceValue,
			OfferPrice:      computedPrice.Offer,
			OfferPriceValue: computedPrice.OfferValue,
			Status:          "ACTIVE",
			Discount:        discount,
			Stock:           stock,
			CreatedAt:       metadata.CreatedAt,
			UpdatedAt:       metadata.UpdatedAt,
		}

		_, err := s.productRepo.Create(payload)
		if err != nil {
			log.Fatal(err)
		}

		s.productUIDS = append(s.productUIDS, payload.UID)
	}
}

func (s *CartUsecaseSuite) SetupTest() {
	env := utils.LoadConfig("../.env")
	pool, resource, db := utils.SetupTestDB(env)

	s.pool = pool
	s.resource = resource
	s.db = db

	ctx := context.Background()
	now := time.Now()
	authRepo := repository.NewAuthRepository(s.db)
	cartRepo := repository.NewCartRepository(s.db)
	productRepo := repository.NewProductRepository(s.db)
	aesEncryptUtil := utils.NewAesEncrypt(env.AesSecret)
	productUtil := utils.NewProductUtil()
	cartUtil := utils.NewCartUtil(productUtil)

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.authRepo = authRepo
	s.cartRepo = cartRepo
	s.productRepo = productRepo
	s.aesEncryptUtil = aesEncryptUtil
	s.cartUtil = cartUtil
	s.productUtil = productUtil
}

func (s *CartUsecaseSuite) TearDownTest() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestCartUsecaseSuite(t *testing.T) {
	suite.Run(t, new(CartUsecaseSuite))
}

func (s *CartUsecaseSuite) TestCartUsecase() {
	s.Run("Create n cart items", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := s.cartRepo.GetCartByUserID(1)
		s.NoError(err)
		s.NotNil(cart)

		for i := 0; i < 3; i++ {
			product, err := s.productRepo.GetByUID(s.productUIDS[i])
			s.NoError(err)

			payload := &domain.CartUsecasePayloadCreateCartItem{
				Cart:     cart,
				Product:  product,
				Quantity: gofakeit.IntRange(1, 10),
			}
			UID, err := uc.CreateCartItem(payload)
			s.NoError(err)
			s.cartItemUID = UID

			calculatedCart, err := s.cartUtil.CalculateCreateCartItem(payload)
			s.NoError(err)
			cart, err = s.cartRepo.GetCartByUserID(1)
			s.NoError(err)
			s.Equal(product.BasePrice, cart.CartItems[i].BasePrice)
			s.Equal(product.BasePriceValue, cart.CartItems[i].BasePriceValue)
			s.Equal(product.OfferPrice, cart.CartItems[i].OfferPrice)
			s.Equal(product.OfferPriceValue, cart.CartItems[i].OfferPriceValue)
			s.Equal(product.Discount, cart.CartItems[i].Discount)
			s.Equal(product.Images[0], cart.CartItems[i].ProductImage)
			s.Equal(product.Name, cart.CartItems[i].ProductName)
			s.Equal(product.Slug, cart.CartItems[i].ProductSlug)
			s.Equal(product.Weight, cart.CartItems[i].ProductWeight)
			s.Equal(product.WeightValue, cart.CartItems[i].ProductWeightValue)
			s.Equal(payload.Quantity, cart.CartItems[i].Quantity)
			s.Equal(calculatedCart.CartQuantity, cart.Quantity)
			s.Equal(calculatedCart.CartTotalPrice, cart.TotalPrice)
			s.Equal(calculatedCart.CartTotalPriceValue, cart.TotalPriceValue)
			s.Equal(calculatedCart.CartTotalWeight, cart.TotalWeight)
			s.Equal(calculatedCart.CartTotalWeightValue, cart.TotalWeightValue)
			s.Equal(calculatedCart.CartItemTotalPrice, cart.CartItems[i].TotalPrice)
			s.Equal(calculatedCart.CartItemTotalPriceValue, cart.CartItems[i].TotalPriceValue)
			s.Equal(calculatedCart.CartItemTotalWeight, cart.CartItems[i].TotalWeight)
			s.Equal(calculatedCart.CartItemTotalWeightValue, cart.CartItems[i].TotalWeightValue)
		}
	})

	s.Run("Update cart item by uid", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := s.cartRepo.GetCartByUserID(1)
		s.NoError(err)
		s.NotNil(cart)
		cartItem, err := s.cartRepo.GetCartItemByUID(s.cartItemUID)
		s.NoError(err)
		s.NotNil(cartItem)

		payload := &domain.CartUsecasePayloadUpdateCartItem{
			UID:      s.cartItemUID,
			Cart:     cart,
			CartItem: cartItem,
			Quantity: 5,
		}
		err = uc.UpdateCartItem(payload)
		s.NoError(err)

		calculatedCart, err := s.cartUtil.CalculateUpdateCartItem(payload)
		s.NoError(err)
		cart, err = s.cartRepo.GetCartByUserID(1)
		s.NoError(err)
		cartItem, err = s.cartRepo.GetCartItemByUID(s.cartItemUID)
		s.NoError(err)

		s.Equal(payload.Quantity, cartItem.Quantity)
		s.Equal(calculatedCart.CartQuantity, cart.Quantity)
		s.Equal(calculatedCart.CartTotalPrice, cart.TotalPrice)
		s.Equal(calculatedCart.CartTotalPriceValue, cart.TotalPriceValue)
		s.Equal(calculatedCart.CartTotalWeight, cart.TotalWeight)
		s.Equal(calculatedCart.CartTotalWeightValue, cart.TotalWeightValue)
		s.Equal(calculatedCart.CartItemTotalPrice, cartItem.TotalPrice)
		s.Equal(calculatedCart.CartItemTotalPriceValue, cartItem.TotalPriceValue)
		s.Equal(calculatedCart.CartItemTotalWeight, cartItem.TotalWeight)
		s.Equal(calculatedCart.CartItemTotalWeightValue, cartItem.TotalWeightValue)
	})

	s.Run("Get cart by user id", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := uc.GetCartByUserID(s.userID)
		s.NoError(err)
		s.NotNil(cart)
	})

	s.Run("Get cart by user id return nil given invalid user id", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := uc.GetCartByUserID(2)
		s.NoError(err)
		s.Nil(cart)
	})

	s.Run("Get cart by user id middleware", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := uc.GetCartByUserIDMiddleware(s.userID)
		s.NoError(err)
		s.Equal(s.userID, cart.UserID)
	})

	s.Run("Get cart by user id middleware return nil given invalid user id", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := uc.GetCartByUserIDMiddleware(2)
		s.NoError(err)
		s.Nil(cart)
	})

	s.Run("Get cart item by uid", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cartItem, err := uc.GetCartItemByUID(s.cartItemUID)
		s.NoError(err)
		s.NotNil(cartItem)
	})

	s.Run("Get cart item by uid return nil given invalid uid", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cartItem, err := uc.GetCartItemByUID("invalid")
		s.NoError(err)
		s.Nil(cartItem)
	})

	s.Run("Get cart item by product id", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		product, err := s.productRepo.GetByUID(s.productUIDS[0])
		s.NoError(err)

		cartItem, err := uc.GetCartItemByProductID(product.ID)
		s.NoError(err)
		s.NotNil(cartItem)
	})

	s.Run("Delete cart item by uid", func() {
		uc := usecase.NewCartUsecase(s.cartRepo, s.cartUtil)

		cart, err := s.cartRepo.GetCartByUserID(1)
		s.NoError(err)
		s.NotNil(cart)

		payload := &domain.CartUsecasePayloadDeleteCartItem{
			UID:      s.cartItemUID,
			Cart:     cart,
			CartItem: &cart.CartItems[0],
		}
		err = uc.DeleteCartItemByUID(payload)
		s.NoError(err)

		calculatedCart, err := s.cartUtil.CalculateDeleteCartItem(payload)
		s.NoError(err)
		cart, err = s.cartRepo.GetCartByUserID(1)
		s.NoError(err)
		s.Equal(calculatedCart.CartQuantity, cart.Quantity)
		s.Equal(calculatedCart.CartTotalPrice, cart.TotalPrice)
		s.Equal(calculatedCart.CartTotalPriceValue, cart.TotalPriceValue)
		s.Equal(calculatedCart.CartTotalWeight, cart.TotalWeight)
		s.Equal(calculatedCart.CartTotalWeightValue, cart.TotalWeightValue)

		cartItem, err := s.cartRepo.GetCartItemByUID(s.cartItemUID)
		s.NoError(err)
		s.Nil(cartItem)
	})
}
