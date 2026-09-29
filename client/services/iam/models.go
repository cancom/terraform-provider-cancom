package client_iam

import "github.com/golang-jwt/jwt/v5"

type UserCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

type UserCreateResponse struct {
	Principal string `json:"principal"`
	Tenant    string `json:"tenant"`
	Message   string `json:"message"`
}

type UserUpdateRequest struct {
	Description string `json:"description"`
	Group       string `json:"group"`
}

type User struct {
	Principal   string `json:"principal"`
	Description string `json:"description"`
	Group       string `json:"group"`
	Tenant      string `json:"tenant"`
	CreatedAt   string `json:"createdAt"`
	CreatedBy   string `json:"createdBy"`
}

type ServiceUserCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

type ServiceUserCreateResponse struct {
	Principal   string `json:"principal"`
	Tenant      string `json:"tenant"`
	Session     string `json:"session"`
	SessionHash string `json:"session_hash"`
	Expiration  int    `json:"expiration"`
	Message     string `json:"message"`
}

type ServiceUserUpdateRequest struct {
	Description string `json:"description"`
	Group       string `json:"group"`
}

type ServiceUser struct {
	Tenant      string `json:"tenant"`
	CreatedAt   string `json:"createdAt"`
	Principal   string `json:"principal"`
	Description string `json:"description"`
	Group       string `json:"group"`
	CreatedBy   string `json:"createdBy"`
	SessionHash string `json:"session_hash"`
}

type RoleCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

type RoleCreateResponse struct {
	Principal string `json:"principal"`
	Tenant    string `json:"tenant"`
	Message   string `json:"message"`
}

type RoleUpdateRequest struct {
	Description string `json:"description"`
	Group       string `json:"group"`
}

type Role struct {
	Tenant      string `json:"tenant"`
	CreatedAt   string `json:"createdAt"`
	Principal   string `json:"principal"`
	Description string `json:"description"`
	Group       string `json:"group"`
	CreatedBy   string `json:"createdBy"`
}

type AssumeRoleRequest struct {
	Role string `json:"role"`
}

type AssumeRoleResponse struct {
	Jwt     string `json:"jwt"`
	Message string `json:"message"`
}

type SessionCreateRequest struct {
	ServiceUser string `json:"serviceUser"`
	Comment     string `json:"comment"`
}

type SessionClaims struct {
	jwt.RegisteredClaims
	SessionID     string `json:"sessionId"`
	Tenant        string `json:"tenant"`
	Subaccount    string `json:"subaccount"`
	PrincipalType string `json:"principalType"`
	PrincipalName string `json:"principalName"`
	PrincipalCRN  string `json:"principalCRN"`
	Typ           string `json:"typ"`
}

type SessionCreateResponse struct {
	Jwt          string         `json:"jwt"`
	SessionID    string         `json:"sessionId,omitempty"`
	PrincipalCRN string         `json:"principalCRN,omitempty"`
	ExpiresAt    int64          `json:"expiresAt,omitempty"`
	Claims       *SessionClaims `json:"claims,omitempty"`
}

type SessionUpdateRequest struct {
	Comment string `json:"comment"`
}

type Session struct {
	Enabled      bool   `json:"enabled"`
	LastIssuedAt string `json:"lastIssuedAt"`
	Comment      string `json:"comment"`
	TTL          string `json:"ttl"`
	Principal    string `json:"principal"`
	SessionID    string `json:"sessionId"`
}

type SessionDeleteResponse struct {
	Message string `json:"message"`
}
