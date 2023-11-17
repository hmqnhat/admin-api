package usecase

import (
	"admin-api/auth"
	"admin-api/entity"
	"admin-api/handler/dto"
	"admin-api/repository"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	AdminUsecase interface {
		CreateAdminUser(req dto.CreateAdminUserRequest) error
		SignIn(req dto.LoginRequest) (res *dto.LoginResponse, err error)
		RefreshToken(req dto.RefreshTokenRequest) (res *dto.RefreshTokenResponse, err error)
	}
	adminUsecase struct {
		repoAdmin repository.AdminUserRepository
	}
)

func NewAdminUsecase() AdminUsecase {
	return &adminUsecase{
		repoAdmin: repository.NewAdminUserRepository(),
	}
}

func (uc *adminUsecase) CreateAdminUser(req dto.CreateAdminUserRequest) error {
	if _, err := uc.repoAdmin.FindByEmail(req.Email); err == nil {
		return errors.New("email already exist")
	}

	hashPass, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = uc.repoAdmin.Create(&entity.AdminUser{
		Email:    req.Email,
		Password: hashPass,
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *adminUsecase) SignIn(req dto.LoginRequest) (res *dto.LoginResponse, err error) {
	adminUser, err := uc.repoAdmin.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("email does not exist")
		}
		return nil, err
	}

	if err := auth.CheckPasswordHash(adminUser.Password, req.Password); err != nil {
		return nil, errors.New("email or password is incorrect")
	}

	now := time.Now()
	claims := auth.NewClaims(strconv.Itoa(int(adminUser.ID)), now)

	accessToken, err := auth.Sign(&claims)
	if err != nil {
		return nil, err
	}

	adminUser.RefreshToken = auth.GenerateRefreshToken(now)

	if err := uc.repoAdmin.UpdateRefreshToken(adminUser.ID, adminUser.RefreshToken); err != nil {
		return nil, err
	}

	res = &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: adminUser.RefreshToken,
	}

	return res, err
}

func (uc *adminUsecase) RefreshToken(req dto.RefreshTokenRequest) (res *dto.RefreshTokenResponse, err error) {
	adminUser, err := uc.repoAdmin.FindByRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("invalid refresh token")
		}
		return nil, err
	}

	expiredTime, err := strconv.Atoi(strings.Split(adminUser.RefreshToken, "_")[1])
	if err != nil {
		return nil, err
	}

	if auth.IsExpired(int64(expiredTime)) {
		return nil, errors.New("refresh token has expired")
	}

	now := time.Now()
	claims := auth.NewClaims(strconv.Itoa(int(adminUser.ID)), now)

	accessToken, err := auth.Sign(&claims)
	if err != nil {
		return nil, err
	}

	res = &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	return res, nil
}
