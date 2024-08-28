// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proto/api/user/v1/user.proto

package user

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on RegisterRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterRequestMultiError, or nil if none found.
func (m *RegisterRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Login

	// no validation rules for Password

	if len(errors) > 0 {
		return RegisterRequestMultiError(errors)
	}

	return nil
}

// RegisterRequestMultiError is an error wrapping multiple validation errors
// returned by RegisterRequest.ValidateAll() if the designated constraints
// aren't met.
type RegisterRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterRequestMultiError) AllErrors() []error { return m }

// RegisterRequestValidationError is the validation error returned by
// RegisterRequest.Validate if the designated constraints aren't met.
type RegisterRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterRequestValidationError) ErrorName() string { return "RegisterRequestValidationError" }

// Error satisfies the builtin error interface
func (e RegisterRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterRequestValidationError{}

// Validate checks the field values on RegisterResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterResponseMultiError, or nil if none found.
func (m *RegisterResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for Message

	// no validation rules for AuthToken

	if len(errors) > 0 {
		return RegisterResponseMultiError(errors)
	}

	return nil
}

// RegisterResponseMultiError is an error wrapping multiple validation errors
// returned by RegisterResponse.ValidateAll() if the designated constraints
// aren't met.
type RegisterResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterResponseMultiError) AllErrors() []error { return m }

// RegisterResponseValidationError is the validation error returned by
// RegisterResponse.Validate if the designated constraints aren't met.
type RegisterResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterResponseValidationError) ErrorName() string { return "RegisterResponseValidationError" }

// Error satisfies the builtin error interface
func (e RegisterResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterResponseValidationError{}

// Validate checks the field values on AuthenticateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AuthenticateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthenticateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuthenticateRequestMultiError, or nil if none found.
func (m *AuthenticateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthenticateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Login

	// no validation rules for Password

	if len(errors) > 0 {
		return AuthenticateRequestMultiError(errors)
	}

	return nil
}

// AuthenticateRequestMultiError is an error wrapping multiple validation
// errors returned by AuthenticateRequest.ValidateAll() if the designated
// constraints aren't met.
type AuthenticateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthenticateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthenticateRequestMultiError) AllErrors() []error { return m }

// AuthenticateRequestValidationError is the validation error returned by
// AuthenticateRequest.Validate if the designated constraints aren't met.
type AuthenticateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthenticateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthenticateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthenticateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthenticateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthenticateRequestValidationError) ErrorName() string {
	return "AuthenticateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e AuthenticateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthenticateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthenticateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthenticateRequestValidationError{}

// Validate checks the field values on AuthenticateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AuthenticateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthenticateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuthenticateResponseMultiError, or nil if none found.
func (m *AuthenticateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthenticateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for AuthToken

	// no validation rules for Message

	if len(errors) > 0 {
		return AuthenticateResponseMultiError(errors)
	}

	return nil
}

// AuthenticateResponseMultiError is an error wrapping multiple validation
// errors returned by AuthenticateResponse.ValidateAll() if the designated
// constraints aren't met.
type AuthenticateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthenticateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthenticateResponseMultiError) AllErrors() []error { return m }

// AuthenticateResponseValidationError is the validation error returned by
// AuthenticateResponse.Validate if the designated constraints aren't met.
type AuthenticateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthenticateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthenticateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthenticateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthenticateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthenticateResponseValidationError) ErrorName() string {
	return "AuthenticateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e AuthenticateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthenticateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthenticateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthenticateResponseValidationError{}

// Validate checks the field values on GetMetadataRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetMetadataRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetMetadataRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetMetadataRequestMultiError, or nil if none found.
func (m *GetMetadataRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetMetadataRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AuthToken

	if len(errors) > 0 {
		return GetMetadataRequestMultiError(errors)
	}

	return nil
}

// GetMetadataRequestMultiError is an error wrapping multiple validation errors
// returned by GetMetadataRequest.ValidateAll() if the designated constraints
// aren't met.
type GetMetadataRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetMetadataRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetMetadataRequestMultiError) AllErrors() []error { return m }

