// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: web/v1/organize.proto

package web

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

// Validate checks the field values on OrganizeDepartmentListRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizeDepartmentListRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizeDepartmentListRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// OrganizeDepartmentListRequestMultiError, or nil if none found.
func (m *OrganizeDepartmentListRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizeDepartmentListRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return OrganizeDepartmentListRequestMultiError(errors)
	}

	return nil
}

// OrganizeDepartmentListRequestMultiError is an error wrapping multiple
// validation errors returned by OrganizeDepartmentListRequest.ValidateAll()
// if the designated constraints aren't met.
type OrganizeDepartmentListRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizeDepartmentListRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizeDepartmentListRequestMultiError) AllErrors() []error { return m }

// OrganizeDepartmentListRequestValidationError is the validation error
// returned by OrganizeDepartmentListRequest.Validate if the designated
// constraints aren't met.
type OrganizeDepartmentListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizeDepartmentListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizeDepartmentListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizeDepartmentListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizeDepartmentListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizeDepartmentListRequestValidationError) ErrorName() string {
	return "OrganizeDepartmentListRequestValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizeDepartmentListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizeDepartmentListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizeDepartmentListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizeDepartmentListRequestValidationError{}

// Validate checks the field values on OrganizeDepartmentListResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizeDepartmentListResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizeDepartmentListResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// OrganizeDepartmentListResponseMultiError, or nil if none found.
func (m *OrganizeDepartmentListResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizeDepartmentListResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrganizeDepartmentListResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrganizeDepartmentListResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrganizeDepartmentListResponseValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OrganizeDepartmentListResponseMultiError(errors)
	}

	return nil
}

// OrganizeDepartmentListResponseMultiError is an error wrapping multiple
// validation errors returned by OrganizeDepartmentListResponse.ValidateAll()
// if the designated constraints aren't met.
type OrganizeDepartmentListResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizeDepartmentListResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizeDepartmentListResponseMultiError) AllErrors() []error { return m }

// OrganizeDepartmentListResponseValidationError is the validation error
// returned by OrganizeDepartmentListResponse.Validate if the designated
// constraints aren't met.
type OrganizeDepartmentListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizeDepartmentListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizeDepartmentListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizeDepartmentListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizeDepartmentListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizeDepartmentListResponseValidationError) ErrorName() string {
	return "OrganizeDepartmentListResponseValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizeDepartmentListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizeDepartmentListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizeDepartmentListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizeDepartmentListResponseValidationError{}

// Validate checks the field values on OrganizePersonnelListRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizePersonnelListRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizePersonnelListRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrganizePersonnelListRequestMultiError, or nil if none found.
func (m *OrganizePersonnelListRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizePersonnelListRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return OrganizePersonnelListRequestMultiError(errors)
	}

	return nil
}

// OrganizePersonnelListRequestMultiError is an error wrapping multiple
// validation errors returned by OrganizePersonnelListRequest.ValidateAll() if
// the designated constraints aren't met.
type OrganizePersonnelListRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizePersonnelListRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizePersonnelListRequestMultiError) AllErrors() []error { return m }

// OrganizePersonnelListRequestValidationError is the validation error returned
// by OrganizePersonnelListRequest.Validate if the designated constraints
// aren't met.
type OrganizePersonnelListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizePersonnelListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizePersonnelListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizePersonnelListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizePersonnelListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizePersonnelListRequestValidationError) ErrorName() string {
	return "OrganizePersonnelListRequestValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizePersonnelListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizePersonnelListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizePersonnelListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizePersonnelListRequestValidationError{}

// Validate checks the field values on OrganizePersonnelListResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizePersonnelListResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizePersonnelListResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// OrganizePersonnelListResponseMultiError, or nil if none found.
func (m *OrganizePersonnelListResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizePersonnelListResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrganizePersonnelListResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrganizePersonnelListResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrganizePersonnelListResponseValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OrganizePersonnelListResponseMultiError(errors)
	}

	return nil
}

