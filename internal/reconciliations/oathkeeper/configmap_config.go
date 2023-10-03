package gateway

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/kyma-project/api-gateway/apis/operator/v1alpha1"
	"github.com/kyma-project/api-gateway/internal/reconciliations"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed configmap_config.yaml
var configmapConfig []byte

const configMapName = "ory-oathkeeper-config"

func reconcileOryOathkeeperConfigConfigMap(ctx context.Context, k8sClient client.Client, apiGatewayCR v1alpha1.APIGateway) error {
	ctrl.Log.Info("Reconciling Ory Config ConfigMap", "name", configMapName, "namespace", namespace)

	if apiGatewayCR.IsInDeletion() {
		return deleteConfigmap(k8sClient, configMapName, namespace)
	}

	templateValues := make(map[string]string)
	templateValues["Name"] = configMapName
	templateValues["Namespace"] = namespace

	return reconciliations.ApplyResource(ctx, k8sClient, configmapConfig, templateValues)
}

func deleteConfigmap(k8sClient client.Client, name, namespace string) error {
	ctrl.Log.Info("Deleting Oathkeeper ConfigMap if it exists", "name", name, "namespace", namespace)
	s := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := k8sClient.Delete(context.Background(), &s)

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Oathkeeper ConfigMap %s/%s: %v", namespace, name, err)
	}

	ctrl.Log.Info("Successfully deleted Oathkeeper ConfigMap", "name", name, "namespace", namespace)

	return nil
}