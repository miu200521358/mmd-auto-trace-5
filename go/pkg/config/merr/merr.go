package merr

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

// BaseError はすべてのカスタムエラーの基本構造体
type BaseError struct {
	msg        string
	stackTrace string
}

// Error はerror interfaceを実装
func (e *BaseError) Error() string {
	return e.msg
}

// StackTrace はエラー発生時のスタックトレースを返す
func (e *BaseError) StackTrace() string {
	return e.stackTrace
}

// captureStackTrace はエラー発生時のスタックトレースを取得する
func captureStackTrace() string {
	buf := make([]byte, 1<<20) // 1MB のバッファを確保
	n := runtime.Stack(buf, true)
	return string(bytes.ReplaceAll(buf[:n], []byte("\n"), []byte("\r\n")))
}

// NameNotFoundError は名前が見つからないエラー
type NameNotFoundError struct {
	*BaseError
	Name string
}

// NewNameNotFoundError は新しいNameNotFoundErrorを生成する
func NewNameNotFoundError(name string, msg string) *NameNotFoundError {
	return &NameNotFoundError{
		BaseError: &BaseError{
			msg: msg,
		},
		Name: name,
	}
}

// IsNameNotFoundError はエラーがNameNotFoundErrorかどうか判定する
func IsNameNotFoundError(err error) bool {
	var nameErr *NameNotFoundError
	return errors.As(err, &nameErr)
}

// ParentNotFoundError は親が見つからないエラー
type ParentNotFoundError struct {
	*BaseError
	Parent string
}

// NewParentNotFoundError は新しいParentNotFoundErrorを生成する
func NewParentNotFoundError(parent, message string) *ParentNotFoundError {
	return &ParentNotFoundError{
		BaseError: &BaseError{
			msg: message,
		},
		Parent: parent,
	}
}

// IsParentNotFoundError はエラーがParentNotFoundErrorかどうか判定する
func IsParentNotFoundError(err error) bool {
	var parentErr *ParentNotFoundError
	return errors.As(err, &parentErr)
}

// TerminateError は終了エラー
type TerminateError struct {
	*BaseError
	Reason string
}

// NewTerminateError は新しいTerminateErrorを生成する
func NewTerminateError(reason string) *TerminateError {
	return &TerminateError{
		BaseError: &BaseError{
			msg:        fmt.Sprintf("terminate error: %s", reason),
			stackTrace: captureStackTrace(),
		},
		Reason: reason,
	}
}

// IsTerminateError はエラーがTerminateErrorかどうか判定する
func IsTerminateError(err error) bool {
	var termErr *TerminateError
	return errors.As(err, &termErr)
}
