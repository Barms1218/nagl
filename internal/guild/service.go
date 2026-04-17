package guild

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
	"unicode"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	paralellism = 2
	saltLen     = 16
	keyLen      = 32
)

type GuildService struct {
	store     *database.Store
	validator *validator.Validate
}

func NewGuildService(s *database.Store, v *validator.Validate, p *ecdsa.PrivateKey) *GuildService {
	return &GuildService{
		store:     s,
		validator: v,
	}
}

func (s *GuildService) GetGuildByName(ctx context.Context, name string) (database.GetGuildByNameRow, error) {
	return s.store.GetGuildByName(ctx, name)
}

func (s *GuildService) GetGuildByID(ctx context.Context, id uuid.UUID) (database.GetGuildByIDRow, error) {
	return s.store.GetGuildByID(ctx, id)
}

func (s *GuildService) RegisterGuild(ctx context.Context, g GuildAuthRequest) (*jwt.Token, error) {
	hashedPassword, err := s.HashPassword(g.GuildKey)
	if err != nil {
		return nil, err
	}
	params := database.InsertGuildParams{
		Name:     g.GuildName,
		Password: hashedPassword,
	}

	guild, err := s.store.InsertGuild(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Error occurred during guild registration: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"guildID": guild.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token, nil

}

func (s *GuildService) EnterGuild(ctx context.Context, g GuildAuthRequest) (uuid.UUID, error) {
	guild, err := s.store.GetGuildByName(ctx, g.GuildName)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("No guild by that name: %w", err)
	}

	match, err := s.VerifyPassword(g.GuildKey, guild.Password)
	if err != nil || !match {
		return uuid.UUID{}, fmt.Errorf("Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"guildID": guild.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return guild.ID, nil
}

func (s *GuildService) VerifyPassword(password, phc string) (bool, error) {
	// Split on "$" — parts are: ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	parts := strings.Split(phc, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var mem, iters uint32
	var threads uint8
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &iters, &threads)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Re-hash the candidate password with the same params
	candidateHash := argon2.IDKey([]byte(password), salt, iters, mem, threads, uint32(len(storedHash)))

	return subtle.ConstantTimeCompare(candidateHash, storedHash) == 1, nil
}

func (s *GuildService) HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, paralellism, keyLen)

	phc := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		paralellism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return phc, nil
}

func (s *GuildService) StrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
