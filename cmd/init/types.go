package init

import "errors"

type AppType string

const (
	AppTypeAPI    AppType = "api"
	AppTypeWeb    AppType = "web"
	AppTypeMobile AppType = "mobile"
)

func (a AppType) String() string {
	return string(a)
}

func (a *AppType) Set(value string) error {
	switch value {
	case "api", "web", "mobile":
		*a = AppType(value)
		return nil
	default:
		return errors.New(`must be one of "api", "web", or "mobile"`)
	}
}

func (a *AppType) Type() string {
	return "app-type"
}

type AppTypeOption struct {
	Name  string
	Value AppType
}

func (a *AppType) Options() []AppTypeOption {
	return []AppTypeOption{
		{Name: "API", Value: AppTypeAPI},
		{Name: "Web", Value: AppTypeWeb},
		{Name: "Mobile", Value: AppTypeMobile},
	}
}
