// Code generated by go-swagger; DO NOT EDIT.

package schema

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/weaviate/weaviate/entities/models"
)

// SchemaDumpReader is a Reader for the SchemaDump structure.
type SchemaDumpReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SchemaDumpReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSchemaDumpOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewSchemaDumpUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewSchemaDumpForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewSchemaDumpInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewSchemaDumpOK creates a SchemaDumpOK with default headers values
func NewSchemaDumpOK() *SchemaDumpOK {
	return &SchemaDumpOK{}
}

/*
SchemaDumpOK describes a response with status code 200, with default header values.

Successfully dumped the database schema.
*/
type SchemaDumpOK struct {
	Payload *models.Schema
}

// IsSuccess returns true when this schema dump o k response has a 2xx status code
func (o *SchemaDumpOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this schema dump o k response has a 3xx status code
func (o *SchemaDumpOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this schema dump o k response has a 4xx status code
func (o *SchemaDumpOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this schema dump o k response has a 5xx status code
func (o *SchemaDumpOK) IsServerError() bool {
	return false
}

// IsCode returns true when this schema dump o k response a status code equal to that given
func (o *SchemaDumpOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the schema dump o k response
func (o *SchemaDumpOK) Code() int {
	return 200
}

func (o *SchemaDumpOK) Error() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpOK  %+v", 200, o.Payload)
}

func (o *SchemaDumpOK) String() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpOK  %+v", 200, o.Payload)
}

func (o *SchemaDumpOK) GetPayload() *models.Schema {
	return o.Payload
}

func (o *SchemaDumpOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Schema)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSchemaDumpUnauthorized creates a SchemaDumpUnauthorized with default headers values
func NewSchemaDumpUnauthorized() *SchemaDumpUnauthorized {
	return &SchemaDumpUnauthorized{}
}

/*
SchemaDumpUnauthorized describes a response with status code 401, with default header values.

Unauthorized or invalid credentials.
*/
type SchemaDumpUnauthorized struct {
}

// IsSuccess returns true when this schema dump unauthorized response has a 2xx status code
func (o *SchemaDumpUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this schema dump unauthorized response has a 3xx status code
func (o *SchemaDumpUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this schema dump unauthorized response has a 4xx status code
func (o *SchemaDumpUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this schema dump unauthorized response has a 5xx status code
func (o *SchemaDumpUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this schema dump unauthorized response a status code equal to that given
func (o *SchemaDumpUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the schema dump unauthorized response
func (o *SchemaDumpUnauthorized) Code() int {
	return 401
}

func (o *SchemaDumpUnauthorized) Error() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpUnauthorized ", 401)
}

func (o *SchemaDumpUnauthorized) String() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpUnauthorized ", 401)
}

func (o *SchemaDumpUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSchemaDumpForbidden creates a SchemaDumpForbidden with default headers values
func NewSchemaDumpForbidden() *SchemaDumpForbidden {
	return &SchemaDumpForbidden{}
}

/*
SchemaDumpForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type SchemaDumpForbidden struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this schema dump forbidden response has a 2xx status code
func (o *SchemaDumpForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this schema dump forbidden response has a 3xx status code
func (o *SchemaDumpForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this schema dump forbidden response has a 4xx status code
func (o *SchemaDumpForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this schema dump forbidden response has a 5xx status code
func (o *SchemaDumpForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this schema dump forbidden response a status code equal to that given
func (o *SchemaDumpForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the schema dump forbidden response
func (o *SchemaDumpForbidden) Code() int {
	return 403
}

func (o *SchemaDumpForbidden) Error() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpForbidden  %+v", 403, o.Payload)
}

func (o *SchemaDumpForbidden) String() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpForbidden  %+v", 403, o.Payload)
}

func (o *SchemaDumpForbidden) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *SchemaDumpForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSchemaDumpInternalServerError creates a SchemaDumpInternalServerError with default headers values
func NewSchemaDumpInternalServerError() *SchemaDumpInternalServerError {
	return &SchemaDumpInternalServerError{}
}

/*
SchemaDumpInternalServerError describes a response with status code 500, with default header values.

An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.
*/
type SchemaDumpInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this schema dump internal server error response has a 2xx status code
func (o *SchemaDumpInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this schema dump internal server error response has a 3xx status code
func (o *SchemaDumpInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this schema dump internal server error response has a 4xx status code
func (o *SchemaDumpInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this schema dump internal server error response has a 5xx status code
func (o *SchemaDumpInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this schema dump internal server error response a status code equal to that given
func (o *SchemaDumpInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the schema dump internal server error response
func (o *SchemaDumpInternalServerError) Code() int {
	return 500
}

func (o *SchemaDumpInternalServerError) Error() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpInternalServerError  %+v", 500, o.Payload)
}

func (o *SchemaDumpInternalServerError) String() string {
	return fmt.Sprintf("[GET /schema][%d] schemaDumpInternalServerError  %+v", 500, o.Payload)
}

func (o *SchemaDumpInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *SchemaDumpInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
