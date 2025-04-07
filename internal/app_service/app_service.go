package appservice

import (
	"context"

	"gorm.io/gorm"

	"github.com/cxcnxl/go-crud/internal/auth_helpers"
	"github.com/cxcnxl/go-crud/internal/dto"
	"github.com/cxcnxl/go-crud/internal/models"
	"github.com/cxcnxl/go-crud/internal/redis"
)

type AppService struct {
	db *gorm.DB
    redis *redis.RedisWrapper
    ctx context.Context
}

func NewAppService(
    db *gorm.DB,
    redis *redis.RedisWrapper,
) *AppService {
    return &AppService{
        db,
        redis,
        context.Background(),
    };
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

func (self *AppService) LoginUser(data dto.PostLoginDto) (models.User, error) {
    var user models.User;
    result := self.db.
        Where(models.User{Email: data.Username}).
        Or(models.User{Username: data.Username}).
        First(&user);

    loginBlocked := self.getLoginBlocked(user);
    if loginBlocked == true {
        return user, LoginBlockedError{};
    }

    if result.Error != nil {
        return user, result.Error;
    }

    passwordValid, err := auth_helpers.VerifyPass(
        data.Password,
        user.PasswordHashed,
    );
    if err != nil {
        return user, err;
    }

    if passwordValid == false {
        self.handleInvalidPasswordAttempt(user);
        return user, InvalidPasswordError{};
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

// TODO: move redis operations into redis service
//
// TODO: create const keys for redis
func (self *AppService) handleInvalidPasswordAttempt(
    user models.User,
) {
    loginAttempts, err := self.redis.GetLoginAttempts(user.ID);
    if err != nil {
        panic(err);
    }

    loginAttempts += 1;

    if loginAttempts >= 3 {
        self.setLoginBlocked(user);
    }

    err = self.redis.SetLoginAttempts(user.ID, loginAttempts);
    if err != nil {
        panic(err);
    }
}

func (self *AppService) setLoginBlocked(
    user models.User,
) {
    err := self.redis.SetLoginBlocked(user.ID, true);
    if err != nil {
        panic(err);
    }
}

func (self *AppService) getLoginBlocked(
    user models.User,
) bool {
    blocked, err := self.redis.GetLoginBlocked(user.ID);
    if err != nil {
        panic(err);
    }

    return blocked;
}

type DuplicateUserEmailError struct {}
func (self DuplicateUserEmailError) Error() string {
    return "duplicate_user_email";
}

type DuplicateUserUsernameError struct {}
func (self DuplicateUserUsernameError) Error() string {
    return "duplicate_user_username";
}

type InvalidPasswordError struct {}
func (self InvalidPasswordError) Error() string {
    return "invalid_password";
}

type LoginBlockedError struct {}
func (self LoginBlockedError) Error() string {
    return "login_blocked";
}
