package repository

import (
	"errors"
	"strings"
	"timeslot-app/models"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
)

type UserRepoImplementation struct {
	db *pgx.Conn
}

func NewUserRepo(dbConn *pgx.Conn) UserRepo {
	return &UserRepoImplementation{
		db: dbConn,
	}
}

type UserRepo interface {
	Create(user models.User) error
	UserExists(userName string) (bool, error)
	Get(userName string) (models.User, error)
}

func (ur *UserRepoImplementation) Create(user models.User) error {

	//check if user with the same name already exists
	exists, err := ur.UserExists(user.Name)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return err
	}

	if exists {
		return errors.New("user with the given name already exists")
	}

	insertQuery := `INSERT INTO users (id, name) VALUES ($1, $2)`
	_, err = ur.db.Exec(insertQuery, user.ID, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepoImplementation) UserExists(userName string) (bool, error) {

	var count int
	err := ur.db.QueryRow("SELECT count(*) FROM users WHERE name = $1", userName).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, errors.New("user does not exist")
	}
	return true, nil
}

func (ur *UserRepoImplementation) Get(userName string) (models.User, error) {

	var name string
	var id uuid.UUID
	err := ur.db.QueryRow("SELECT * FROM users WHERE name = $1", userName).Scan(&id, &name)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:   id,
		Name: name,
	}, nil
}
