package arcade

type Config struct {
	URL    string `conf:"URL"    default:"https://api.arcade.dev/v1/tools/execute"`
	APIKey string `conf:"APIKEY"`
}
