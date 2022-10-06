package http

import (
	"fmt"
	"time"
)

const (
	// users

	// ErrCodeUserIDInvalid response error code.
	ErrCodeUserIDInvalid = "users/user-id-invalid"

	// ErrCodeUserNotFound response error code.
	ErrCodeUserNotFound = "users/user-not-found"

	// ErrCodeUserIDAttribInvalid response error code.
	ErrCodeUserIDAttribInvalid = "users/user-id-invalid"

	// projects

	// ErrCodeProjectIDInvalid response error code.
	ErrCodeProjectIDInvalid = "projects/project-id-invalid"

	// ErrCodeProjectNotFound response error code.
	ErrCodeProjectNotFound = "projects/project-not-found"

	// ErrCodeProjectExist response error code.
	ErrCodeProjectExist = "projects/project-exists"

	// transports

	// ErrCodeTransportIDInvalid response error code.
	ErrCodeTransportIDInvalid = "transports/transport-id-invalid"

	// ErrCodeTransportNotFound response error code.
	ErrCodeTransportNotFound = "transports/transport-not-found"

	// ErrCodeTransportCodeInvalid response error code.
	ErrCodeTransportCodeInvalid = "transports/transport-code-invalid"

	// ErrCodeActiveTransportNotFound response error code.
	ErrCodeActiveTransportNotFound = "mail/active-transport-not-found"

	// groups

	// ErrCodeGroupIDInvalid response error code.
	ErrCodeGroupIDInvalid = "groups/group-id-invalid"

	// ErrCodeGroupNotFound response error code.
	ErrCodeGroupNotFound = "groups/group-not-found"

	// ErrCodeGroupExists response error code.
	ErrCodeGroupExists = "groups/group-exists"

	// ErrCodeGroupContainsTemplates response error code.
	ErrCodeGroupContainsTemplates = "groups/group-contains-templates"

	// templates

	// ErrCodeTemplateIDInvalid response error code.
	ErrCodeTemplateIDInvalid = "templates/template-id-invalid"

	// ErrCodeTemplateNotFound response error code.
	ErrCodeTemplateNotFound = "templates/template-not-found"

	// ErrCodeTemplateExists response error code.
	ErrCodeTemplateExists = "templates/template-exists"

	// mail

	// ErrCodeMailIDInvalid response error code.
	ErrCodeMailIDInvalid = "mail/mail-id-invalid"

	// ErrCodeMailNotFound response error code.
	ErrCodeMailNotFound = "mail/mail-not-found"

	// ErrCodeMailTemplateParse response error code.
	ErrCodeMailTemplateParse = "mail/mail-template-parse-failure"

	// ErrCodeMailTemplateExecute response error code.
	ErrCodeMailTemplateExecute = "mail/mail-template-execute-failure"

	// general

	// ErrCodeBadRequest response error code.
	ErrCodeBadRequest = "bad-request"
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
	Active       bool      `json:"active"`
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
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"projectId"`
	GroupID     string                 `json:"groupId"`
	Txt         string                 `json:"txt"`
	HTML        string                 `json:"html"`
	TxtDigest   string                 `json:"txtDigest"`
	HTMLDigest  string                 `json:"htmlDigest"`
	TxtActions  map[string]interface{} `json:"txtActions"`
	HTMLActions map[string]interface{} `json:"htmlActions"`
	CreatedAt   time.Time              `json:"createdAt"`
	ModifiedAt  time.Time              `json:"modifiedAt"`
}

// CreateTemplateParams parameters to create a new template.
type CreateTemplateParams struct {
	ID        string
	ProjectID string
	GroupID   string
	Txt       string
	HTML      string
}

type createTemplateRequest struct {
	GroupID string `json:"groupId,omitempty"`
	Txt     string `json:"txt"`
	HTML    string `json:"html"`
}

// Mail resource.
type Mail struct {
	ID           string     `json:"id"`
	TemplateID   string     `json:"templateId"`
	ProjectID    string     `json:"projectId"`
	Status       string     `json:"status"`
	EmailTo      string     `json:"emailTo"`
	EmailFrom    string     `json:"emailFrom"`
	EmailReplyTo string     `json:"emailReplyTo"`
	Subject      string     `json:"subject"`
	CreatedAt    time.Time  `json:"createdAt"`
	SentAt       *time.Time `json:"sentAt"`
	ModifiedAt   time.Time  `json:"modifiedAt"`
}

// MailLog type.
type MailLog struct {
	ID        string                 `json:"id"`
	MailID    string                 `json:"mailId"`
	ProjectID string                 `json:"projectId"`
	Status    string                 `json:"status"`
	SMTPCode  int                    `json:"smtpCode"`
	Msg       string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"createdAt"`
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
