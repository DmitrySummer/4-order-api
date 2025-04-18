package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *User) error {
	return r.db.Create(product).Error
}

func (r *Repository) Update(product *User) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *Repository) GetByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByPhone(phone string) (*User, error) {
	var user User
	err := r.db.First(&user, phone).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetAll() ([]User, error) {
	var user []User
	err := r.db.Find(&user).Error
	return user, err
}
