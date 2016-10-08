package csv

import "fmt"

type EmptyReportError struct {
	action string
}

func NewEmptyReportError(action string) EmptyReportError {
	return EmptyReportError{
		action: action,
	}
}

func (e EmptyReportError) Error() string {
	errString := "Report is empty"
	if e.action != "" {
		errString = fmt.Sprintf("%s: %s returned no valid reports", errString, e.action)
	}
	return errString
}

type ReadCleanError struct {
	pack string
	err  error
}

func NewReadCleanError(pack string, err error) ReadCleanError {
	return ReadCleanError{
		pack: pack,
		err:  err,
	}
}

func (e ReadCleanError) Error() string {
	return fmt.Sprintf("Failed to \"Read\" or \"Clean\" %s reports: %s", e.pack, e.err.Error())
}

type ReportParseError struct {
	pack string
	err  error
}

func NewReportParseError(pack string, err error) ReportParseError {
	return ReportParseError{
		pack: pack,
		err:  err,
	}
}

func (e ReportParseError) Error() string {
	return fmt.Sprintf("Failed to parse reports for %s: %s", e.pack, e.err.Error())
}
