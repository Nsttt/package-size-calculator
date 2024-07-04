package npm

import (
	"encoding/json"
	"fmt"
	"unicode"

	npm_version "github.com/aquasecurity/go-npm-version/pkg"
	"github.com/rs/zerolog/log"
)

type PackageJSON struct {
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Dependencies NormalDependencies `json:"dependencies"`
}

func (p PackageJSON) String() string {
	return p.AsDependency().String()
}

func (p PackageJSON) AsDependency() DependencyInfo {
	return DependencyInfo{
		Name:    p.Name,
		Version: p.Version,
	}
}

type NormalDependencies []Dependency

func (d *NormalDependencies) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for name, rawConstraint := range raw {
		l := log.With().Str("name", name).Str("rawConstraint", rawConstraint).Logger()

		l.Trace().Msg("Adding normal dependency")

		if len(rawConstraint) == 0 || unicode.IsLetter(rune(rawConstraint[0])) {
			// l.Warn().Msg("Skipping dependency with improper version constraint format")
			continue
		}

		c, err := npm_version.NewConstraints(rawConstraint)
		if err != nil {
			// l.Warn().Err(err).Msg("Failed to create constraints")
			continue
		}

		*d = append(*d, Dependency{
			Name:          name,
			RawConstraint: rawConstraint,
			Constraint:    c,
		})
	}

	return nil
}

type Dependency struct {
	Name          string
	RawConstraint string
	Constraint    npm_version.Constraints
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s %s", d.Name, d.RawConstraint)
}
