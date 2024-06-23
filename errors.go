package duckdb

import "C"
import (
	"errors"
	"fmt"
	"strings"
)

func getError(errDriver error, err error) error {
	if err == nil {
		return fmt.Errorf("%s: %w", driverErrMsg, errDriver)
	}
	return fmt.Errorf("%s: %w: %s", driverErrMsg, errDriver, err.Error())
}

func duckdbError(err *C.char) error {
	return fmt.Errorf("%s: %w", duckdbErrMsg, errors.New(C.GoString(err)))
}

func castError(actual string, expected string) error {
	return fmt.Errorf("%s: cannot cast %s to %s", castErrMsg, actual, expected)
}

func structFieldError(actual string, expected string) error {
	return fmt.Errorf("%s: expected %s, got %s", structFieldErrMsg, expected, actual)
}

func columnError(err error, colIdx int) error {
	return fmt.Errorf("%w: %s: %d", err, columnErrMsg, colIdx)
}

func columnCountError(actual int, expected int) error {
	return fmt.Errorf("%s: expected %d, got %d", columnCountErrMsg, expected, actual)
}

func unsupportedTypeError(name string) error {
	return fmt.Errorf("%s: %s", unsupportedTypeErrMsg, name)
}

func invalidatedAppenderError(err error) error {
	if err == nil {
		return fmt.Errorf(invalidatedAppenderMsg)
	}
	return fmt.Errorf("%w: %s", err, invalidatedAppenderMsg)
}

const (
	driverErrMsg           = "database/sql/driver"
	duckdbErrMsg           = "duckdb error"
	castErrMsg             = "cast error"
	structFieldErrMsg      = "invalid STRUCT field"
	columnErrMsg           = "column index"
	columnCountErrMsg      = "invalid column count"
	unsupportedTypeErrMsg  = "unsupported data type"
	invalidatedAppenderMsg = "appended data has been invalidated due to corrupt row"
)

var (
	errAPI        = errors.New("API error")
	errVectorSize = errors.New("data chunks cannot exceed duckdb's internal vector size")

	errParseDSN  = errors.New("could not parse DSN for database")
	errOpen      = errors.New("could not open database")
	errSetConfig = errors.New("could not set invalid or local option for global database config")

	errUnsupportedMapKeyType = errors.New("MAP key type not supported")

	errAppenderInvalidCon       = errors.New("could not create appender: not a DuckDB driver connection")
	errAppenderClosedCon        = errors.New("could not create appender: appender creation on a closed connection")
	errAppenderCreation         = errors.New("could not create appender")
	errAppenderDoubleClose      = errors.New("could not close appender: already closed")
	errAppenderAppendRow        = errors.New("could not append row")
	errAppenderAppendAfterClose = errors.New("could not append row: appender already closed")
	errAppenderClose            = errors.New("could not close appender")
	errAppenderFlush            = errors.New("could not flush appender")

	// Errors not covered in tests.
	errConnect      = errors.New("could not connect to database")
	errCreateConfig = errors.New("could not create config for database")
)

type DuckDBErrorType int

