// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cli

import (
	"os"
	"time"

	"github.com/fatih/color"

	"github.com/uber/cadence/common/types"
)

const (
	localHost    = "127.0.0.1"
	tchannelPort = localHost + ":7933"
	grpcPort     = localHost + ":7833"

	grpcTransport   = "grpc"
	thriftTransport = "tchannel"

	maxOutputStringLength = 200 // max length for output string
	maxWorkflowTypeLength = 32  // max item length for output workflow type in table
	defaultMaxFieldLength = 500 // default max length for each attribute field
	maxWordLength         = 120 // if text length is larger than maxWordLength, it will be inserted spaces

	// regex expression for parsing time durations, shorter, longer notations and numeric value respectively
	defaultDateTimeRangeShortRE = "^[1-9][0-9]*[smhdwMy]$"                                // eg. 1s, 20m, 300h etc.
	defaultDateTimeRangeLongRE  = "^[1-9][0-9]*(second|minute|hour|day|week|month|year)$" // eg. 1second, 20minute, 300hour etc.
	defaultDateTimeRangeNum     = "^[1-9][0-9]*"                                          // eg. 1, 20, 300 etc.

	// time ranges
	day   = 24 * time.Hour
	week  = 7 * day
	month = 30 * day
	year  = 365 * day

	defaultTimeFormat                            = "15:04:05"   // used for converting UnixNano to string like 16:16:36 (only time)
	defaultDateTimeFormat                        = time.RFC3339 // used for converting UnixNano to string like 2018-02-15T16:16:36-08:00
	defaultDomainRetentionDays                   = 3
	defaultContextTimeoutInSeconds               = 5
	defaultContextTimeout                        = defaultContextTimeoutInSeconds * time.Second
	defaultContextTimeoutForLongPoll             = 2 * time.Minute
	defaultContextTimeoutForListArchivedWorkflow = 3 * time.Minute

	defaultDecisionTimeoutInSeconds = 10
	defaultPageSizeForList          = 500
	defaultPageSizeForScan          = 2000
	defaultWorkflowIDReusePolicy    = types.WorkflowIDReusePolicyAllowDuplicateFailedOnly

	workflowStatusNotSet = -1
	showErrorStackEnv    = `CADENCE_CLI_SHOW_STACKS`

	searchAttrInputSeparator = "|"

	defaultGracefulFailoverTimeoutInSeconds = 60
)

var envKeysForUserName = []string{
	"USER",
	"LOGNAME",
	"HOME",
}

const resetTypeFirstDecisionCompleted = "FirstDecisionCompleted"
const resetTypeLastDecisionCompleted = "LastDecisionCompleted"
const resetTypeLastContinuedAsNew = "LastContinuedAsNew"
const resetTypeBadBinary = "BadBinary"
const resetTypeDecisionCompletedTime = "DecisionCompletedTime"
const resetTypeFirstDecisionScheduled = "FirstDecisionScheduled"
const resetTypeLastDecisionScheduled = "LastDecisionScheduled"

var resetTypesMap = map[string]string{
	resetTypeFirstDecisionCompleted: "",
	resetTypeLastDecisionCompleted:  "",
	resetTypeLastContinuedAsNew:     "",
	resetTypeBadBinary:              FlagResetBadBinaryChecksum,
	resetTypeDecisionCompletedTime:  FlagEarliestTime,
	resetTypeFirstDecisionScheduled: "",
	resetTypeLastDecisionScheduled:  "",
}

type jsonType int

const (
	jsonTypeInput jsonType = iota
	jsonTypeMemo
	jsonTypeHeader
	jsonTypeSignal
)

var (
	colorRed     = color.New(color.FgRed).SprintFunc()
	colorMagenta = color.New(color.FgMagenta).SprintFunc()
	colorGreen   = color.New(color.FgGreen).SprintFunc()

	optionErr               = "there is something wrong with your command options"
	osExit                  = os.Exit
	workflowClosedStatusMap = map[string]types.WorkflowExecutionCloseStatus{
		"completed":        types.WorkflowExecutionCloseStatusCompleted,
		"failed":           types.WorkflowExecutionCloseStatusFailed,
		"canceled":         types.WorkflowExecutionCloseStatusCanceled,
		"terminated":       types.WorkflowExecutionCloseStatusTerminated,
		"continued_as_new": types.WorkflowExecutionCloseStatusContinuedAsNew,
		"timed_out":        types.WorkflowExecutionCloseStatusTimedOut,
		// below are some alias
		"c":              types.WorkflowExecutionCloseStatusCompleted,
		"complete":       types.WorkflowExecutionCloseStatusCompleted,
		"f":              types.WorkflowExecutionCloseStatusFailed,
		"fail":           types.WorkflowExecutionCloseStatusFailed,
		"cancel":         types.WorkflowExecutionCloseStatusCanceled,
		"terminate":      types.WorkflowExecutionCloseStatusTerminated,
		"term":           types.WorkflowExecutionCloseStatusTerminated,
		"continue":       types.WorkflowExecutionCloseStatusContinuedAsNew,
		"cont":           types.WorkflowExecutionCloseStatusContinuedAsNew,
		"continuedasnew": types.WorkflowExecutionCloseStatusContinuedAsNew,
		"continueasnew":  types.WorkflowExecutionCloseStatusContinuedAsNew,
		"timedout":       types.WorkflowExecutionCloseStatusTimedOut,
		"timeout":        types.WorkflowExecutionCloseStatusTimedOut,
	}
)
