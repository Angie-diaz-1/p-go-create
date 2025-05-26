#!/bin/bash

set -e  # Detener en caso de error

# Cargar variables del entorno desde el archivo .env
source .env

NAMESPACE="p-go-create"
CLUSTER_NAME="parcial-cluster"

echo "🔧 Verificando si el clúster $CLUSTER_NAME existe..."
if ! kind get clusters | grep -q "$CLUSTER_NAME"; then
  echo "🚀 Creando clúster $CLUSTER_NAME..."
  kind create cluster --name $CLUSTER_NAME
else
  echo "✅ El clúster $CLUSTER_NAME ya existe."
fi

echo -e "\n📁 Aplicando manifiestos YAML del microservicio en $NAMESPACE..."

# 1. Namespace
kubectl apply -f k8s/namespace.yml

# 2. Secret de acceso a GHCR desde variable $GHCR_TOKEN
echo "🔐 Creando el secreto ghcr-secret en el namespace $NAMESPACE..."
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=angie-diaz-1 \
  --docker-password="$GHCR_TOKEN" \
  --namespace=$NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# 3. Aplicar manifiestos
kubectl apply -f k8s/pv.yml
kubectl apply -f k8s/pvc.yml
kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/service.yml
kubectl apply -f k8s/ingress.yml
kubectl apply -f k8s/pod.yml
kubectl apply -f k8s/replicaset.yml

# 4. Esperar a que el Deployment esté listo
echo -e "\n⏳ Esperando que el deployment esté listo..."
kubectl wait --namespace=$NAMESPACE \
  --for=condition=available deployment/p-go-create-deployment \
  --timeout=120s || echo "⚠️  No se completó el deployment a tiempo"

# 5. Verificación
echo -e "\n📦 Estado del namespace $NAMESPACE:"
kubectl get all -n $NAMESPACE

echo -e "\n🌐 Ingress configurado:"
kubectl get ingress -n $NAMESPACE

echo -e "\n📂 Volúmenes:"
kubectl get pvc -n $NAMESPACE
kubectl get pv
