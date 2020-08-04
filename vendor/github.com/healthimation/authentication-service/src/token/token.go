package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// role claims
const (
	RolePatient = "patient"
	RoleAdmin   = "admin"
	RoleCS      = "customer_service"
	RoleRN      = "registered_nutritionist"
	RoleCoach   = "coach"
)

// GenerateToken generates a token string for the given claims
func GenerateToken(claims map[string]interface{}, signingKey []byte) (string, error) {
	return GenerateTokenWithDuration(claims, signingKey, time.Hour*48) // default to 2 days
}

// GenerateTokenWithDuration generates a token string for the given claims with a specific duration
func GenerateTokenWithDuration(claims map[string]interface{}, signingKey []byte, duration time.Duration) (string, error) {
	if signingKey == nil || len(signingKey) == 0 {
		return "", fmt.Errorf("signingKey must not be empty")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	token.Claims = claims
	token.Claims["exp"] = strconv.FormatInt(time.Now().Add(duration).UTC().Unix(), 10)
	// Sign and get the complete encoded token as a string
	return token.SignedString(signingKey)
}

// GenerateAdminToken generates a new token with the admin role.
func GenerateAdminToken(signingKey []byte) (string, error) {
	adminClaims := make(map[string]interface{})
	adminClaims["sub"] = "00000000-0000-0000-0000-000000000000"
	adminClaims["iss"] = "service"
	adminClaims["role"] = RoleAdmin
	adminClaims["vfd"] = true
	return GenerateToken(adminClaims, signingKey)
}

// VerifyToken checks if a token is valid, if so it returns its claims, if not it returns an error
func VerifyToken(tokenString string, signingKey []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token not valid")
	}

	if exp, ok := token.Claims["exp"]; ok {
		if expStr, ok := exp.(string); ok {
			expInt, err := strconv.ParseInt(expStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Token expiration not properly formatted")
			}
			if expInt >= time.Now().UTC().Unix() {
				return token.Claims, nil
			}
			return nil, fmt.Errorf("Token expired")
		}
		return nil, fmt.Errorf("Token expiration not properly formatted")
	}
	return nil, fmt.Errorf("Token missing expiration")
}

// VerifyAndExtract will verify a JWT and then extract and return the user ID and role from it's claims
func VerifyAndExtract(tokenString string, signingKey []byte) (userID string, role string, veried bool, err error) {
	c, err := VerifyToken(tokenString, signingKey)
	if err != nil {
		return "", "", false, err
	}

	u, subOk := c["sub"]
	r, roleOk := c["role"]
	v, verifiedOk := c["vfd"]

	if !subOk || !roleOk || !verifiedOk {
		return "", "", false, fmt.Errorf("Sub or role or verified not set")
	}

	userID, subOk = u.(string)
	role, roleOk = r.(string)
	if !subOk || !roleOk {
		return "", "", false, fmt.Errorf("Sub or role not a string")
	}

	verified, verifiedOk := v.(bool)
	if !verifiedOk {
		return "", "", false, fmt.Errorf("Verrified is not a bool")
	}

	return userID, role, verified, nil
}

// VerifyForPerms verifies the token and returns a permission object
func VerifyForPerms(tokenString string, signingKey []byte, logger Logger) (*Permission, error) {
	c, err := VerifyToken(tokenString, signingKey)
	if err != nil {
		return nil, err
	}

	// subscription, role, verrified are required
	u, subOk := c["sub"]
	r, roleOk := c["role"]
	v, verifiedOk := c["vfd"]

	if !subOk || !roleOk || !verifiedOk {
		return nil, fmt.Errorf("Sub, role, verified not set")
	}

	userID, subOk := u.(string)
	role, roleOk := r.(string)
	if !subOk || !roleOk {
		return nil, fmt.Errorf("Sub or role not a string")
	}

	verified, verifiedOk := v.(bool)
	if !verifiedOk {
		return nil, fmt.Errorf("Verrified is not a bool")
	}
	result := Permission{
		id:       userID,
		role:     role,
		logger:   logger,
		verified: verified,
	}

	if ua, ok := c["user_access"]; ok {
		if userAccess, ok := ua.(map[string]interface{}); ok {
			result.userAccess = userAccess
		}
	}
	return &result, nil
}
