//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BatchDeleteResponse Delete Objects response.
//
// swagger:model BatchDeleteResponse
type BatchDeleteResponse struct {

	// If true, objects will not be deleted yet, but merely listed. Defaults to false.
	DryRun *bool `json:"dryRun,omitempty"`

	// match
	Match *BatchDeleteResponseMatch `json:"match,omitempty"`

	// Controls the verbosity of the output, possible values are: "minimal", "verbose". Defaults to "minimal".
	Output *string `json:"output,omitempty"`

	// results
	Results *BatchDeleteResponseResults `json:"results,omitempty"`
}

// Validate validates this batch delete response
func (m *BatchDeleteResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMatch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponse) validateMatch(formats strfmt.Registry) error {
	if swag.IsZero(m.Match) { // not required
		return nil
	}

	if m.Match != nil {
		if err := m.Match.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("match")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("match")
			}
			return err
		}
	}

	return nil
}

func (m *BatchDeleteResponse) validateResults(formats strfmt.Registry) error {
	if swag.IsZero(m.Results) { // not required
		return nil
	}

	if m.Results != nil {
		if err := m.Results.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("results")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("results")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this batch delete response based on the context it is used
func (m *BatchDeleteResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMatch(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponse) contextValidateMatch(ctx context.Context, formats strfmt.Registry) error {

	if m.Match != nil {
		if err := m.Match.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("match")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("match")
			}
			return err
		}
	}

	return nil
}

func (m *BatchDeleteResponse) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	if m.Results != nil {
		if err := m.Results.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("results")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("results")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BatchDeleteResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BatchDeleteResponse) UnmarshalBinary(b []byte) error {
	var res BatchDeleteResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BatchDeleteResponseMatch Outlines how to find the objects to be deleted.
//
// swagger:model BatchDeleteResponseMatch
type BatchDeleteResponseMatch struct {

	// Class (name) which objects will be deleted.
	// Example: City
	Class string `json:"class,omitempty"`

	// Filter to limit the objects to be deleted.
	Where *WhereFilter `json:"where,omitempty"`
}

// Validate validates this batch delete response match
func (m *BatchDeleteResponseMatch) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateWhere(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseMatch) validateWhere(formats strfmt.Registry) error {
	if swag.IsZero(m.Where) { // not required
		return nil
	}

	if m.Where != nil {
		if err := m.Where.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("match" + "." + "where")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("match" + "." + "where")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this batch delete response match based on the context it is used
func (m *BatchDeleteResponseMatch) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateWhere(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseMatch) contextValidateWhere(ctx context.Context, formats strfmt.Registry) error {

	if m.Where != nil {
		if err := m.Where.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("match" + "." + "where")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("match" + "." + "where")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BatchDeleteResponseMatch) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BatchDeleteResponseMatch) UnmarshalBinary(b []byte) error {
	var res BatchDeleteResponseMatch
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BatchDeleteResponseResults batch delete response results
//
// swagger:model BatchDeleteResponseResults
type BatchDeleteResponseResults struct {

	// How many objects should have been deleted but could not be deleted.
	Failed int64 `json:"failed"`

	// The most amount of objects that can be deleted in a single query, equals QUERY_MAXIMUM_RESULTS.
	Limit int64 `json:"limit"`

	// How many objects were matched by the filter.
	Matches int64 `json:"matches"`

	// With output set to "minimal" only objects with error occurred will the be described. Successfully deleted objects would be omitted. Output set to "verbose" will list all of the objets with their respective statuses.
	Objects []*BatchDeleteResponseResultsObjectsItems0 `json:"objects"`

	// How many objects were successfully deleted in this round.
	Successful int64 `json:"successful"`
}

// Validate validates this batch delete response results
func (m *BatchDeleteResponseResults) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateObjects(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseResults) validateObjects(formats strfmt.Registry) error {
	if swag.IsZero(m.Objects) { // not required
		return nil
	}

	for i := 0; i < len(m.Objects); i++ {
		if swag.IsZero(m.Objects[i]) { // not required
			continue
		}

		if m.Objects[i] != nil {
			if err := m.Objects[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("results" + "." + "objects" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("results" + "." + "objects" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this batch delete response results based on the context it is used
func (m *BatchDeleteResponseResults) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateObjects(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseResults) contextValidateObjects(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Objects); i++ {

		if m.Objects[i] != nil {
			if err := m.Objects[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("results" + "." + "objects" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("results" + "." + "objects" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BatchDeleteResponseResults) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BatchDeleteResponseResults) UnmarshalBinary(b []byte) error {
	var res BatchDeleteResponseResults
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BatchDeleteResponseResultsObjectsItems0 Results for this specific Object.
//
// swagger:model BatchDeleteResponseResultsObjectsItems0
type BatchDeleteResponseResultsObjectsItems0 struct {

	// errors
	Errors *ErrorResponse `json:"errors,omitempty"`

	// ID of the Object.
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// status
	// Enum: [SUCCESS DRYRUN FAILED]
	Status *string `json:"status,omitempty"`
}

// Validate validates this batch delete response results objects items0
func (m *BatchDeleteResponseResultsObjectsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateErrors(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseResultsObjectsItems0) validateErrors(formats strfmt.Registry) error {
	if swag.IsZero(m.Errors) { // not required
		return nil
	}

	if m.Errors != nil {
		if err := m.Errors.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("errors")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("errors")
			}
			return err
		}
	}

	return nil
}

func (m *BatchDeleteResponseResultsObjectsItems0) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

var batchDeleteResponseResultsObjectsItems0TypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["SUCCESS","DRYRUN","FAILED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		batchDeleteResponseResultsObjectsItems0TypeStatusPropEnum = append(batchDeleteResponseResultsObjectsItems0TypeStatusPropEnum, v)
	}
}

const (

	// BatchDeleteResponseResultsObjectsItems0StatusSUCCESS captures enum value "SUCCESS"
	BatchDeleteResponseResultsObjectsItems0StatusSUCCESS string = "SUCCESS"

	// BatchDeleteResponseResultsObjectsItems0StatusDRYRUN captures enum value "DRYRUN"
	BatchDeleteResponseResultsObjectsItems0StatusDRYRUN string = "DRYRUN"

	// BatchDeleteResponseResultsObjectsItems0StatusFAILED captures enum value "FAILED"
	BatchDeleteResponseResultsObjectsItems0StatusFAILED string = "FAILED"
)

// prop value enum
func (m *BatchDeleteResponseResultsObjectsItems0) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, batchDeleteResponseResultsObjectsItems0TypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *BatchDeleteResponseResultsObjectsItems0) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", *m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this batch delete response results objects items0 based on the context it is used
func (m *BatchDeleteResponseResultsObjectsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateErrors(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchDeleteResponseResultsObjectsItems0) contextValidateErrors(ctx context.Context, formats strfmt.Registry) error {

	if m.Errors != nil {
		if err := m.Errors.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("errors")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("errors")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BatchDeleteResponseResultsObjectsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BatchDeleteResponseResultsObjectsItems0) UnmarshalBinary(b []byte) error {
	var res BatchDeleteResponseResultsObjectsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
