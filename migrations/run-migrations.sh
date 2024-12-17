#!/bin/bash

psql -U postgres -d postgres -f /migrations/000001_orders.up.sql
psql -U postgres -d postgres -f /migrations/000002_deliveries.up.sql
psql -U postgres -d postgres -f /migrations/000003_payments.up.sql
psql -U postgres -d postgres -f /migrations/000004_items.up.sql