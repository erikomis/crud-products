package usecase

import (
	"goravel/model"
	"goravel/repository"
)

type ProductUseCase struct {
	repository repository.ProductRepository
}

func NewProductUseCase(
	repository repository.ProductRepository,
) ProductUseCase {
	return ProductUseCase{
		repository: repository,
	}
}

func (p *ProductUseCase) GetProduct() ([]model.Product, error) {

	return p.repository.GetProduct()
}

func (p *ProductUseCase) CreateProduct(product model.Product) (model.Product, error) {

	productId, err := p.repository.CreateProduct(product)
	if err != nil {
		return model.Product{}, err
	}

	product.ID = productId

	return product, nil
}

func (p *ProductUseCase) GetProductById(id int) (*model.Product, error) {
	product, err := p.repository.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p ProductUseCase) DeleteProduct(id int) error {
	err := p.repository.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductUseCase) UpdateProduct(product model.Product) (*model.Product, error) {
	_, err := p.repository.GetProductByID(product.ID)
	if err != nil {
		return nil, err
	}

	// Atualiza o produto e retorna o produto atualizado
	productUpdate, err := p.repository.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return &productUpdate, nil
}
