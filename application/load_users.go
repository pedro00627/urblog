package application

import (
	"bufio"
	"github.com/pedro00627/urblog/infrastructure/db"
	"log"
	"os"
	"strings"

	"github.com/pedro00627/urblog/domain"
)

//go:generate mockgen -destination=./mocks/mock_load_users.go -package=mocks github.com/pedro00627/urblog/application LoadUsers
type LoadUsers interface {
	Execute(filePath string) ([]domain.User, error)
}
type LoadUsersUseCase struct {
	userRepo db.UserRepository
}

func NewLoadUsersUseCase(userRepo db.UserRepository) *LoadUsersUseCase {
	return &LoadUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *LoadUsersUseCase) Execute(filePath string) ([]domain.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error loading file: %v", err)
		return nil, err
	}
	defer file.Close()

	users, err := uc.parseUsers(file)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *LoadUsersUseCase) parseUsers(file *os.File) ([]domain.User, error) {
	var users []domain.User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		user, err := uc.parseUser(line)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, *user)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *LoadUsersUseCase) parseUser(line string) (*domain.User, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return nil, nil
	}

	_, err := uc.userRepo.FindByName(parts[0])
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}

	user := domain.NewUser(parts[0], parts[0])
	for _, following := range parts[1:] {
		user.Following[following] = true
	}

	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}
