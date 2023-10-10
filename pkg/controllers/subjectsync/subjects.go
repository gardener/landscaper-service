package subjectsync

import (
	"context"
	"fmt"

	"github.com/gardener/landscaper/controller-utils/pkg/logging"
	rbacv1 "k8s.io/api/rbac/v1"

	lssv1alpha2 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha2"
)

// CreateSubjectsForSubjectList converts the subjects of the SubjectList custom resource into rbac subjects.
func CreateSubjectsForSubjectList(ctx context.Context, subjectList *lssv1alpha2.SubjectList) []rbacv1.Subject {
	logger, _ := logging.FromContextOrNew(ctx, nil)

	subjects := []rbacv1.Subject{}

	for _, subject := range subjectList.Spec.Subjects {
		rbacSubject, err := createSubjectForSubjectListEntry(subject)
		if err != nil {
			logger.Error(err, "could not create rbac.Subject from SubjectList.spec.subject")
			continue
		}
		subjects = append(subjects, *rbacSubject)
	}

	return subjects
}

// createSubjectForSubjectListEntry converts a single subject of the SubjectList custom resource into an rbac subject.
func createSubjectForSubjectListEntry(subjectListEntry lssv1alpha2.Subject) (*rbacv1.Subject, error) {
	switch subjectListEntry.Kind {
	case SUBJECT_LIST_ENTRY_USER, SUBJECT_LIST_ENTRY_GROUP:
		// if the entry has a namespace, we ignore it
		return &rbacv1.Subject{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     subjectListEntry.Kind,
			Name:     subjectListEntry.Name,
		}, nil
	case SUBJECT_LIST_ENTRY_SERVICE_ACCOUNT:
		// if the entry has no namespace, we use the LS_USER_NAMESPACE
		namespace := subjectListEntry.Namespace
		if namespace == "" {
			namespace = LS_USER_NAMESPACE
		}
		return &rbacv1.Subject{
			APIGroup:  "", //defaults to "" for service accounts as per rbacv1.Subject doc
			Kind:      subjectListEntry.Kind,
			Name:      subjectListEntry.Name,
			Namespace: namespace,
		}, nil
	default:
		return nil, fmt.Errorf("subject kind %s unknown", subjectListEntry.Kind)
	}
}
