package repository

import (
	"database/sql"

	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

// SearchQuery defines the which and how data are queried
type SearchQuery struct {
	Columns eventstore.Columns

	SubQueries            [][]*Filter
	Tx                    *sql.Tx
	AllowTimeTravel       bool
	AwaitOpenTransactions bool
	Limit                 uint64
	Desc                  bool

	InstanceID        *Filter
	ExcludedInstances *Filter
	Creator           *Filter
	Owner             *Filter
	Position          *Filter
	Sequence          *Filter
	CreatedAfter      *Filter
	CreatedBefore     *Filter
}

// Filter represents all fields needed to compare a field of an event with a value
type Filter struct {
	Field     Field
	Value     interface{}
	Operation Operation
}

// Operation defines how fields are compared
type Operation int32

const (
	// OperationEquals compares two values for equality
	OperationEquals Operation = iota + 1
	// OperationGreater compares if the given values is greater than the stored one
	OperationGreater
	// OperationLess compares if the given values is less than the stored one
	OperationLess
	//OperationIn checks if a stored value matches one of the passed value list
	OperationIn
	//OperationJSONContains checks if a stored value matches the given json
	OperationJSONContains
	//OperationNotIn checks if a stored value does not match one of the passed value list
	OperationNotIn

	operationCount
)

// Field is the representation of a field from the event
type Field int32

const (
	//FieldAggregateType represents the aggregate type field
	FieldAggregateType Field = iota + 1
	//FieldAggregateID represents the aggregate id field
	FieldAggregateID
	//FieldSequence represents the sequence field
	FieldSequence
	//FieldResourceOwner represents the resource owner field
	FieldResourceOwner
	//FieldInstanceID represents the instance id field
	FieldInstanceID
	//FieldEditorService represents the editor service field
	FieldEditorService
	//FieldEditorUser represents the editor user field
	FieldEditorUser
	//FieldEventType represents the event type field
	FieldEventType
	//FieldEventData represents the event data field
	FieldEventData
	//FieldCreationDate represents the creation date field
	FieldCreationDate
	// FieldPosition represents the field of the global sequence
	FieldPosition

	fieldCount
)

// NewFilter is used in tests. Use searchQuery.*Filter() instead
func NewFilter(field Field, value interface{}, operation Operation) *Filter {
	return &Filter{
		Field:     field,
		Value:     value,
		Operation: operation,
	}
}

// Validate checks if the fields of the filter have valid values
func (f *Filter) Validate() error {
	if f == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-z6KcG", "filter is nil")
	}
	if f.Field <= 0 || f.Field >= fieldCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-zw62U", "field not definded")
	}
	if f.Value == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-GJ9ct", "no value definded")
	}
	if f.Operation <= 0 || f.Operation >= operationCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-RrQTy", "operation not definded")
	}
	return nil
}

func QueryFromBuilder(builder *eventstore.SearchQueryBuilder) (*SearchQuery, error) {
	if builder == nil ||
		builder.GetColumns().Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "builder invalid")
	}

	query := &SearchQuery{
		Columns:               builder.GetColumns(),
		Limit:                 builder.GetLimit(),
		Desc:                  builder.GetDesc(),
		Tx:                    builder.GetTx(),
		AllowTimeTravel:       builder.GetAllowTimeTravel(),
		AwaitOpenTransactions: builder.GetAwaitOpenTransactions(),
		// Queries:               make([]*Filter, 0, 7),
		SubQueries: make([][]*Filter, len(builder.GetQueries())),
	}

	for _, f := range []func(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter{
		instanceIDFilter,
		excludedInstanceIDFilter,
		editorUserFilter,
		resourceOwnerFilter,
		positionAfterFilter,
		eventSequenceGreaterFilter,
		creationDateAfterFilter,
		creationDateBeforeFilter,
	} {
		filter := f(builder, query)
		if filter == nil {
			continue
		}
		if err := filter.Validate(); err != nil {
			return nil, err
		}
	}

	for i, q := range builder.GetQueries() {
		for _, f := range []func(query *eventstore.SearchQuery) *Filter{
			aggregateTypeFilter,
			aggregateIDFilter,
			eventTypeFilter,
			eventDataFilter,
		} {
			filter := f(q)
			if filter == nil {
				continue
			}
			if err := filter.Validate(); err != nil {
				return nil, err
			}
			query.SubQueries[i] = append(query.SubQueries[i], filter)
		}
	}

	return query, nil
}

func eventSequenceGreaterFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetEventSequenceGreater() == 0 {
		return nil
	}
	sortOrder := OperationGreater
	if builder.GetDesc() {
		sortOrder = OperationLess
	}
	query.Sequence = NewFilter(FieldSequence, builder.GetEventSequenceGreater(), sortOrder)
	return query.Sequence
}

func excludedInstanceIDFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if len(builder.GetExcludedInstanceIDs()) == 0 {
		return nil
	}
	query.ExcludedInstances = NewFilter(FieldInstanceID, database.TextArray[string](builder.GetExcludedInstanceIDs()), OperationNotIn)
	return query.ExcludedInstances
}

func creationDateAfterFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetCreationDateAfter().IsZero() {
		return nil
	}
	query.CreatedAfter = NewFilter(FieldCreationDate, builder.GetCreationDateAfter(), OperationGreater)
	return query.CreatedAfter
}

func creationDateBeforeFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetCreationDateBefore().IsZero() {
		return nil
	}
	query.CreatedBefore = NewFilter(FieldCreationDate, builder.GetCreationDateBefore(), OperationLess)
	return query.CreatedBefore
}

func resourceOwnerFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetResourceOwner() == "" {
		return nil
	}
	query.Owner = NewFilter(FieldResourceOwner, builder.GetResourceOwner(), OperationEquals)
	return query.Owner
}

func editorUserFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetEditorUser() == "" {
		return nil
	}
	query.Creator = NewFilter(FieldEditorUser, builder.GetEditorUser(), OperationEquals)
	return query.Creator
}

func instanceIDFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetInstanceID() == nil {
		return nil
	}
	query.InstanceID = NewFilter(FieldInstanceID, *builder.GetInstanceID(), OperationEquals)
	return query.InstanceID
}

func positionAfterFilter(builder *eventstore.SearchQueryBuilder, query *SearchQuery) *Filter {
	if builder.GetPositionAfter() == 0 {
		return nil
	}
	query.Position = NewFilter(FieldPosition, builder.GetPositionAfter(), OperationGreater)
	return query.Position
}

func aggregateIDFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetAggregateIDs()) < 1 {
		return nil
	}
	if len(query.GetAggregateIDs()) == 1 {
		return NewFilter(FieldAggregateID, query.GetAggregateIDs()[0], OperationEquals)
	}
	return NewFilter(FieldAggregateID, database.TextArray[string](query.GetAggregateIDs()), OperationIn)
}

func eventTypeFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetEventTypes()) < 1 {
		return nil
	}
	if len(query.GetEventTypes()) == 1 {
		return NewFilter(FieldEventType, query.GetEventTypes()[0], OperationEquals)
	}
	eventTypes := make(database.TextArray[eventstore.EventType], len(query.GetEventTypes()))
	for i, eventType := range query.GetEventTypes() {
		eventTypes[i] = eventType
	}
	return NewFilter(FieldEventType, eventTypes, OperationIn)
}

func aggregateTypeFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetAggregateTypes()) < 1 {
		return nil
	}
	if len(query.GetAggregateTypes()) == 1 {
		return NewFilter(FieldAggregateType, query.GetAggregateTypes()[0], OperationEquals)
	}
	aggregateTypes := make(database.TextArray[eventstore.AggregateType], len(query.GetAggregateTypes()))
	for i, aggregateType := range query.GetAggregateTypes() {
		aggregateTypes[i] = aggregateType
	}
	return NewFilter(FieldAggregateType, aggregateTypes, OperationIn)
}

func eventDataFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetEventData()) == 0 {
		return nil
	}
	return NewFilter(FieldEventData, query.GetEventData(), OperationJSONContains)
}
