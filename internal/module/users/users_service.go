package users

type UsersService struct {
	repo *UsersRepository
}

func NewUsersService (repo *UsersRepository) *UsersService{
	return &UsersService {
		repo: repo, 
	}
}



