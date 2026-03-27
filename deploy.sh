#!/bin/bash

# Configuration
REMOTE_IP="10.0.0.206"
REMOTE_USER="barikhan" # Update this if the username is different

echo "--- Building Docker images ---"
docker build -t orders-api:latest --build-arg APP_NAME=orders-api .
docker build -t reports-api:latest --build-arg APP_NAME=reports-api .

echo "--- Saving images to tar files ---"
docker save orders-api:latest > orders-api.tar
docker save reports-api:latest > reports-api.tar

echo "--- Transferring images to $REMOTE_IP ---"
scp orders-api.tar reports-api.tar ${REMOTE_USER}@${REMOTE_IP}:.

echo "--- Cleaning up local tar files ---"
rm orders-api.tar reports-api.tar

echo "--- Deployment script finished ---"
echo "Next steps on the remote machine:"
echo "1. ssh ${REMOTE_USER}@${REMOTE_IP}"
echo "2. docker load < orders-api.tar"
echo "3. docker load < reports-api.tar"
echo "4. kubectl rollout restart deployment/orders-api"
echo "5. kubectl rollout restart deployment/reports-api"
