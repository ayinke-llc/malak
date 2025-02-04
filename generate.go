package malak

// Swagger generation
//
//go:generate swag init --output swagger -g swagger.go
//
//switch swagger v2 to v3
//go:generate go run tools/v3gen/main.go
//
//
// Mocks generation
//go:generate mockgen -source=internal/pkg/socialauth/social.go -destination=internal/pkg/socialauth/mocks/social.go -package=socialauth_mocks
//go:generate mockgen -source=internal/pkg/jwttoken/jwt.go -destination=internal/pkg/jwttoken/mocks/token.go
//go:generate mockgen -source=internal/pkg/queue/queue.go -destination=mocks/queue.go -package=malak_mocks
//go:generate mockgen -source=internal/pkg/cache/cache.go -destination=mocks/cache.go -package=malak_mocks
//go:generate mockgen -source=user.go -destination=mocks/user.go -package=malak_mocks
//go:generate mockgen -source=plan.go -destination=mocks/plan.go -package=malak_mocks
//go:generate mockgen -source=workspace.go -destination=mocks/workspace.go -package=malak_mocks
//go:generate mockgen -source=contact.go -destination=mocks/contact.go -package=malak_mocks
//go:generate mockgen -source=contact_list.go -destination=mocks/contact_list.go -package=malak_mocks
//go:generate mockgen -source=update.go -destination=mocks/update.go -package=malak_mocks
//go:generate mockgen -source=reference.go -destination=mocks/reference.go -package=malak_mocks
//go:generate mockgen -source=deck.go -destination=mocks/deck.go -package=malak_mocks
//go:generate mockgen -source=uuid.go -destination=mocks/uuid.go -package=malak_mocks
//go:generate mockgen -source=share.go -destination=mocks/share.go -package=malak_mocks
//go:generate mockgen -source=preferences.go -destination=mocks/preferences.go -package=malak_mocks
//go:generate mockgen -source=integration.go -destination=mocks/integration.go -package=malak_mocks
//go:generate mockgen -source=internal/pkg/billing/billing.go -destination=mocks/billing.go -package=malak_mocks
//go:generate mockgen -source=internal/secret/secret.go -destination=mocks/secret.go -package=malak_mocks
