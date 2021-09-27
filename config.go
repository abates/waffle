package waffle

type Config struct {
	Pkg        string `json:"pkg"`
	Maintainer string `json:"maintainer"`
	Repo       string `json:"repo"`
	Org        string `json:"org"`
}
