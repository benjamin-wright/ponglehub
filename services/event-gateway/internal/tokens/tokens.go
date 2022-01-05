package tokens

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Subject string
	Kind    string
}

type Tokens struct {
	key   []byte
	redis *redis.Client
}

func New(keyfile string, redisUrl string) (*Tokens, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %+v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	tokens := Tokens{
		key:   key,
		redis: rdb,
	}

	return &tokens, nil
}

func (t *Tokens) DeleteToken(id string, kind string) error {
	key := fmt.Sprintf("%s.%s", id, kind)
	err := t.redis.Del(context.Background(), key).Err()

	if err != nil {
		return fmt.Errorf("error deleting token %s: %+v", key, err)
	}

	return nil
}

func (t *Tokens) GetToken(id string, kind string) (string, error) {
	key := fmt.Sprintf("%s.%s", id, kind)
	value, err := t.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch token: %+v", err)
	}

	return value, nil
}

func (t *Tokens) NewToken(id string, kind string, expiration time.Duration) (string, error) {
	key := fmt.Sprintf("%s.%s", id, kind)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Subject": id,
			"Kind":    kind,
		},
	)

	tokenString, err := token.SignedString(t.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %+v", err)
	}

	err = t.redis.Set(context.Background(), key, tokenString, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save token: %+v", err)
	}

	return tokenString, nil
}

func (t *Tokens) Parse(token string) (Claims, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.key, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("failed to parse token: %+v", err)
	}

	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		return Claims{
			Subject: claims["Subject"].(string),
			Kind:    claims["Kind"].(string),
		}, nil

	} else {
		return Claims{}, fmt.Errorf("invalid claims in parsed token")
	}
}

func (t *Tokens) AddPasswordHash(id string, password string) error {
	key := fmt.Sprintf("%s.password", id)

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %+v", err)
	}
	hash := string(bytes)

	err = t.redis.Set(context.Background(), key, hash, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to send hashed password to redis: %+v", err)
	}

	return nil
}

func (t *Tokens) CheckPassword(id string, password string) (bool, error) {
	key := fmt.Sprintf("%s.password", id)

	hash, err := t.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to fetch password: %+v", err)
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}
