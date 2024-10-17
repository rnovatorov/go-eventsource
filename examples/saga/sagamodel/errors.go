package sagamodel

import (
	"errors"
)

var (
	ErrSagaAlreadyBegun                                    = errors.New("saga already begun")
	ErrSagaEnded                                           = errors.New("saga ended")
	ErrTaskDefinitionsMissing                              = errors.New("task definitions missing")
	ErrDependencyGraphCyclic                               = errors.New("dependency graph cyclic")
	ErrDependencyGraphBuildingFailed                       = errors.New("dependency graph building failed")
	ErrTaskArgumentAlreadyDefined                          = errors.New("task argument already defined")
	ErrTaskCompensationArgumentAlreadyDefined              = errors.New("task compensation argument already defined")
	ErrTaskArgumentReferencesIndependentResult             = errors.New("task argument references independent result")
	ErrTaskCompensationArgumentReferencesIndependentResult = errors.New("task compensation argument references independent result")

	ErrTaskAbortReasonMissing = errors.New("task abort reason missing")
	ErrTaskAborted            = errors.New("task aborted")
	ErrTaskAlreadyAborted     = errors.New("task already aborted")
	ErrTaskEnded              = errors.New("task ended")
	ErrTaskAlreadyEnded       = errors.New("task already ended")
	ErrTaskNotDefined         = errors.New("task not defined")
	ErrTaskNotBegun           = errors.New("task not begun")
	ErrTaskResultMissing      = errors.New("task result missing")

	ErrCompensationAlreadyEnded = errors.New("compensation already ended")
	ErrCompensationNotBegun     = errors.New("compensation not begun")

	ErrVertexAlreadyAdded = errors.New("vertex already added")
	ErrVertexNotFound     = errors.New("vertex not found")
)
