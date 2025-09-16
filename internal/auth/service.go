package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	secret []byte
}

type User struct {
	ID    string
	Email string
	Hash  []byte
}

func New(secret string) *Service {
	return &Service{secret: []byte(secret)}
}

func (s *Service) CreateUser(ctx context.Context, db *pgxpool.Pool, email, password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	id := uuid.NewString()
	_, err = db.Exec(ctx, `insert into users(id,email,password_hash) values($1,$2,$3)`, id, email, h)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) VerifyUser(ctx context.Context, db *pgxpool.Pool, email, password string) (User, error) {
	var u User
	if err := db.QueryRow(ctx, `select id,email,password_hash from users where email=$1`, email).Scan(&u.ID, &u.Email, &u.Hash); err != nil {
		return User{}, err
	}
	if bcrypt.CompareHashAndPassword(u.Hash, []byte(password)) != nil {
		return User{}, errors.New("invalid")
	}
	return u, nil
}

func (s *Service) Token(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Service) FromRequest(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("no auth")
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || subtle.ConstantTimeCompare([]byte(parts[0]), []byte("Bearer")) != 1 {
		return "", errors.New("bad auth")
	}
	tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) { return s.secret, nil })
	if err != nil || !tok.Valid {
		return "", errors.New("invalid")
	}
	m, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid")
	}
	sub, _ := m["sub"].(string)
	if sub == "" {
		return "", errors.New("invalid")
	}
	return sub, nil
}
