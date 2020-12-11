package implementations

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type UserRepo struct {
	db *pgx.ConnPool
}

func NewUserRepo(pool *pgx.ConnPool) interfaces.UserRepository {
	return &UserRepo{db: pool}
}

func (u *UserRepo) CreateUser(user *models.User, nickname string) (models.Users, error) {
	tx, err := u.db.Begin()
	defer tx.Commit()

	result, err := tx.Exec("insert_user", user.Email, user.Fullname, &nickname, user.About)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		existingUsers := models.Users{}

		rows, _ := tx.Query("get_users_by_nickname_or_email", user.Email, &nickname)

		for rows.Next() {
			existingUser := models.User{}
			_ = rows.Scan(&existingUser.Email, &existingUser.Fullname, &existingUser.Nickname, &existingUser.About)

			existingUsers = append(existingUsers, existingUser)
		}
		rows.Close()

		return existingUsers, models.UserExists
	}

	return nil, nil
}

func (u *UserRepo) GetUserByNickname(nickname string) (models.User, error) {
	tx, _ := u.db.Begin()

	user := models.User{}
	if err := tx.QueryRow("get_user_by_nickname", &nickname).
		Scan(&user.Email, &user.Fullname, &user.Nickname, &user.About); err != nil {
		return models.User{}, models.UserNotExists
	}

	return user, nil
}

func (u *UserRepo) UpdateUserByNickname(newInfo *models.UserUpdate, nickname string) (models.User, error) {
	tx, _ := u.db.Begin()
	defer tx.Commit()

	if newInfo.Email != "" {
		existingUser := models.User{}
		if err := tx.QueryRow("get_user_by_email", &newInfo.Email).
			Scan(&existingUser.Email, &existingUser.Fullname, &existingUser.Nickname, &existingUser.About); err == nil {
			return models.User{}, models.UserConflict
		}
	}

	user := models.User{}

	if err := tx.QueryRow("update_user",
		sql.NullString{String: newInfo.Email.String(), Valid: newInfo.Email != ""},
		sql.NullString{String: newInfo.Fullname, Valid: newInfo.Fullname != ""},
		&nickname,
		sql.NullString{String: newInfo.About, Valid: newInfo.About != ""}).
		Scan(&user.Email, &user.Fullname, &user.Nickname, &user.About); err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepo) PrepareStatements() error {
	// INSERT NEW USER
	_, err := u.db.Prepare("insert_user",
		"INSERT INTO users (email, fullname, nickname, about) "+
			"VALUES ($1, $2, $3, $4) "+
			"ON CONFLICT DO NOTHING ",
	)
	if err != nil {
		return err
	}

	// GET USERs BY NICKNAME OR EMAIL
	_, err = u.db.Prepare("get_users_by_nickname_or_email",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM users "+
			"WHERE email = $1 OR nickname = $2",
	)
	if err != nil {
		return err
	}

	// GET USER BY NICKNAME
	_, err = u.db.Prepare("get_user_by_nickname",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM users "+
			"WHERE nickname = $1",
	)
	if err != nil {
		return err
	}

	// GET USER BY EMAIL
	_, err = u.db.Prepare("get_user_by_email",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM users "+
			"WHERE email = $1",
	)
	if err != nil {
		return err
	}

	// UPDATE EXISTING USER BY NICKNAME
	_, err = u.db.Prepare("update_user",
		"UPDATE users SET "+
			"email = COALESCE($1, users.email), "+
			"fullname = COALESCE($2, users.fullname), "+
			"about = COALESCE($4, users.about) "+
			"WHERE nickname = $3 "+
			"RETURNING users.email, users.fullname, users.nickname, users.about ",
	)
	if err != nil {
		return err
	}

	return nil
}
