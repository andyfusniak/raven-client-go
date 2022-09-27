package http

import "time"

type container struct {
	Data []interface{} `json:"data"`
}

// TemplateActions name value pairs for template parameters.
type TemplateActions map[string]interface{}

// Template resource object.
type Template struct {
	ID                     string                 `json:"id"`
	ProjectID              string                 `json:"projectId"`
	GroupID                string                 `json:"groupId"`
	Txt                    string                 `json:"txt"`
	HTML                   string                 `json:"html"`
	TxtDigest              string                 `json:"txtDigest"`
	HTMLDigest             string                 `json:"htmlDigest"`
	TxtActions             map[string]interface{} `json:"txtActions"`
	HTMLActions            map[string]interface{} `json:"htmlActions"`
	TxtTemplateCompiledOK  bool                   `json:"txtTemplateCompiledOk"`
	HTMLTemplateCompiledOK bool                   `json:"htmlTemplateCompiledOk"`
	CreatedAt              time.Time              `json:"createdAt"`
	ModifiedAt             time.Time              `json:"modifiedAt"`
}
