package user

import (
	"context"
	"fmt"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/dana-team/provider-gitlab/apis/gitlab/v1alpha1"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strconv"
)

// +kubebuilder:webhook:path=/mutate-gitlab-dana-io-v1alpha1-user,mutating=true,failurePolicy=fail,sideEffects=None,groups=gitlab.dana.io,resources=users,verbs=create;update,versions=v1alpha1,name=mutate.user.gitlab.dana.io,admissionReviewVersions=v1

type UserDefaulter struct {
	Client client.Client
}

var _ admission.CustomDefaulter = &UserDefaulter{}

func (d *UserDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	user, ok := obj.(*v1alpha1.User)
	if !ok {
		return fmt.Errorf("expected a User but got a %T", obj)
	}

	// Skip if external name already set
	if meta.GetExternalName(user) != "" {
		return nil
	}

	gClient, err := NewClientFromProviderConfig(ctx, d.Client, user.Spec.ProviderConfigReference.Name)
	if err != nil {
		return fmt.Errorf("failed to create GitLab client: %w", err)
	}

	// Lookup user by username
	users, _, err := gClient.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: user.Spec.ForProvider.Username,
	})
	if err != nil {
		return fmt.Errorf("failed to lookup GitLab user: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("GitLab user %s not found", *user.Spec.ForProvider.Username)
	}

	// Set the external name to the found ID
	meta.SetExternalName(user, strconv.Itoa(users[0].ID))
	return nil
}

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	defaulter := &UserDefaulter{
		Client: mgr.GetClient(),
	}

	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.User{}).
		WithDefaulter(defaulter).
		Complete()
}