const (
	ErrorTypeInvalid              DuckDBErrorType = iota // invalid type
	ErrorTypeOutOfRange                                  // value out of range error
	ErrorTypeConversion                                  // conversion/casting error
	ErrorTypeUnknownType                                 // unknown type
	ErrorTypeDecimal                                     // decimal related
	ErrorTypeMismatchType                                // type mismatch
	ErrorTypeDivideByZero                                // divide by 0
	ErrorTypeObjectSize                                  // object size exceeded
	ErrorTypeInvalidType                                 // incompatible for operation
	ErrorTypeSerialization                               // serialization
	ErrorTypeTransaction                                 // transaction management
	ErrorTypeNotImplemented                              // method not implemented
	ErrorTypeExpression                                  // expression parsing
	ErrorTypeCatalog                                     // catalog related
	ErrorTypeParser                                      // parser related
	ErrorTypePlanner                                     // planner related
	ErrorTypeScheduler                                   // scheduler related
	ErrorTypeExecutor                                    // executor related
	ErrorTypeConstraint                                  // constraint related
	ErrorTypeIndex                                       // index related
	ErrorTypeStat                                        // stat related
	ErrorTypeConnection                                  // connection related
	ErrorTypeSyntax                                      // syntax related
	ErrorTypeSettings                                    // settings related
	ErrorTypeBinder                                      // binder related
	ErrorTypeNetwork                                     // network related
	ErrorTypeOptimizer                                   // optimizer related
	ErrorTypeNullPointer                                 // nullptr exception
	ErrorTypeIO                                          // IO exception
	ErrorTypeInterrupt                                   // interrupt
	ErrorTypeFatal                                       // Fatal exceptions are non-recoverable, and render the entire DB in an unusable state
	ErrorTypeInternal                                    // Internal exceptions indicate something went wrong internally (i.e. bug in the code base)
	ErrorTypeInvalidInput                                // Input or arguments error
	ErrorTypeOutOfMemory                                 // out of memory
	ErrorTypePermission                                  // insufficient permissions
	ErrorTypeParameterNotResolved                        // parameter types could not be resolved
	ErrorTypeParameterNotAllowed                         // parameter types not allowed
	ErrorTypeDependency                                  // dependency
	ErrorTypeHTTP
	ErrorTypeMissingExtension // Thrown when an extension is used but not loaded
	ErrorTypeAutoLoad         // Thrown when an extension is used but not loaded
	ErrorTypeSequence
	DuckDBExceptionUnknown DuckDBErrorType = -1
)

var exceptionPrefixMap = map[DuckDBErrorType]string{
	ErrorTypeInvalid:              "Invalid",
	ErrorTypeOutOfRange:           "Out of Range",
	ErrorTypeConversion:           "Conversion",
	ErrorTypeUnknownType:          "Unknown Type",
	ErrorTypeDecimal:              "Decimal",
	ErrorTypeMismatchType:         "Mismatch Type",
	ErrorTypeDivideByZero:         "Divide by Zero",
	ErrorTypeObjectSize:           "Object Size",
	ErrorTypeInvalidType:          "Invalid type",
	ErrorTypeSerialization:        "Serialization",
	ErrorTypeTransaction:          "TransactionContext",
	ErrorTypeNotImplemented:       "Not implemented",
	ErrorTypeExpression:           "Expression",
	ErrorTypeCatalog:              "Catalog",
	ErrorTypeParser:               "Parser",
	ErrorTypePlanner:              "Planner",
	ErrorTypeScheduler:            "Scheduler",
	ErrorTypeExecutor:             "Executor",
	ErrorTypeConstraint:           "Constraint",
	ErrorTypeIndex:                "Index",
	ErrorTypeStat:                 "Stat",
	ErrorTypeConnection:           "Connection",
	ErrorTypeSyntax:               "Syntax",
	ErrorTypeSettings:             "Settings",
	ErrorTypeBinder:               "Binder",
	ErrorTypeNetwork:              "Network",
	ErrorTypeOptimizer:            "Optimizer",
	ErrorTypeNullPointer:          "NullPointer",
	ErrorTypeIO:                   "IO",
	ErrorTypeInterrupt:            "INTERRUPT",
	ErrorTypeFatal:                "FATAL",
	ErrorTypeInternal:             "INTERNAL",
	ErrorTypeInvalidInput:         "Invalid Input",
	ErrorTypeOutOfMemory:          "Out of Memory",
	ErrorTypePermission:           "Permission",
	ErrorTypeParameterNotResolved: "Parameter Not Resolved",
	ErrorTypeParameterNotAllowed:  "Parameter Not Allowed",
	ErrorTypeDependency:           "Dependency",
	ErrorTypeHTTP:                 "HTTP",
	ErrorTypeMissingExtension:     "Missing Extension",
	ErrorTypeAutoLoad:             "Extension Autoloading",
	ErrorTypeSequence:             "Sequence",
}

type DuckDBError struct {
	Type DuckDBErrorType
	Msg  string
}

func (de *DuckDBError) Error() string {
	return de.Msg
}

func (de *DuckDBError) Is(err error) bool {
	if derr, ok := err.(*DuckDBError); ok {
		return derr.Msg == de.Msg
	}
	return false
}

func getDuckDBError(errMsg string) error {
	errType := DuckDBExceptionUnknown
	for k, v := range exceptionPrefixMap {
		if strings.HasPrefix(errMsg, v+" Error") {
			errType = k
			break
		}
	}
	return &DuckDBError{
		Type: errType,
		Msg:  errMsg,
	}
}
