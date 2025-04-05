package appservice

import (
	"gorm.io/gorm"

	"github.com/cxcnxl/go-crud/internal/auth_helpers"
	"github.com/cxcnxl/go-crud/internal/dto"
	"github.com/cxcnxl/go-crud/internal/models"
)

type AppService struct {
	db *gorm.DB
}

func NewAppService(db *gorm.DB) *AppService {
    return &AppService{ db };
}

func (self *AppService) CreateUser(data dto.CreateUserDto) (models.User, error) {
    user := models.User{
        Email: data.Email,
        Username: data.Username,
        PasswordHashed: auth_helpers.HashPassword(
            data.Password,
            auth_helpers.GenerateRandomSalt(),
        ),
    };

    if duplicate, _ := self.GetUserByEmail(data.Email); duplicate.ID != 0 {
        return duplicate, DuplicateUserEmailError{};
    }
    if duplicate, _ := self.GetUserByUsername(data.Username); duplicate.ID != 0 {
        return duplicate, DuplicateUserUsernameError{};
    }

    result := self.db.Create(&user);
    if result.Error != nil {
        return user, result.Error;
    }

    return user, nil;
}

func (self *AppService) GetUserByEmail(email string) (models.User, error) {
    user := models.User{
        Email: email,
    };

    result := self.db.Limit(1).Where(&user).First(&user);
    if result.Error != nil {
        return user, result.Error;
    }

    return user, nil;
}

func (self *AppService) GetUserByUsername(username string) (models.User, error) {
    user := models.User{
        Username: username,
    };

    result := self.db.Limit(1).Where(&user).First(&user);
    if result.Error != nil {
        return user, result.Error;
    }

    return user, nil;
}

func (self *AppService) GetUserById(id uint) (models.User, error) {
    user := models.User{
        ID: id,
    };

    result := self.db.Limit(1).Where(&user).First(&user);
    if result.Error != nil {
        return user, result.Error;
    }

    return user, nil;
}

type DuplicateUserEmailError struct {}
func (self DuplicateUserEmailError) Error() string {
    return "duplicate_user_email";
}

type DuplicateUserUsernameError struct {}
func (self DuplicateUserUsernameError) Error() string {
    return "duplicate_user_username";
}
