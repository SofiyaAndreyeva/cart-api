#!/bin/bash

# Usage:
# ./load-hpa.sh <namespace> <endpoint> [method] [json-data]
# Examples:
# POST:   ./load-hpa.sh dev "/carts" POST
# POST:   ./load-hpa.sh dev "/carts/1/items" POST '{"product":"apple","price":10.5}'
# DELETE: ./load-hpa.sh dev "/carts/1/items/1" DELETE
# GET:    ./load-hpa.sh dev "/carts/1" GET
# GET:    ./load-hpa.sh dev "/carts/1/price" GET

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
  echo "Error: specify namespace, endpoint, and method (GET, POST, DELETE)"
  echo "Example: ./load-hpa.sh dev '/carts/1/items' POST '{\"product\":\"apple\",\"price\":10.5}'"
  exit 1
fi

NAMESPACE=$1
ENDPOINT=$2
METHOD=$(echo "$3" | tr '[:lower:]' '[:upper:]')
DATA=$4

PORT=8080
while lsof -i :$PORT >/dev/null; do
  PORT=$((PORT+1))
done

SERVICE_URL="http://localhost:$PORT$ENDPOINT"
DURATION="30s"
CONCURRENCY=50

echo "Starting HPA demo in namespace: $NAMESPACE"
echo "Endpoint: $ENDPOINT, Method: $METHOD, Port: $PORT"

kubectl port-forward svc/cart-api $PORT:80 -n $NAMESPACE >/dev/null 2>&1 &
PF_PID=$!
sleep 2
echo "Port-forward started (PID=$PF_PID)"

kubectl get hpa -n $NAMESPACE -w &
HPA_PID=$!
echo "HPA monitoring started (PID=$HPA_PID)"

echo "Sending load to $SERVICE_URL"
case $METHOD in
  POST)
    if [ -n "$DATA" ]; then
      hey -z $DURATION -c $CONCURRENCY -m POST -d "$DATA" $SERVICE_URL
    else
      hey -z $DURATION -c $CONCURRENCY -m POST $SERVICE_URL
    fi
    ;;
  GET)
    hey -z $DURATION -c $CONCURRENCY -m GET $SERVICE_URL
    ;;
  DELETE)
    hey -z $DURATION -c $CONCURRENCY -m DELETE $SERVICE_URL
    ;;
  *)
    echo "Error: unknown method $METHOD. Use GET, POST, or DELETE."
    kill $PF_PID
    kill $HPA_PID
    exit 1
    ;;
esac

kill $PF_PID
kill $HPA_PID
echo "Load test finished, port-forward and HPA monitoring stopped"
