package waffle

import "fmt"

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalJSON(input []byte) error {
	_, err := fmt.Sscanf(string(input), "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return err
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Config struct {
	Pkg        string     `json:"pkg"`
	Maintainer Maintainer `json:"maintainer"`
	URL        string     `json:"url"`
	Repo       string     `json:"repo"`
	Org        string     `json:"org"`
	Version    Version    `json:"version"`
}
