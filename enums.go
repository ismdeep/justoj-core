package main

const ResultPending = 0
const ResultPendingRejudge = 1
const ResultCompiling = 2
const ResultRunning = 3
const ResultAccepted = 4
const ResultPresentationError = 5
const ResultWrongAnswer = 6
const ResultTimeLimit = 7
const ResultMemLimit = 8
const ResultOutputLimit = 9
const ResultRuntimeError = 10
const ResultCompileError = 11
const ResultCompilePassed = 12
const ResultTestRunning = 13
const ResultSubmitOK = 14

var ResultText map[int]string

func init() {
	ResultText = make(map[int]string)
	ResultText[ResultPending] = "Pending"
	ResultText[ResultPendingRejudge] = "Rejudging"
	ResultText[ResultCompiling] = "Compiling"
	ResultText[ResultRunning] = "Running"
	ResultText[ResultAccepted] = "Accepted"
	ResultText[ResultPresentationError] = "PresentationError"
	ResultText[ResultWrongAnswer] = "WrongAnswer"
	ResultText[ResultTimeLimit] = "TimeLimit"
	ResultText[ResultMemLimit] = "MemLimit"
	ResultText[ResultOutputLimit] = "OutputLimit"
	ResultText[ResultRuntimeError] = "RuntimeError"
	ResultText[ResultCompileError] = "CompileError"
	ResultText[ResultCompilePassed] = "CompilePassed"
	ResultText[ResultTestRunning] = "TestRunning"
	ResultText[ResultSubmitOK] = "SubmitOK"
}
