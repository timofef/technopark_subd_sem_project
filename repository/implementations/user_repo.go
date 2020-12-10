package implementations

import (
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
	tx,_ := u.db.Begin()

	user := models.User{}
	if err := tx.QueryRow("getUserProfileQuery", &nickname).
		Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
		return models.User{}, models.UserNotExists
	}

	return user, nil
}

func (u *UserRepo) UpdateUserByNickname(newInfo *models.User, nickname string) (models.User, error) {
	tx, _ := u.db.Begin()
	defer tx.Commit()

	user := models.User{}

	if err := tx.QueryRow("update_user",  newInfo.Email, newInfo.Fullname, &nickname, newInfo.About).
		Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
		if _, ok := err.(pgx.PgError); ok {
			return models.User{}, models.UserConflict
		}
		return models.User{}, models.UserNotExists
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
		"SELECT email, fullname, nickname, about "+
			"FROM users"+
			"WHERE nickname = $1",
	)
	if err != nil {
		return err
	}

	// UPDATE EXISTING USER BY NICKNAME
	_, err = u.db.Prepare("update_user",
		"UPDATE users SET " +
		"email = COALESCE($2, users.email), " +
		"fullname = COALESCE($3, users.fullname), " +
		"about = COALESCE($1, users.about) " +
		"WHERE nickname = $4 " +
		"RETURNING email, fullname, nickname, about",
	)
	if err != nil {
		return err
	}

	return nil
}
