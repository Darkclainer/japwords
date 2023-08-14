package anki

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// ValidationError is a wrapper for error to make possible to distinguish this error in API layer.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

type MappingValidationError struct {
	Key string
	Msg string
}

func (e *MappingValidationError) Error() string {
	return fmt.Sprintf("mapping validation of field %q failed: %s", e.Key, e.Msg)
}

type MappingValidationErrors struct {
	KeyErrors   []*MappingValidationError
	ValueErrors []*MappingValidationError
}

func (*MappingValidationErrors) Error() string {
	return "mapping validation failed"
}

var hostnameRFC1123Regex = regexp.MustCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`)

// validateAddr accepts strings that consists of hostname compaining RFC1123 and port separated by colon.
func validateAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	portVal, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if portVal < 1 || portVal > 65535 {
		return fmt.Errorf("port value should be above 0 and less than 65536")
	}
	if !hostnameRFC1123Regex.MatchString(host) {
		return fmt.Errorf("hostname is not valid")
	}
	return nil
}

// validateMappingKeys checks that mapping keys is valid anki note type field names,
// see validateFieldName. It returns only key errors.
func validateMappingKeys(mapping map[string]string) []*MappingValidationError {
	var errs []*MappingValidationError
	for key := range mapping {
		err := validateFieldName(key)
		if err != nil {
			errs = append(errs, &MappingValidationError{
				Key: key,
				Msg: err.Error(),
			})
		}
	}
	return errs
}

// validateDeckName checks that name is non-empty string that doesn't contain symbol
// `"` and doesn't start or end with spaces.
func validateDeckName(name string) error {
	return validateDeckNoteType(name)
}

// validateNoteType checks that name is non-empty string that doesn't contain symbol
// `"` and doesn't start or end with spaces.
func validateNoteType(name string) error {
	return validateDeckNoteType(name)
}

var (
	deckNoteTypeRegex      = regexp.MustCompile(`(^[^ \t\n\v"][^"\n]*[^ \t\n\v"]$)|(^[^ \t\n\v"]$)`)
	errDeckNoteTypeInvalid = errors.New("must not be empty string, start or end with spaces or contain '\"'")
)

func validateDeckNoteType(name string) error {
	if !deckNoteTypeRegex.MatchString(name) {
		return errDeckNoteTypeInvalid
	}
	return nil
}

var (
	fieldNameRegex      = regexp.MustCompile(`(^[^ \t\n\v{}:"][^{}:"\n]*[^ \t\n\v{}:"]$)|(^[^ \t\n\v{}:"]$)`)
	errFieldNameInvalid = errors.New("must not be empty string, start or end with spaces or contain ':', '\"', '{', '}'")
)

// validateFieldName checks that name is non-empty string that doesn't contain symbols
// `:`, `"`, `{` or `}` and doesn't start or end with spaces.
func validateFieldName(name string) error {
	if !fieldNameRegex.MatchString(name) {
		return errFieldNameInvalid
	}
	return nil
}
