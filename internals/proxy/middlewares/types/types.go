// Pretty useless, but go requires separate package
// since import cycles are not allowed

package types

type DataAlias struct {
	Alias string `koanf:"alias"`
	Score int    `koanf:"score"`
}