// GetMetadataRequestValidationError is the validation error returned by
// GetMetadataRequest.Validate if the designated constraints aren't met.
type GetMetadataRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMetadataRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMetadataRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMetadataRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMetadataRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMetadataRequestValidationError) ErrorName() string {
	return "GetMetadataRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetMetadataRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMetadataRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMetadataRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMetadataRequestValidationError{}

// Validate checks the field values on GetMetadataResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetMetadataResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetMetadataResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetMetadataResponseMultiError, or nil if none found.
func (m *GetMetadataResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetMetadataResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetMetadata() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetMetadataResponseValidationError{
						field:  fmt.Sprintf("Metadata[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetMetadataResponseValidationError{
						field:  fmt.Sprintf("Metadata[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetMetadataResponseValidationError{
					field:  fmt.Sprintf("Metadata[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetMetadataResponseMultiError(errors)
	}

	return nil
}

// GetMetadataResponseMultiError is an error wrapping multiple validation
// errors returned by GetMetadataResponse.ValidateAll() if the designated
// constraints aren't met.
type GetMetadataResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetMetadataResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetMetadataResponseMultiError) AllErrors() []error { return m }

// GetMetadataResponseValidationError is the validation error returned by
// GetMetadataResponse.Validate if the designated constraints aren't met.
type GetMetadataResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMetadataResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMetadataResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMetadataResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMetadataResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMetadataResponseValidationError) ErrorName() string {
	return "GetMetadataResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetMetadataResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMetadataResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMetadataResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMetadataResponseValidationError{}

// Validate checks the field values on CreateSessionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateSessionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateSessionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateSessionRequestMultiError, or nil if none found.
func (m *CreateSessionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateSessionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AuthToken

	// no validation rules for DeviceInfo

	if len(errors) > 0 {
		return CreateSessionRequestMultiError(errors)
	}

	return nil
}

// CreateSessionRequestMultiError is an error wrapping multiple validation
// errors returned by CreateSessionRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateSessionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateSessionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateSessionRequestMultiError) AllErrors() []error { return m }

// CreateSessionRequestValidationError is the validation error returned by
// CreateSessionRequest.Validate if the designated constraints aren't met.
type CreateSessionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateSessionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateSessionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateSessionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateSessionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateSessionRequestValidationError) ErrorName() string {
	return "CreateSessionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateSessionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateSessionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateSessionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateSessionRequestValidationError{}

// Validate checks the field values on CreateSessionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateSessionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateSessionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateSessionResponseMultiError, or nil if none found.
func (m *CreateSessionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateSessionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for SessionId

	// no validation rules for Message

	if len(errors) > 0 {
		return CreateSessionResponseMultiError(errors)
	}

	return nil
}

// CreateSessionResponseMultiError is an error wrapping multiple validation
// errors returned by CreateSessionResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateSessionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateSessionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateSessionResponseMultiError) AllErrors() []error { return m }

// CreateSessionResponseValidationError is the validation error returned by
// CreateSessionResponse.Validate if the designated constraints aren't met.
type CreateSessionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateSessionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateSessionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateSessionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateSessionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateSessionResponseValidationError) ErrorName() string {
	return "CreateSessionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateSessionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateSessionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateSessionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateSessionResponseValidationError{}

// Validate checks the field values on EndSessionRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *EndSessionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EndSessionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// EndSessionRequestMultiError, or nil if none found.
func (m *EndSessionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *EndSessionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for SessionId

	if len(errors) > 0 {
		return EndSessionRequestMultiError(errors)
	}

	return nil
}

// EndSessionRequestMultiError is an error wrapping multiple validation errors
// returned by EndSessionRequest.ValidateAll() if the designated constraints
// aren't met.
type EndSessionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EndSessionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EndSessionRequestMultiError) AllErrors() []error { return m }

// EndSessionRequestValidationError is the validation error returned by
// EndSessionRequest.Validate if the designated constraints aren't met.
type EndSessionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EndSessionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EndSessionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EndSessionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EndSessionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EndSessionRequestValidationError) ErrorName() string {
	return "EndSessionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e EndSessionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEndSessionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EndSessionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EndSessionRequestValidationError{}

// Validate checks the field values on EndSessionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *EndSessionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EndSessionResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// EndSessionResponseMultiError, or nil if none found.
func (m *EndSessionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *EndSessionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for Message

	if len(errors) > 0 {
		return EndSessionResponseMultiError(errors)
	}

	return nil
}

// EndSessionResponseMultiError is an error wrapping multiple validation errors
// returned by EndSessionResponse.ValidateAll() if the designated constraints
// aren't met.
type EndSessionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EndSessionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EndSessionResponseMultiError) AllErrors() []error { return m }

// EndSessionResponseValidationError is the validation error returned by
// EndSessionResponse.Validate if the designated constraints aren't met.
type EndSessionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EndSessionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EndSessionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EndSessionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EndSessionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EndSessionResponseValidationError) ErrorName() string {
	return "EndSessionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e EndSessionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEndSessionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EndSessionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EndSessionResponseValidationError{}

// Validate checks the field values on Metadata with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Metadata) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Metadata with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MetadataMultiError, or nil
// if none found.
func (m *Metadata) ValidateAll() error {
	return m.validate(true)
}

func (m *Metadata) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for MetadataId

	// no validation rules for MetadataKey

	// no validation rules for MetadataValue

	if len(errors) > 0 {
		return MetadataMultiError(errors)
	}

	return nil
}

// MetadataMultiError is an error wrapping multiple validation errors returned
// by Metadata.ValidateAll() if the designated constraints aren't met.
type MetadataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetadataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetadataMultiError) AllErrors() []error { return m }

// MetadataValidationError is the validation error returned by
// Metadata.Validate if the designated constraints aren't met.
type MetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataValidationError) ErrorName() string { return "MetadataValidationError" }

// Error satisfies the builtin error interface
func (e MetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataValidationError{}
