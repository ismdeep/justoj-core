package main

// ResultPending pending
const ResultPending = 0

// ResultPendingRejudge rejudge
const ResultPendingRejudge = 1

// ResultCompiling compiling
const ResultCompiling = 2

// ResultRunning running
const ResultRunning = 3

// ResultAccepted accepted
const ResultAccepted = 4

// ResultPresentationError presentation error
const ResultPresentationError = 5

// ResultWrongAnswer wrong answer
const ResultWrongAnswer = 6

// ResultTimeLimit time limit
const ResultTimeLimit = 7

// ResultMemLimit memory limit
const ResultMemLimit = 8

// ResultOutputLimit output limit
const ResultOutputLimit = 9

// ResultRuntimeError runtime error
const ResultRuntimeError = 10

// ResultCompileError compile error
const ResultCompileError = 11

// ResultCompilePassed compile passed
const ResultCompilePassed = 12

// ResultTestRunning test running
const ResultTestRunning = 13

// ResultSubmitOK submit ok
const ResultSubmitOK = 14

// ResultText result text map
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
