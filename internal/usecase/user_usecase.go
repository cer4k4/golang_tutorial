package usecase

import (
    "errors"
    "shop/internal/domain"
    "shop/internal/repository"
    "shop/pkg/jwt"
    "golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
    return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) Register(req *domain.UserRequest) (*domain.User, error) {
    // Check if user already exists
    _, err := u.userRepo.GetByEmail(req.Email)
    if err == nil {
        return nil, errors.New("user already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &domain.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     "user",
    }

    if err := u.userRepo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

func (u *UserUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
    user, err := u.userRepo.GetByEmail(req.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        return nil, err
    }

    return &domain.LoginResponse{
        Token: token,
        User:  *user,
    }, nil
}

func (u *UserUsecase) GetProfile(userID uint) (*domain.User, error) {
    return u.userRepo.GetByID(userID)
}