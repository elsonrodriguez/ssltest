#!/bin/bash
kubectl create secret generic sslcerts --from-file=./server.key --from-file=./server.crt
