package oathkeeper

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/kyma-project/api-gateway/apis/operator/v1alpha1"
	"github.com/kyma-project/api-gateway/internal/clusterconfig"
	"github.com/kyma-project/api-gateway/internal/reconciliations"
	"github.com/kyma-project/api-gateway/internal/reconciliations/oathkeeper/maester"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	deploymentName = "ory-oathkeeper"
)

//go:embed deployment_light.yaml
var deploymentLight []byte

//go:embed deployment.yaml
var deployment []byte

func reconcileOathkeeperDeployment(ctx context.Context, k8sClient client.Client, apiGatewayCR v1alpha1.APIGateway) error {

	clusterSize, err := clusterconfig.EvaluateClusterSize(ctx, k8sClient)
	if err != nil {
		return err
	}

	ctrl.Log.Info("Reconciling Ory Oathkeeper Deployment", "Cluster size", clusterSize, "name", deploymentName, "Namespace", reconciliations.Namespace)

	if apiGatewayCR.IsInDeletion() {
		return deleteDeployment(ctx, k8sClient, deploymentName)
	}

	if clusterSize == clusterconfig.Evaluation {
		return reconcileDeployment(ctx, k8sClient, deploymentName, &deploymentLight)
	}
	return reconcileDeployment(ctx, k8sClient, deploymentName, &deployment)
}

func reconcileDeployment(ctx context.Context, k8sClient client.Client, name string, deploymentManifest *[]byte) error {

	ctrl.Log.Info("Reconciling Deployment", "name", name, "Namespace", reconciliations.Namespace)
	templateValues := make(map[string]string)
	templateValues["Name"] = name
	templateValues["Namespace"] = reconciliations.Namespace
	templateValues["ServiceAccountName"] = maester.ServiceAccountName

	return reconciliations.ApplyResource(ctx, k8sClient, *deploymentManifest, templateValues)
}

func deleteDeployment(ctx context.Context, k8sClient client.Client, name string) error {
	ctrl.Log.Info("Deleting Deployment if it exists", "name", name, "Namespace", reconciliations.Namespace)
	c := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: reconciliations.Namespace,
		},
	}
	err := k8sClient.Delete(ctx, &c)

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Deployment %s/%s: %v", reconciliations.Namespace, name, err)
	}

	if k8serrors.IsNotFound(err) {
		ctrl.Log.Info("Skipped deletion of Deployment as it wasn't present", "name", name, "Namespace", reconciliations.Namespace)
	} else {
		ctrl.Log.Info("Successfully deleted Deployment", "name", name, "Namespace", deploymentName)
	}

	return nil
}

func waitForDeploymentToBeReady(ctx context.Context, k8sClient client.Client) error {
	return retry.Do(func() error {
		var dep appsv1.Deployment
		err := k8sClient.Get(ctx, client.ObjectKey{
			Namespace: reconciliations.Namespace,
			Name:      deploymentName,
		}, &dep)

		if err != nil {
			return err
		}
		ctrl.Log.Info("replicas:", "unavailable replicas", dep.Status.UnavailableReplicas)

		if dep.Status.UnavailableReplicas > 0 {
			return fmt.Errorf("unavailable replicas %d", dep.Status.UnavailableReplicas)
		}

		return nil
	}, retry.Attempts(60), retry.Delay(2*time.Second), retry.DelayType(retry.FixedDelay))
}
