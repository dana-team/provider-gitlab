package user

import (
	"context"
	"fmt"
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/dana-team/provider-gitlab/apis/v1beta1"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewClientFromProviderConfig(ctx context.Context, k8sClient client.Client, providerConfigName string) (*gitlab.Client, error) {
	// Get the ProviderConfig
	providerConfig := &v1beta1.ProviderConfig{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: providerConfigName}, providerConfig); err != nil {
		return nil, fmt.Errorf("failed to get ProviderConfig: %w", err)
	}

	// Get the token from the secret reference
	token, err := getSecretValue(ctx, k8sClient, *providerConfig.Spec.Credentials.SecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from secret: %w", err)
	}

	// Create GitLab client
	gClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(*providerConfig.Spec.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return gClient, nil
}

func getSecretValue(ctx context.Context, k8sClient client.Client, ref v1.SecretKeySelector) (string, error) {
	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}, secret); err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	tokenBytes, ok := secret.Data[ref.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret", ref.Key)
	}

	return string(tokenBytes), nil
}
