package services

import (
	"errors"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"math/rand"
	"strconv"
)

type UserService struct {
	repo                    ports.UserRepository
	addressRepo             ports.AddressRepository
	verificationCodeManager ports.VerificationCodeManager
	passwordManager         ports.PasswordManager
	jwtManager              ports.JWTManager
}

func NewUserService(
	repo ports.UserRepository, addressRepo ports.AddressRepository,
	verificationCodeManager ports.VerificationCodeManager,
	passwordManager ports.PasswordManager, jwtManager ports.JWTManager,
) UserService {
	return UserService{
		repo:                    repo,
		addressRepo:             addressRepo,
		verificationCodeManager: verificationCodeManager,
		passwordManager:         passwordManager,
		jwtManager:              jwtManager,
	}
}

func (s *UserService) List(
	filters map[string]string, limit int, offset int,
) ([]users.User, int) {
	return s.repo.List(filters, limit, offset)
}

func (s *UserService) GetByID(ID users.UserID) (users.User, bool) {
	user, ok := s.repo.GetByID(ID)
	if ok == true {
		user.Addresses = s.addressRepo.List(ID)
	}
	return user, ok
}

func (s *UserService) Login(email string, password string) (ports.Token, error) {
	user, ok := s.repo.GetByEmail(email)
	if ok == false {
		return ports.Token{}, ports.InvalidCredentials
	}

	encryptedPassword, ok := s.repo.GetPassword(user.ID)
	if ok == false {
		return ports.Token{}, ports.InvalidCredentials
	}

	ok, err := s.passwordManager.Verify(password, encryptedPassword)
	if err != nil || ok == false {
		return ports.Token{}, ports.InvalidCredentials
	}

	token, err := s.jwtManager.Create(user)
	if err != nil || ok == false {
		return ports.Token{}, errors.New("error to create JWT")
	}

	return token, nil
}

func (s *UserService) Add(
	name string, email string, password string, phone string, isActive bool, scopes []users.ScopeName,
) (users.User, error) {
	hashPassword, err := s.passwordManager.Encrypt(password)
	if err != nil {
		return users.User{}, err
	}
	return s.repo.Add(name, email, hashPassword, phone, isActive, scopes)
}

func (s *UserService) Update(
	ID users.UserID, name string, email string, phone string, scopes []users.ScopeName,
) (users.User, error) {
	return s.repo.Update(ID, name, email, phone, scopes)
}

func (s *UserService) ChangePassword(ID users.UserID, currentPassword string, newPassword string) error {
	encryptedPassword, ok := s.repo.GetPassword(ID)
	if ok == false {
		return ports.UserDoesNotExists
	}
	ok, err := s.passwordManager.Verify(currentPassword, encryptedPassword)
	if err != nil {
		return err
	}
	if ok == false {
		return ports.InvalidCredentials
	}

	newPass, err := s.passwordManager.Encrypt(newPassword)
	if err != nil {
		return err
	}
	err = s.repo.ChangePassword(ID, newPass)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Deactivate(ID users.UserID) error {
	return s.repo.Deactivate(ID)
}

func (s *UserService) ListAddresses(id users.UserID) []users.Address {
	return s.addressRepo.List(id)
}

func (s *UserService) AddAddress(
	ID users.UserID, department string, city string, address string, receiverPhone string,
	receiverName string,
) (users.Address, error) {
	createdAddress, err := s.addressRepo.Add(ID, department, city, address, receiverPhone, receiverName)
	if err != nil {
		return users.Address{}, err
	}

	return createdAddress, nil
}

func (s *UserService) DeleteAddress(ID users.UserID, addressID users.AddressID) error {
	err := s.addressRepo.Delete(addressID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateAddress(
	addressID users.AddressID, department string, city string, address string,
	receiverPhone string, receiverName string,
) (users.Address, error) {
	return s.addressRepo.Update(addressID, department, city, address, receiverPhone, receiverName)
}

func (s *UserService) SendAccountVerificationCode(user users.User) error {
	codeRangeMin := 1000
	codeRangeMax := 9999
	code := strconv.Itoa(rand.Intn(codeRangeMax-codeRangeMin) + codeRangeMin)

	err := s.verificationCodeManager.SendEmailToVerifyAccount(code, user.Email)
	if err != nil {
		return err
	}

	err = s.repo.SaveAccountVerificationCode(user.ID, code)
	return err
}

func (s *UserService) ValidateAccountVerificationCode(ID users.UserID, code string) bool {
	isValid := s.repo.ValidateAccountVerificationCode(ID, code)

	if isValid == true {
		ok := s.repo.Activate(ID)
		return ok
	}
	return false
}

func (s *UserService) SendRecoveryPasswordRequest(email string) error {
	codeRangeMin := 100000000
	codeRangeMax := 999999999
	code := strconv.Itoa(rand.Intn(codeRangeMax-codeRangeMin) + codeRangeMin)

	user, ok := s.repo.GetByEmail(email)
	if ok == false {
		return ports.UserDoesNotExists
	}

	err := s.verificationCodeManager.SendEmailToRecoverPassword(code, email)
	if err != nil {
		return err
	}

	err = s.repo.SaveRecoveryPasswordCode(user.ID, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RecoveryPassword(email string, newPassword string, code string) error {
	user, ok := s.repo.GetByEmail(email)
	if ok == false {
		return ports.UserDoesNotExists
	}
	ok = s.repo.ValidateRecoveryPasswordCode(user.ID, code)
	if ok == false {
		return ports.InvalidValidationCode
	}

	encryptedPassword, err := s.passwordManager.Encrypt(newPassword)
	if err != nil {
		return err
	}
	err = s.repo.ChangePassword(user.ID, encryptedPassword)
	if err != nil {
		return err
	}

	return nil
}