// OrganizePersonnelListResponseMultiError is an error wrapping multiple
// validation errors returned by OrganizePersonnelListResponse.ValidateAll()
// if the designated constraints aren't met.
type OrganizePersonnelListResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizePersonnelListResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizePersonnelListResponseMultiError) AllErrors() []error { return m }

// OrganizePersonnelListResponseValidationError is the validation error
// returned by OrganizePersonnelListResponse.Validate if the designated
// constraints aren't met.
type OrganizePersonnelListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizePersonnelListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizePersonnelListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizePersonnelListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizePersonnelListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizePersonnelListResponseValidationError) ErrorName() string {
	return "OrganizePersonnelListResponseValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizePersonnelListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizePersonnelListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizePersonnelListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizePersonnelListResponseValidationError{}

// Validate checks the field values on OrganizeDepartmentListResponse_Item with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *OrganizeDepartmentListResponse_Item) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizeDepartmentListResponse_Item
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// OrganizeDepartmentListResponse_ItemMultiError, or nil if none found.
func (m *OrganizeDepartmentListResponse_Item) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizeDepartmentListResponse_Item) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for DeptId

	// no validation rules for ParentId

	// no validation rules for DeptName

	// no validation rules for Ancestors

	// no validation rules for Count

	if len(errors) > 0 {
		return OrganizeDepartmentListResponse_ItemMultiError(errors)
	}

	return nil
}

// OrganizeDepartmentListResponse_ItemMultiError is an error wrapping multiple
// validation errors returned by
// OrganizeDepartmentListResponse_Item.ValidateAll() if the designated
// constraints aren't met.
type OrganizeDepartmentListResponse_ItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizeDepartmentListResponse_ItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizeDepartmentListResponse_ItemMultiError) AllErrors() []error { return m }

// OrganizeDepartmentListResponse_ItemValidationError is the validation error
// returned by OrganizeDepartmentListResponse_Item.Validate if the designated
// constraints aren't met.
type OrganizeDepartmentListResponse_ItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizeDepartmentListResponse_ItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizeDepartmentListResponse_ItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizeDepartmentListResponse_ItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizeDepartmentListResponse_ItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizeDepartmentListResponse_ItemValidationError) ErrorName() string {
	return "OrganizeDepartmentListResponse_ItemValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizeDepartmentListResponse_ItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizeDepartmentListResponse_Item.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizeDepartmentListResponse_ItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizeDepartmentListResponse_ItemValidationError{}

// Validate checks the field values on OrganizePersonnelListResponse_Position
// with the rules defined in the proto definition for this message. If any
// rules are violated, the first error encountered is returned, or nil if
// there are no violations.
func (m *OrganizePersonnelListResponse_Position) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on
// OrganizePersonnelListResponse_Position with the rules defined in the proto
// definition for this message. If any rules are violated, the result is a
// list of violation errors wrapped in
// OrganizePersonnelListResponse_PositionMultiError, or nil if none found.
func (m *OrganizePersonnelListResponse_Position) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizePersonnelListResponse_Position) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for Name

	// no validation rules for Sort

	if len(errors) > 0 {
		return OrganizePersonnelListResponse_PositionMultiError(errors)
	}

	return nil
}

// OrganizePersonnelListResponse_PositionMultiError is an error wrapping
// multiple validation errors returned by
// OrganizePersonnelListResponse_Position.ValidateAll() if the designated
// constraints aren't met.
type OrganizePersonnelListResponse_PositionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizePersonnelListResponse_PositionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizePersonnelListResponse_PositionMultiError) AllErrors() []error { return m }

