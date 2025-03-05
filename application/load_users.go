package application

import (
	"bufio"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/repositories"
)

type LoadUsersUseCase struct {
	userRepo repositories.UserRepository
}

func NewLoadUsersUseCase(userRepo repositories.UserRepository) *LoadUsersUseCase {
	return &LoadUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *LoadUsersUseCase) Execute(filePath string) ([]domain.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	users := []domain.User{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		// Validate if user already exists
		_, err := uc.userRepo.FindByName(parts[0])
		if err != nil && err != domain.ErrUserNotFound {
			return nil, err
		}

		user := domain.NewUser(uuid.New().String(), parts[0])
		// Agregar usuarios seguidos
		for _, following := range parts[1:] {
			user.Following[following] = true
		}

		savingErr := uc.userRepo.Save(user)
		if savingErr != nil {
			return nil, savingErr
		}

		users = append(users, *user)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
