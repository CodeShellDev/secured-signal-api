// Pretty useless, but go requires seperate package
// since import cycles are not allowed

package middlewareTypes

type MessageAlias struct {
	Alias    string `koanf:"alias"`
	Score 	 int	`koanf:"score"`
}