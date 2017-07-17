package zendesk

import (
	"fmt"
	"time"
	"net/url"
	"net/http"
	"encoding/json"
)

const ( 
	urlTemplate = "https://%s:%s@%s.zendesk.com/api/v2/%s"
	searchTemplate = "search.json?query=%s"
)

// Used to construct URL for Zendesk Core API
// OAuth currently not supported, patches are welcome
type Config struct {
	Subdomain	string
	User		string
	Password	string
}

type client struct {
	httpClient	*http.Client
	baseURL		*url.URL
	baseRequest	*http.Request
}

// Container for seach results
type searchResponse struct {
	Count			int				`json: "count"`
	Type			int				`json: "-"`
	Next_page		string			`json:"next_page"`
	Results		    json.RawMessage	`json: "results"`
}

// All tickets are assigned to and from groups
// https://developer.zendesk.com/rest_api/docs/core/groups for more information
type Group struct {
	ID			uint64		`json: "id"`
	Name		string		`json: "name"`	
	Created_at	*time.Time	`json: "created_at"`
	Updated_at	*time.Time	`json: "updated_at"`
}

// All end-users map to an Organizational unit, usually identified by domain_name
// https://developer.zendesk.com/rest_api/docs/core/organizations
type Organization struct {
	ID					uint64		`json: "id"`
	Name				string		`json: "name"`
	Created_at			*time.Time	`json: "created_at"`
	Updated_at			*time.Time	`json: "updated_at"`
	Domain_names		[]string	`json: "domain_names"`
	Group_id			uint64		`json: "group_id"`
	Organization_fields	[]string	`json: "organization_fields"`
}

type Users struct {
	ID				uint64				`json: "id"`
	Email			string				`json: "email"`
	Name			string				`json: "name"`
	Created_at		*time.Time			`json: "created_at"`
	Last_login		*time.Time			`json: "last_login_at"`
	Organization_id	uint64				`json: "organization_id"`
	Time_zone		string				`json: "time_zone"`
	Updated_at		*time.Time			`json: "updated_at"`
	User_fields		map[string]string	`json: "user_felds"`
}

// Returns zendesk handler, you can provide your own http client or use the default 
func Open(conf *Config) (*client, error) {
	baseURL, err := url.Parse(fmt.Sprintf("https://%s:%s@%s.zendesk.com/api/v2/",
		conf.User, conf.Password, conf.Subdomain))

	baseRequest := &http.Request{
		URL: baseURL,
	}

	return &client{
		http.DefaultClient, 
		baseURL,
		baseRequest,
	}, err
}

// Creates Search request, returns top-level data and results json blob
// https://developer.zendesk.com/rest_api/docs/core/search
func (c *client) newSearchRequest(qry string) (*searchResponse) {
	u, _ := url.Parse(fmt.Sprintf(searchTemplate, qry))
	c.baseRequest.URL = c.baseURL.ResolveReference(u)
	resp, _ := c.httpClient.Do(c.baseRequest)
	
	var res searchResponse 
	json.NewDecoder(resp.Body).Decode(&res)
	
	return &res
}

// Returns slice of all groups 
func (c *client) Groups() ([]Group) {
	res := c.newSearchRequest("type:group")

	// Decode results 
	groups := make([]Group, res.Count)
	json.Unmarshal(res.Results, &groups)
	
	return groups
}

//
func (c *client) Organizations() ([]Organization) {
	res := c.newSearchRequest("type:organization")	

    organizations := make([]Organization, res.Count)
    json.Unmarshal(res.Results, &organizations)

	return organizations
}

// Returns slice of all users
func (c *client) Users() ([]Users) {
	res := c.newSearchRequest("type:user")

	users := make([]Users, res.Count)
	json.Unmarshal(res.Results, &users)

	return users
}
