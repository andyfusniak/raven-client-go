package http

import (
	"fmt"
	"time"
)

const (
	ErrCodeUserNotFound           = "users/user-not-found"
	ErrCodeUserIDAttribInvalid    = "users/user-id-invalid"
	ErrCodeProjectNotFound        = "projects/project-not-found"
	ErrCodeProjectSlugExist       = "projects/project-slug-exists"
	ErrCodeProjectSlugInvalid     = "projects/project-slug-invalid"
	ErrCodeGroupIDInvalid         = "groups/group-id-invalid"
	ErrCodeGroupNotFound          = "groups/group-not-found"
	ErrCodeGroupContainsTemplates = "groups/group-contains-templates"
	ErrCodeTransportNotFound      = "transports/transport-not-found"
	ErrCodeTemplateNotFound       = "templates/template-not-found"
	ErrCodeTemplateExists         = "templates/template-exists"
	ErrCodeMailTemplateParse      = "email/email-template-parse-failure"
	ErrCodeMailTemplateExecute    = "email/email-template-execute-failure"
	ErrCodeBadRequest             = "bad-request"
)

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

// CreateTemplateParams parameters to create a new template.
type CreateTemplateParams struct {
	ID      string
	GroupID string
	Txt     string
	HTML    string
}

type createTemplateRequest struct {
	GroupID string `json:"groupId,omitempty"`
	Txt     string `json:"txt"`
	HTML    string `json:"html"`
}

// APIError standard response format for Raven Mailer errors.
type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error string representation of an APIError.
func (e *APIError) Error() string {
	return fmt.Sprintf("Status: %d Code: %s Message: %s",
		e.Status, e.Code, e.Message)
}