// OrganizePersonnelListResponse_PositionValidationError is the validation
// error returned by OrganizePersonnelListResponse_Position.Validate if the
// designated constraints aren't met.
type OrganizePersonnelListResponse_PositionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizePersonnelListResponse_PositionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizePersonnelListResponse_PositionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizePersonnelListResponse_PositionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizePersonnelListResponse_PositionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizePersonnelListResponse_PositionValidationError) ErrorName() string {
	return "OrganizePersonnelListResponse_PositionValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizePersonnelListResponse_PositionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizePersonnelListResponse_Position.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizePersonnelListResponse_PositionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizePersonnelListResponse_PositionValidationError{}

// Validate checks the field values on OrganizePersonnelListResponse_Dept with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *OrganizePersonnelListResponse_Dept) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizePersonnelListResponse_Dept
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// OrganizePersonnelListResponse_DeptMultiError, or nil if none found.
func (m *OrganizePersonnelListResponse_Dept) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizePersonnelListResponse_Dept) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for DeptId

	// no validation rules for DeptName

	// no validation rules for Ancestors

	if len(errors) > 0 {
		return OrganizePersonnelListResponse_DeptMultiError(errors)
	}

	return nil
}

// OrganizePersonnelListResponse_DeptMultiError is an error wrapping multiple
// validation errors returned by
// OrganizePersonnelListResponse_Dept.ValidateAll() if the designated
// constraints aren't met.
type OrganizePersonnelListResponse_DeptMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizePersonnelListResponse_DeptMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizePersonnelListResponse_DeptMultiError) AllErrors() []error { return m }

// OrganizePersonnelListResponse_DeptValidationError is the validation error
// returned by OrganizePersonnelListResponse_Dept.Validate if the designated
// constraints aren't met.
type OrganizePersonnelListResponse_DeptValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizePersonnelListResponse_DeptValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizePersonnelListResponse_DeptValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizePersonnelListResponse_DeptValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizePersonnelListResponse_DeptValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizePersonnelListResponse_DeptValidationError) ErrorName() string {
	return "OrganizePersonnelListResponse_DeptValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizePersonnelListResponse_DeptValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizePersonnelListResponse_Dept.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizePersonnelListResponse_DeptValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizePersonnelListResponse_DeptValidationError{}

// Validate checks the field values on OrganizePersonnelListResponse_Item with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *OrganizePersonnelListResponse_Item) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizePersonnelListResponse_Item
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// OrganizePersonnelListResponse_ItemMultiError, or nil if none found.
func (m *OrganizePersonnelListResponse_Item) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizePersonnelListResponse_Item) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for Nickname

	// no validation rules for Gender

	for idx, item := range m.GetPositionItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrganizePersonnelListResponse_ItemValidationError{
						field:  fmt.Sprintf("PositionItems[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrganizePersonnelListResponse_ItemValidationError{
						field:  fmt.Sprintf("PositionItems[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrganizePersonnelListResponse_ItemValidationError{
					field:  fmt.Sprintf("PositionItems[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetDeptItem()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, OrganizePersonnelListResponse_ItemValidationError{
					field:  "DeptItem",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, OrganizePersonnelListResponse_ItemValidationError{
					field:  "DeptItem",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDeptItem()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OrganizePersonnelListResponse_ItemValidationError{
				field:  "DeptItem",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Avatar

	if len(errors) > 0 {
		return OrganizePersonnelListResponse_ItemMultiError(errors)
	}

	return nil
}

// OrganizePersonnelListResponse_ItemMultiError is an error wrapping multiple
// validation errors returned by
// OrganizePersonnelListResponse_Item.ValidateAll() if the designated
// constraints aren't met.
type OrganizePersonnelListResponse_ItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizePersonnelListResponse_ItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizePersonnelListResponse_ItemMultiError) AllErrors() []error { return m }

// OrganizePersonnelListResponse_ItemValidationError is the validation error
// returned by OrganizePersonnelListResponse_Item.Validate if the designated
// constraints aren't met.
type OrganizePersonnelListResponse_ItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizePersonnelListResponse_ItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizePersonnelListResponse_ItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizePersonnelListResponse_ItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizePersonnelListResponse_ItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizePersonnelListResponse_ItemValidationError) ErrorName() string {
	return "OrganizePersonnelListResponse_ItemValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizePersonnelListResponse_ItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizePersonnelListResponse_Item.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizePersonnelListResponse_ItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizePersonnelListResponse_ItemValidationError{}
