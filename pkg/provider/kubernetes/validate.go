/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kubernetes

import (
	"context"
	"fmt"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/external-secrets/external-secrets/pkg/utils"
)

func (p *ProviderKubernetes) ValidateStore(store esv1beta1.GenericStore) error {
	storeSpec := store.GetSpec()
	k8sSpec := storeSpec.Provider.Kubernetes
	if k8sSpec.Server.CABundle == nil && k8sSpec.Server.CAProvider == nil {
		return fmt.Errorf("a CABundle or CAProvider is required")
	}
	if k8sSpec.Auth.Cert != nil {
		if k8sSpec.Auth.Cert.ClientCert.Name == "" {
			return fmt.Errorf("ClientCert.Name cannot be empty")
		}
		if k8sSpec.Auth.Cert.ClientCert.Key == "" {
			return fmt.Errorf("ClientCert.Key cannot be empty")
		}
		if err := utils.ValidateSecretSelector(store, k8sSpec.Auth.Cert.ClientCert); err != nil {
			return err
		}
	}
	if k8sSpec.Auth.Token != nil {
		if k8sSpec.Auth.Token.BearerToken.Name == "" {
			return fmt.Errorf("BearerToken.Name cannot be empty")
		}
		if k8sSpec.Auth.Token.BearerToken.Key == "" {
			return fmt.Errorf("BearerToken.Key cannot be empty")
		}
		if err := utils.ValidateSecretSelector(store, k8sSpec.Auth.Token.BearerToken); err != nil {
			return err
		}
	}
	if k8sSpec.Auth.ServiceAccount != nil {
		if err := utils.ValidateReferentServiceAccountSelector(store, *k8sSpec.Auth.ServiceAccount); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProviderKubernetes) Validate() (esv1beta1.ValidationResult, error) {
	// when using referent namespace we can not validate the token
	// because the namespace is not known yet when Validate() is called
	// from the SecretStore controller.
	if p.Namespace == "" {
		return esv1beta1.ValidationResultUnknown, nil
	}
	ctx := context.Background()
	authReview, err := p.ReviewClient.Create(ctx, &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Resource:  "secrets",
				Namespace: p.Namespace,
				Verb:      "get",
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return esv1beta1.ValidationResultUnknown, fmt.Errorf("could not verify if client is valid: %w", err)
	}
	if !authReview.Status.Allowed {
		return esv1beta1.ValidationResultError, fmt.Errorf("client is not allowed to get secrets")
	}
	return esv1beta1.ValidationResultReady, nil
}