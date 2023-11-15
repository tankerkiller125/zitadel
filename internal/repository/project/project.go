package project

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

const (
	UniqueProjectnameType  = "project_names"
	projectEventTypePrefix = eventstore.EventType("project.")
	ProjectAddedType       = projectEventTypePrefix + "added"
	ProjectChangedType     = projectEventTypePrefix + "changed"
	ProjectDeactivatedType = projectEventTypePrefix + "deactivated"
	ProjectReactivatedType = projectEventTypePrefix + "reactivated"
	ProjectRemovedType     = projectEventTypePrefix + "removed"
)

func NewAddProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueProjectnameType,
		projectName+resourceOwner,
		"Errors.Project.AlreadyExists")
}

func NewRemoveProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueProjectnameType,
		projectName+resourceOwner)
}

type ProjectAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                   string                        `json:"name,omitempty"`
	ProjectRoleAssertion   bool                          `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck       bool                          `json:"projectRoleCheck,omitempty"`
	HasProjectCheck        bool                          `json:"hasProjectCheck,omitempty"`
	PrivateLabelingSetting domain.PrivateLabelingSetting `json:"privateLabelingSetting,omitempty"`
}

func (e *ProjectAddedEvent) Payload() interface{} {
	return e
}

func (e *ProjectAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddProjectNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func NewProjectAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	projectRoleAssertion,
	projectRoleCheck,
	hasProjectCheck bool,
	privateLabelingSetting domain.PrivateLabelingSetting,
) *ProjectAddedEvent {
	return &ProjectAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectAddedType,
		),
		Name:                   name,
		ProjectRoleAssertion:   projectRoleAssertion,
		ProjectRoleCheck:       projectRoleCheck,
		HasProjectCheck:        hasProjectCheck,
		PrivateLabelingSetting: privateLabelingSetting,
	}
}

func ProjectAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ProjectAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-Bfg2f", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectChangeEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                   *string                        `json:"name,omitempty"`
	ProjectRoleAssertion   *bool                          `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck       *bool                          `json:"projectRoleCheck,omitempty"`
	HasProjectCheck        *bool                          `json:"hasProjectCheck,omitempty"`
	PrivateLabelingSetting *domain.PrivateLabelingSetting `json:"privateLabelingSetting,omitempty"`
	oldName                string
}

func (e *ProjectChangeEvent) Payload() interface{} {
	return e
}

func (e *ProjectChangeEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.Name != nil {
		return []*eventstore.UniqueConstraint{
			NewRemoveProjectNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
			NewAddProjectNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
		}
	}
	return nil
}

func NewProjectChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldName string,
	changes []ProjectChanges,
) (*ProjectChangeEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-mV9xc", "Errors.NoChangesFound")
	}
	changeEvent := &ProjectChangeEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectChangedType,
		),
		oldName: oldName,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type ProjectChanges func(event *ProjectChangeEvent)

func ChangeName(name string) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.Name = &name
	}
}

func ChangeProjectRoleAssertion(projectRoleAssertion bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.ProjectRoleAssertion = &projectRoleAssertion
	}
}

func ChangeProjectRoleCheck(projectRoleCheck bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.ProjectRoleCheck = &projectRoleCheck
	}
}

func ChangeHasProjectCheck(ChangeHasProjectCheck bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.HasProjectCheck = &ChangeHasProjectCheck
	}
}

func ChangePrivateLabelingSetting(ChangePrivateLabelingSetting domain.PrivateLabelingSetting) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.PrivateLabelingSetting = &ChangePrivateLabelingSetting
	}
}

func ProjectChangeEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ProjectChangeEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-M9osd", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ProjectDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *ProjectDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewProjectDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ProjectDeactivatedEvent {
	return &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectDeactivatedType,
		),
	}
}

func ProjectDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ProjectReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ProjectReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *ProjectReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewProjectReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ProjectReactivatedEvent {
	return &ProjectReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectReactivatedType,
		),
	}
}

func ProjectReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &ProjectReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ProjectRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                     string
	entityIDUniqueContraints []*eventstore.UniqueConstraint
}

func (e *ProjectRemovedEvent) Payload() interface{} {
	return nil
}

func (e *ProjectRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	constraints := []*eventstore.UniqueConstraint{NewRemoveProjectNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
	if e.entityIDUniqueContraints != nil {
		for _, constraint := range e.entityIDUniqueContraints {
			constraints = append(constraints, constraint)
		}
	}
	return constraints
}

func NewProjectRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	entityIDUniqueContraints []*eventstore.UniqueConstraint,
) *ProjectRemovedEvent {
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectRemovedType,
		),
		Name:                     name,
		entityIDUniqueContraints: entityIDUniqueContraints,
	}
}

func ProjectRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
