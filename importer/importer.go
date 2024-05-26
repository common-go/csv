package importer

import (
	"context"
	"fmt"
	"io"
)

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

type ErrorHandler[T any] struct {
	HandleError func(ctx context.Context, format string, fields map[string]interface{})
	FileName    string
	LineNumber  string
	Map         *map[string]interface{}
}

func NewErrorHandler[T any](logger func(ctx context.Context, format string, fields map[string]interface{}), fileName string, lineNumber string, mp *map[string]interface{}) *ErrorHandler[T] {
	if len(fileName) <= 0 {
		fileName = "filename"
	}
	if len(lineNumber) <= 0 {
		lineNumber = "lineNumber"
	}
	return &ErrorHandler[T]{
		HandleError: logger,
		FileName:    fileName,
		LineNumber:  lineNumber,
		Map:         mp,
	}
}

func (e *ErrorHandler[T]) HandlerError(ctx context.Context, raw string, rs T, err []ErrorMessage, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v", raw, rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v line: %d", raw, rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v filename:%s", raw, rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v filename:%s line: %d", raw, rs, err, fileName, i), ext)
	}
}

func (e *ErrorHandler[T]) HandlerException(ctx context.Context, raw string, rs T, err error, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v", raw, rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v line: %d", raw, rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v filename:%s", raw, rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Error to write:  %s %+v . Error: %v filename: %s line: %d", raw, rs, err, fileName, i), ext)
	}
}

func NewImportRepository[T any](
	transform func(ctx context.Context, lines string) (T, error),
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs *T, err error, i int, fileName string),
	validate func(ctx context.Context, model *T) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs *T, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data *T) error,
	opt ...func(ctx context.Context) error,
) *Importer[T] {
	return NewImporter[T](transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImportAdapter[T any](
	transform func(ctx context.Context, lines string) (T, error),
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs *T, err error, i int, fileName string),
	validate func(ctx context.Context, model *T) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs *T, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data *T) error,
	opt ...func(ctx context.Context) error,
) *Importer[T] {
	return NewImporter[T](transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImportService[T any](
	transform func(ctx context.Context, lines string) (T, error),
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs *T, err error, i int, fileName string),
	validate func(ctx context.Context, model *T) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs *T, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data *T) error,
	opt ...func(ctx context.Context) error,
) *Importer[T] {
	return NewImporter[T](transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImporter[T any](
	transform func(ctx context.Context, lines string) (T, error),
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs *T, err error, i int, fileName string),
	validate func(ctx context.Context, model *T) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs *T, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data *T) error,
	opt ...func(ctx context.Context) error,
) *Importer[T] {
	var flush func(ctx context.Context) error
	if len(opt) > 0 {
		flush = opt[0]
	}
	return &Importer[T]{Transform: transform, Write: write, Flush: flush, Read: read, Validate: validate, HandleError: handleError, HandleException: handleException, Filename: filename}
}

type Importer[T any] struct {
	Transform       func(ctx context.Context, lines string) (T, error)
	Read            func(next func(line string, err error, numLine int) error) error
	Validate        func(ctx context.Context, model *T) ([]ErrorMessage, error)
	HandleError     func(ctx context.Context, raw string, rs *T, err []ErrorMessage, i int, fileName string)
	HandleException func(ctx context.Context, raw string, rs *T, err error, i int, fileName string)
	Filename        string
	Write           func(ctx context.Context, data *T) error
	Flush           func(ctx context.Context) error
}

func (s *Importer[T]) Import(ctx context.Context) (total int, success int, err error) {
	err = s.Read(func(line string, err error, numLine int) error {
		if err == io.EOF {
			if s.Flush != nil {
				return s.Flush(ctx)
			}
			return nil
		}
		total++
		record, err := s.Transform(ctx, line)
		if err != nil {
			return err
		}
		if s.Validate != nil {
			errs, err := s.Validate(ctx, &record)
			if err != nil {
				return err
			}
			if len(errs) > 0 {
				s.HandleError(ctx, line, &record, errs, numLine, s.Filename)
				return nil
			}
		}
		err = s.Write(ctx, &record)
		if err != nil {
			if s.HandleException != nil {
				s.HandleException(ctx, line, &record, err, numLine, s.Filename)
				return nil
			} else {
				return err
			}
		}
		success++
		return nil
	})
	if err != nil && err != io.EOF {
		return total, success, err
	}
	return total, success, nil
}
