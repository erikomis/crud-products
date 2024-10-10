package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"goravel/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (p *ProductRepository) GetProduct() ([]model.Product, error) {

	query := "SELECT * FROM products"
	rows, err := p.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}
	var products []model.Product

	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(&productObj.ID, &productObj.Name, &productObj.Description)
		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}
		products = append(products, productObj)
	}
	rows.Close()
	return products, nil
}

func (p *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int

	query, err := p.connection.Prepare("INSERT INTO products (name, description) VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Description).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	query.Close()
	return id, nil

}

func (pr *ProductRepository) GetProductByID(id int) (*model.Product, error) {

	query, err := pr.connection.Prepare("SELECT * FROM products WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product model.Product

	err = query.QueryRow(id).Scan(&product.ID, &product.Name, &product.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query.Close()

	return &product, nil

}

var ErrNotFound = errors.New("product not found")

func (p *ProductRepository) DeleteProduct(id int) error {
	result, err := p.connection.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *ProductRepository) UpdateProduct(product model.Product) (model.Product, error) {

	result, err := p.connection.Exec("UPDATE products SET name = $1, description = $2 WHERE id = $3",
		product.Name, product.Description, product.ID)
	if err != nil {
		return model.Product{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Product{}, err
	}

	if rowsAffected == 0 {
		return model.Product{}, ErrNotFound
	}

	return product, nil
}
