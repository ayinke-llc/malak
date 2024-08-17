package malak

// Swagger generation
//
//go:generate swag init -g swagger.go
//
//
// Mocks generation
//go:generate mockgen -source=internal/pkg/socialauth/social.go -destination=internal/pkg/socialauth/mocks/social.go -package=socialauth_mocks
//go:generate mockgen -source=user.go -destination=mocks/user.go -package=malak_mocks
