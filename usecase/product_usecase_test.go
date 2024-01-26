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

type ProductUsecaseSuite struct {
	suite.Suite
	db             *sqlx.DB
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	ctx            context.Context
	now            time.Time
	nowUTC         time.Time
	repo           domain.ProductRepository
	aesEncryptUtil domain.AesEncryptUtil
	productUtil    domain.ProductUtil
	productUIDS    []string
}

func (s *ProductUsecaseSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestReadDeleteProductUsecase" {
		for i := 1; i <= 12; i++ {
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

			_, err := s.repo.Create(payload)
			if err != nil {
				log.Fatal(err)
			}

			s.productUIDS = append(s.productUIDS, payload.UID)
		}
	}
}

func (s *ProductUsecaseSuite) SetupTest() {
	env := utils.LoadConfig("../.env")
	pool, resource, db := utils.SetupTestDB(env)

	s.pool = pool
	s.resource = resource
	s.db = db

	ctx := context.Background()
	now := time.Now()
	repo := repository.NewProductRepository(s.db)
	aesEncryptUtil := utils.NewAesEncrypt(env.AesSecret)
	productUtil := utils.NewProductUtil()

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.repo = repo
	s.aesEncryptUtil = aesEncryptUtil
	s.productUtil = productUtil
}

func (s *ProductUsecaseSuite) TearDownTest() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestProductUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseSuite))
}

func (s *ProductUsecaseSuite) TestCreateUpdateProductUsecase() {
	payload := &domain.ProductUsecasePayloadCreateProduct{
		Name:           "Product Test 1",
		Description:    "Test 1",
		WeightValue:    1000.0,
		BasePriceValue: 10000,
		Discount:       0,
		Stock:          10,
		Status:         "ACTIVE",
		Images:         domain.StringSlice{"test.jpg"},
		SKU:            "TEST123",
	}
	var createdProductUID string

	s.Run("Create product without discount", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		UID, err := uc.Create(payload)
		s.NoError(err)
		createdProductUID = UID

		product, err := s.repo.GetByUID(UID)
		s.NoError(err)
		s.Equal(payload.Name, product.Name)
		s.Equal("product-test-1", product.Slug)
		s.Equal(payload.Description, product.Description)
		s.Equal(payload.WeightValue, product.WeightValue)
		s.Equal("1.00kg", product.Weight)
		s.Equal(payload.Images, product.Images)
		s.Equal(payload.Stock, product.Stock)
		s.Equal(payload.Status, product.Status)
		s.Equal(payload.BasePriceValue, product.BasePriceValue)
		s.Equal(10000, product.OfferPriceValue)
		s.Equal(payload.Discount, product.Discount)
		s.Equal(payload.SKU, product.SKU.String)
	})

	s.Run("Create product with discount 10%", func() {
		payload.Name = "Product Test 2"
		payload.Description = "Test 2"
		payload.Discount = 10
		payload.SKU = "TEST321"

		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		UID, err := uc.Create(payload)
		s.NoError(err)
		createdProductUID = UID

		product, err := s.repo.GetByUID(UID)
		s.NoError(err)
		s.Equal(payload.Name, product.Name)
		s.Equal("product-test-2", product.Slug)
		s.Equal(payload.Description, product.Description)
		s.Equal(payload.WeightValue, product.WeightValue)
		s.Equal("1.00kg", product.Weight)
		s.Equal(payload.Images, product.Images)
		s.Equal(payload.Stock, product.Stock)
		s.Equal(payload.Status, product.Status)
		s.Equal(payload.BasePriceValue, product.BasePriceValue)
		s.Equal(9000, product.OfferPriceValue)
		s.Equal(payload.Discount, product.Discount)
		s.Equal(payload.SKU, product.SKU.String)
	})

	s.Run("Update product", func() {
		payload := &domain.ProductUsecasePayloadUpdateProduct{
			Name:           "Product Test 2 Updated",
			Description:    "Test 2 Updated",
			WeightValue:    1500.0,
			BasePriceValue: 15000,
			Discount:       0,
			Stock:          30,
			Status:         "ACTIVE",
			Images:         domain.StringSlice{"test.jpg", "test2.jpg"},
			SKU:            "TEST2UPDATED",
		}

		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		err := uc.UpdateByUID(createdProductUID, payload)
		s.NoError(err)

		product, err := s.repo.GetByUID(createdProductUID)
		s.NoError(err)
		s.Equal(payload.Name, product.Name)
		s.Equal("product-test-2-updated", product.Slug)
		s.Equal(payload.Description, product.Description)
		s.Equal(payload.WeightValue, product.WeightValue)
		s.Equal("1.50kg", product.Weight)
		s.Equal(payload.Images, product.Images)
		s.Equal(payload.Stock, product.Stock)
		s.Equal(payload.Status, product.Status)
		s.Equal(payload.BasePriceValue, product.BasePriceValue)
		s.Equal(15000, product.OfferPriceValue)
		s.Equal(payload.Discount, product.Discount)
		s.Equal(payload.SKU, product.SKU.String)
	})

}

