#!/bin/bash


BASE_URL="http://localhost:8080"
ID="h890feb7b2b84b6test"


echo "Тест для GET /orders"
wrk -t12 -c400 -d30s "$BASE_URL/orders"


echo "Тест для GET /orders/:id"
wrk -t12 -c400 -d30s "$BASE_URL/orders/$ID"


echo "Тест для POST /order/:id"
wrk -t1 -c1 -d30s -s post.lua -- "$BASE_URL/orders" 

