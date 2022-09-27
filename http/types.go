package http

import "time"

// Project resource.
type Project struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

// Transport type.
type Transport struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	Name         string    `json:"name"`
	Host         string    `json:"host"`
	Port         int       `json:"port"`
	Username     string    `json:"username"`
	EmailFrom    string    `json:"emailFrom"`
	EmailReplyTo string    `json:"emailReplyTo"`
	IsActive     bool      `json:"isActiveTransport"`
	CreatedAt    time.Time `json:"createdAt"`
	ModifiedAt   time.Time `json:"modifiedAt"`
}

// Group resource.
type Group struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"projectId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

// TemplateActions name value pairs for template parameters.
type TemplateActions map[string]interface{}

// TemplateParseError type.
type TemplateParseError struct {
	TemplateName string `json:"templateID,omitempty"`
	LineNumber   string `json:"lineNum,omitempty"`
	Msg          string `json:"msg,omitempty"`
}

// Template resource.
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
	ErrorsTxt              TemplateParseError     `json:"txtErrors"`
	ErrorsHTML             TemplateParseError     `json:"htmlErrors"`
	CreatedAt              time.Time              `json:"createdAt"`
	ModifiedAt             time.Time              `json:"modifiedAt"`
}
