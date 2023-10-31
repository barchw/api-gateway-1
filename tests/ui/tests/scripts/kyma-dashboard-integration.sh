#!/usr/bin/env bash

set -ex
export CYPRESS_DOMAIN=http://localhost:3001
export NO_COLOR=1
export REPO_IMG_DEV="k3d-registry.localhost:5000/kyma-dashboard"
export TAG="test-dev"

apt-get install -y gettext-base

function deploy_k3d_kyma (){
curl -Lo kyma https://storage.googleapis.com/kyma-cli-unstable/kyma-linux
chmod +x ./kyma

echo "Provisioning k3d cluster for Kyma"
./kyma provision k3d --ci

./kyma deploy

./kyma alpha deploy

echo "Apply and enable keda module"
kubectl apply -f https://github.com/kyma-project/keda-manager/releases/latest/download/moduletemplate.yaml

echo "Apply and enable serverless module"
kubectl apply -f https://github.com/kyma-project/serverless-manager/releases/latest/download/moduletemplate.yaml
./kyma alpha enable module serverless --channel fast

echo "Apply api-gateway"
kubectl apply -f https://github.com/kyma-project/api-gateway/releases/latest/download/api-gateway-manager.yaml
kubectl apply -f https://github.com/kyma-project/api-gateway/releases/latest/download/apigateway-default-cr.yaml

if [[ ${JOB_NAME} =~ .*smoke.* ]]; then
    echo "Apply and enable telemetry module"
    kubectl apply -f https://github.com/kyma-project/telemetry-manager/releases/latest/download/moduletemplate.yaml
    ./kyma alpha enable module telemetry --channel fast
fi

echo "Apply gardener resources"
echo "Certificates"
kubectl apply -f https://raw.githubusercontent.com/gardener/cert-management/master/pkg/apis/cert/crds/cert.gardener.cloud_certificates.yaml
echo "DNS Providers"
kubectl apply -f https://raw.githubusercontent.com/gardener/external-dns-management/master/pkg/apis/dns/crds/dns.gardener.cloud_dnsproviders.yaml
echo "DNS Entries"
kubectl apply -f https://raw.githubusercontent.com/gardener/external-dns-management/master/pkg/apis/dns/crds/dns.gardener.cloud_dnsentries.yaml
echo "Issuers"
kubectl apply -f https://raw.githubusercontent.com/gardener/cert-management/master/pkg/apis/cert/crds/cert.gardener.cloud_issuers.yaml

echo "Apply OAuth2 Hydra CRD"
kubectl apply -f https://raw.githubusercontent.com/ory/hydra-maester/master/config/crd/bases/hydra.ory.sh_oauth2clients.yaml

k3d kubeconfig get kyma > tests/fixtures/kubeconfig.yaml
}

function build_and_run_busola() {
echo "Create k3d registry..."
k3d registry create registry.localhost --port=5000

echo "Make release-dev..."
make release-dev

echo "Running kyma-dashboard..."
docker run -d --rm --net=host --pid=host --name kyma-dashboard "$REPO_IMG_DEV-local-dev:$TAG"

echo "waiting for server to be up..."
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$CYPRESS_DOMAIN")" != "200" ]]; do sleep 5; done
sleep 10
}

echo 'Waiting for deploy_k3d_kyma and build_and_run_busola'
wait -n
echo "First process finished"
wait -n
echo "Second process finished"

cd tests
npm ci && npm run "test:$SCOPE"