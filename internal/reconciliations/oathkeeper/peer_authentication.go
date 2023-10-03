package gateway

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/kyma-project/api-gateway/apis/operator/v1alpha1"
	"github.com/kyma-project/api-gateway/internal/reconciliations"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed peer_authentication.yaml
var peerAuthentication []byte

const peerAuthenticationName = "ory-oathkeeper-maester-metrics"

func reconcileOryOathkeeperPeerAuthentication(ctx context.Context, k8sClient client.Client, apiGatewayCR v1alpha1.APIGateway) error {
	ctrl.Log.Info("Reconciling Ory Config PeerAuthentication", "name", peerAuthenticationName, "namespace", namespace)

	if apiGatewayCR.IsInDeletion() {
		return deletePeerAuthentication(k8sClient, peerAuthenticationName, namespace)
	}

	templateValues := make(map[string]string)
	templateValues["Name"] = peerAuthenticationName
	templateValues["Namespace"] = namespace

	return reconciliations.ApplyResource(ctx, k8sClient, peerAuthentication, templateValues)
}

func deletePeerAuthentication(k8sClient client.Client, name, namespace string) error {
	ctrl.Log.Info("Deleting Oathkeeper PeerAuthentication if it exists", "name", name, "namespace", namespace)
	s := securityv1beta1.PeerAuthentication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := k8sClient.Delete(context.Background(), &s)

	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Oathkeeper PeerAuthentication %s/%s: %v", namespace, name, err)
	}

	ctrl.Log.Info("Successfully deleted Oathkeeper PeerAuthentication", "name", name, "namespace", namespace)

	return nil
}
