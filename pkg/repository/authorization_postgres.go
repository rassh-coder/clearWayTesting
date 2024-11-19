package repository

import (
	"clearWayTest/pkg/models"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
)

type AuthRepository struct {
	db *pgx.Conn
}

func NewAuthorization(db *pgx.Conn) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByLogin(login string) (*models.User, error) {
	var user models.User

	query := `
 		SELECT * FROM users WHERE login = @login
	`
	args := pgx.NamedArgs{
		"login": login,
	}

	row := r.db.QueryRow(context.Background(), query, args)

	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) SaveSession(session *models.Session) (string, error) {
	var token string
	var lastSession models.Session
	tx, err := r.db.BeginTx(context.TODO(), pgx.TxOptions{})
	selectQuery := `
		SELECT DISTINCT on (uid) * FROM public.sessions
		where uid = @uid
		ORDER BY uid, created_at desc
	`
	selectArgs := pgx.NamedArgs{
		"uid": session.UID,
	}

	selectRow := tx.QueryRow(context.TODO(), selectQuery, selectArgs)
	err = selectRow.Scan(&lastSession.ID, &lastSession.UID, &lastSession.CreatedAt, &lastSession.Ip, &lastSession.IsActive)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(context.TODO())
			if err != nil {
				log.Printf("can't rollback transaction: %s", err)
			}
		} else {
			err = tx.Commit(context.TODO())
			if err != nil {
				log.Printf("can't commit transaction: %s", err)
			}
		}
	}()
	if lastSession.ID != "" {
		updateQuery := `
			UPDATE sessions SET is_active = false
			WHERE sessions.id = @sessionId
		`
		updateArgs := pgx.NamedArgs{
			"sessionId": lastSession.ID,
		}
		_, err = tx.Exec(context.Background(), updateQuery, updateArgs)
		if err != nil {
			return "", err
		}
	}

	insertQuery := `
		INSERT INTO sessions (uid, ip, is_active) VALUES(@uid, @ip, true)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"uid": session.UID,
		"ip":  session.Ip,
	}
	row := tx.QueryRow(context.TODO(), insertQuery, args)

	err = row.Scan(&token)
	if err != nil {
		log.Printf("error save session: %s", err)
		return "", err
	}

	return token, nil
}

func (r *AuthRepository) GetSession(sessionId string) (*models.Session, error) {
	var session models.Session

	query := `
		SELECT * FROM sessions WHERE id = @sessionId		
	`
	args := pgx.NamedArgs{
		"sessionId": sessionId,
	}

	row := r.db.QueryRow(context.Background(), query, args)

	err := row.Scan(&session.ID, &session.UID, &session.CreatedAt, &session.Ip, &session.IsActive)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
