package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string    `gorm:"size:50;not null json:"firstName"`
	LastName  string    `gorm:"size:50;not null json:"lastName"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Role      string    `gorm:"size:50;not null;default:user" json:"role"`
	Photo     string    `gorm:"size:150" json:"photo"`
	Active    bool      `gorm:"default:true;" json:"active"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Encrypt passowrd
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Make sure password provided is valid
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Encrypt password before save it
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Data normalization
func (u *User) Prepare() {
	u.ID = 0
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate User's access info case email/password are blank ou email format is invalid
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if u.Password == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil
	case "update":
		if u.Password == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil

	default:
		if u.Password == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil
	}
}

// Save User
func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// Find all Users
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

// Find User by ID
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// Update User
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"firstName":  u.FirstName,
			"lastName":   u.LastName,
			"email":      u.Email,
			"password":   u.Password,
			"role":       u.Role,
			"photo":      u.Photo,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// To display updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// Delete User
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
