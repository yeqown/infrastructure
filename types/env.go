package types

// Envrion .
type Envrion string

const (
	// DevEnv defination .
	DevEnv Envrion = "dev"
	// ProdEnv defination .
	ProdEnv Envrion = "prod"
	// SimuEnv defination.
	SimuEnv Envrion = "simu"
)

func (env Envrion) String() string {
	return string(env)
}

// ParseEnvrion srting to Envrion
func ParseEnvrion(s string) Envrion {
	switch s {
	case "":
		return DevEnv
	case "dev":
		return DevEnv
	case "prod":
		return ProdEnv
	case "simu":
		return SimuEnv
	default:
		// not supported env, default set to DevEnv
		return Envrion(s)
	}
}