func (s *ProductUsecaseSuite) TestReadDeleteProductUsecase() {
	s.Run("List products pagination for first page", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)

		paginationRes, err := uc.List(5, "", "")
		s.NoError(err)
		s.True(paginationRes.IsFirstPage)
		s.Equal(5, paginationRes.Limit)
		s.Empty(paginationRes.PrevCursor)
		s.NotEmpty(paginationRes.NextCursor)

		s.Equal("Product Test 1", paginationRes.Products[0].Name)
		s.Equal("Product Test 2", paginationRes.Products[1].Name)
		s.Equal("Product Test 3", paginationRes.Products[2].Name)
		s.Equal("Product Test 4", paginationRes.Products[3].Name)
		s.Equal("Product Test 5", paginationRes.Products[4].Name)
	})

	s.Run("List products pagination for next page", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)

		cursor, err := s.aesEncryptUtil.Encrypt("5")
		s.NoError(err)
		paginationRes, err := uc.List(5, cursor, "next")
		s.NoError(err)
		s.False(paginationRes.IsFirstPage)
		s.Equal(5, paginationRes.Limit)
		s.Equal(cursor, paginationRes.PrevCursor)
		s.NotEmpty(paginationRes.NextCursor)

		s.Equal("Product Test 6", paginationRes.Products[0].Name)
		s.Equal("Product Test 7", paginationRes.Products[1].Name)
		s.Equal("Product Test 8", paginationRes.Products[2].Name)
		s.Equal("Product Test 9", paginationRes.Products[3].Name)
		s.Equal("Product Test 10", paginationRes.Products[4].Name)
	})

	s.Run("List products pagination for last page", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)

		cursor, err := s.aesEncryptUtil.Encrypt("10")
		s.NoError(err)
		paginationRes, err := uc.List(5, cursor, "next")
		s.NoError(err)
		s.False(paginationRes.IsFirstPage)
		s.Equal(5, paginationRes.Limit)
		s.Equal(cursor, paginationRes.PrevCursor)
		s.NotEmpty(paginationRes.NextCursor)

		s.Equal("Product Test 11", paginationRes.Products[0].Name)
		s.Equal("Product Test 12", paginationRes.Products[1].Name)
	})

	s.Run("List products pagination for prev page", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)

		cursor, err := s.aesEncryptUtil.Encrypt("11")
		s.NoError(err)
		paginationRes, err := uc.List(5, cursor, "prev")
		s.NoError(err)
		s.False(paginationRes.IsFirstPage)
		s.Equal(5, paginationRes.Limit)
		s.Equal(cursor, paginationRes.NextCursor)
		s.NotEmpty(paginationRes.PrevCursor)

		s.Equal("Product Test 6", paginationRes.Products[0].Name)
		s.Equal("Product Test 7", paginationRes.Products[1].Name)
		s.Equal("Product Test 8", paginationRes.Products[2].Name)
		s.Equal("Product Test 9", paginationRes.Products[3].Name)
		s.Equal("Product Test 10", paginationRes.Products[4].Name)
	})

	s.Run("Get product", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		product, err := uc.GetByUID(s.productUIDS[0])
		s.NoError(err)
		s.NotNil(product)
	})

	s.Run("Get product return nil if product not found", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		product, err := uc.GetByUID("123")
		s.NoError(err)
		s.Nil(product)
	})

	s.Run("Delete product", func() {
		uc := usecase.NewProductUsecase(s.repo, s.aesEncryptUtil, s.productUtil)
		err := uc.DeleteByUID(s.productUIDS[0])
		s.NoError(err)

		product, err := s.repo.GetByUID(s.productUIDS[0])
		s.NoError(err)
		s.Nil(product)
	})
}
