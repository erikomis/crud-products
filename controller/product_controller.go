package controller

import (
	"errors"
	"goravel/model"
	"goravel/repository"
	"goravel/usecase"
	"goravel/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(
	productUseCase usecase.ProductUseCase,
) productController {
	return productController{
		productUseCase: productUseCase,
	}
}

func (p *productController) GetProduct(ctx *gin.Context) {

	products, err := p.productUseCase.GetProduct()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

	}

	ctx.JSON(http.StatusOK, products)
}

func (p *productController) CreateProduct(ctx *gin.Context) {

	var product model.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	product, err = p.productUseCase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, product)
}

func (p *productController) GetProductById(ctx *gin.Context) {
	idParams := ctx.Param("id")

	if idParams == "" {
		response := model.Response{
			Status:  http.StatusBadRequest,
			Message: "ID is required",
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id, err := strconv.Atoi(idParams)

	if err != nil {
		response := model.Response{
			Status:  http.StatusBadRequest,
			Message: "ID must be a number",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	product, err := p.productUseCase.GetProductById(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Status:  http.StatusNotFound,
			Message: "Product not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	idParams := ctx.Param("id")

	id, ok := validation.ValideParamId(idParams, ctx)

	if !ok {
		return
	}

	if err := p.productUseCase.DeleteProduct(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response := model.Response{
				Status:  http.StatusNotFound,
				Message: "Product not found",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete product",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := model.Response{
		Status:  http.StatusOK,
		Message: "Product deleted successfully",
	}
	ctx.JSON(http.StatusOK, response)
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	idParams := ctx.Param("id")

	id, ok := validation.ValideParamId(idParams, ctx)

	if !ok {
		return
	}

	var product model.Product

	if err := ctx.BindJSON(&product); err != nil {
		response := model.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid product data",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product.ID = id

	updatedProduct, err := p.productUseCase.UpdateProduct(product)
	if err != nil {
		if err.Error() == "product not found" {
			response := model.Response{
				Status:  http.StatusNotFound,
				Message: "Product not found",
			}
			ctx.JSON(http.StatusNotFound, response)
		} else {
			response := model.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := model.Response{
		Status:  http.StatusOK,
		Message: "Product updated successfully",
		Data:    updatedProduct,
	}
	ctx.JSON(http.StatusOK, response)
}
