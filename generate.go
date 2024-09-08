package malak

// Swagger generation
//
//go:generate swag init -g swagger.go  --requiredByDefault
//
//
// Mocks generation
//go:generate mockgen -source=internal/pkg/socialauth/social.go -destination=internal/pkg/socialauth/mocks/social.go -package=socialauth_mocks
//go:generate mockgen -source=internal/pkg/jwttoken/jwt.go -destination=internal/pkg/jwttoken/mocks/token.go
//go:generate mockgen -source=user.go -destination=mocks/user.go -package=malak_mocks
//go:generate mockgen -source=plan.go -destination=mocks/plan.go -package=malak_mocks
//go:generate mockgen -source=workspace.go -destination=mocks/workspace.go -package=malak_mocks
//go:generate mockgen -source=contact.go -destination=mocks/contact.go -package=malak_mocks
