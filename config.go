package waffle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi3"
)

// VersionList is a list that is sortable
// by semantic naming
type VersionList []Version

// Len of the version list.  This satisfies one of the
// requirements of the sort.Interface interface
func (vl VersionList) Len() int { return len(vl) }

// Less determines if version[i] is less than version[j]
// This satisfies one of the requirements of the sort.Interface
// interface
func (vl VersionList) Less(i, j int) bool {
	if vl[i].Major > vl[j].Major {
		return false
	}

	if vl[i].Minor > vl[j].Minor {
		return false
	}

	if vl[i].Patch > vl[j].Patch {
		return false
	}
	return true
}

// Swap will swap the versions at indices i and j.  This
// satisfies one of the requirements of the sort.Interface interface
func (vl VersionList) Swap(i, j int) { vl[i], vl[j] = vl[j], vl[i] }

// Version is a struct representation of a semantic version
type Version struct {
	// Major release number
	Major int
	// Minor release number
	Minor int
	// Patch release number
	Patch int
}

// String will convert the Version so a string in the form
// Major.Minor.Patch ie 1.0.0
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// MarshalJSON will convert the Version to a JSON string
// such as "1.0.0"
func (v *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// Set will parse the string and populate
// the appropriate Version fields
func (v *Version) Set(str string) error {
	_, err := fmt.Sscanf(str, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return err
}

// UnmarshalJSON will parse a JSON string and populate
// the fields Major, Minor and Parse.  The only valid
// input is "x.y.z" where x is the major number (integer)
// y is the minor number (integer) and z is the patch
// number (also integer).
func (v *Version) UnmarshalJSON(input []byte) error {
	str := ""
	err := json.Unmarshal(input, &str)
	if err == nil {
		err = v.Set(str)
	}
	return err
}

// Maintainer contains information about the maintainer of a project
type Maintainer struct {
	// Name is some identifying information about the individual
	Name string `json:"name"`
	// Email should contain a contact email for the project
	Email string `json:"email"`
	// Org indicates if an organization is responsible for the contact
	Org string `json:"org"`
}

type Module struct {
	// Path is the projects go module path
	Path string `json:"path"`

	// Version is the semantic version for the project
	Version Version `json:"version"`
}

const (
	// DefConfigFile is the default filename for the project config file
	DefConfigFile = "project.json"
	// DefAPIFile is the default filename for the openapi/swagger file for the project
	DefAPIFile = "openapi.json"

	// OpenAPIVersion is the version of JSON that is written to the api config file
	OpenAPIVersion = "3.1.0"
)

// Config represents all the information about a
// waffle project
type Config struct {
	// Name is the title of the project
	Name string `json:"name"`

	// Desc is a short description of the project
	Desc string `json:"desc"`

	// Maintainer is contact information for the person who
	// is responsible for the project
	Maintainer Maintainer `json:"maintainer"`

	// URL is a link to the project webpage
	URL string `json:"url"`

	// Module includes information about the current code version
	// and repository
	Module Module `json:"mod"`

	apiConfig *openapi3.T // not exported so it's easer to marshal the config to json
}

func (c *Config) AddController(name string) error {
	return nil
}

func (c *Config) APIConfig() *openapi3.T {
	return c.apiConfig
}

func save(fileType, filename string, data interface{}) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(filename, content, 0644)
		if err != nil {
			err = fmt.Errorf("Failed to write %q: %w", filename, err)
		}
	} else {
		err = fmt.Errorf("Failed to marshal %s: %w", fileType, err)
	}
	return err
}

func (c *Config) Save(projectFile, apiFile string) error {
	err := save("config", projectFile, c)
	if err == nil {
		c.apiConfig.OpenAPI = OpenAPIVersion
		if c.apiConfig.Info == nil {
			c.apiConfig.Info = &openapi3.Info{}
		}
		c.apiConfig.Info.Title = c.Name
		c.apiConfig.Info.Description = c.Desc
		c.apiConfig.Info.Contact = &openapi3.Contact{
			Name:  c.Maintainer.Name,
			URL:   c.URL,
			Email: c.Maintainer.Email,
		}
		c.apiConfig.Info.Version = c.Module.Version.String()

		err = save("api config", apiFile, c.apiConfig)
	}
	return err
}

func (c *Config) SaveDef() error {
	return c.Save(DefConfigFile, DefAPIFile)
}

func (c *Config) LoadDef() error {
	return c.Load(DefConfigFile, DefAPIFile)
}

func (c *Config) Load(projectFile, apiFile string) error {
	content, err := ioutil.ReadFile(projectFile)
	if err == nil {
		err = json.Unmarshal(content, c)
	} else {
		err = fmt.Errorf("Failed to load %q: %w", projectFile, err)
	}

	if err == nil || errors.Is(err, fs.ErrNotExist) {
		var err2 error
		c.apiConfig, err2 = openapi3.NewLoader().LoadFromFile(apiFile)
		if err != nil {
			if errors.Is(err2, fs.ErrNotExist) {
				c.apiConfig = &openapi3.T{}
				err2 = nil
			} else {
				err2 = fmt.Errorf("Failed to load %q: %w", apiFile, err2)
			}
		}

		if err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}
