package product

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) Update(product *Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

func (r *Repository) GetByID(id uint) (*Product, error) {
	var product Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) GetAll() ([]Product, error) {
	var products []Product
	err := r.db.Find(&products).Error
	return products, err
}
