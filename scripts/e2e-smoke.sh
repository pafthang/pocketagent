#!/usr/bin/env bash
set -euo pipefail

GATE_URL="${GATE_URL:-http://127.0.0.1:8080}"
EMAIL="${EMAIL:-smoke-$(date +%s)@example.com}"
PASSWORD="${PASSWORD:-secret12345}"

echo "==> register"
curl -sf -X POST "$GATE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" >/dev/null

echo "==> login"
TOKEN=$(curl -sf -X POST "$GATE_URL/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" | jq -r .token)

echo "==> create space"
SPACE_ID=$(curl -sf -X POST "$GATE_URL/spaces" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Smoke Space","slug":"smoke-'"$(date +%s)"'"}' | jq -r .id)

echo "==> create agent"
AGENT_ID=$(curl -sf -X POST "$GATE_URL/agents" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Space-Id: $SPACE_ID" \
  -H 'Content-Type: application/json' \
  -d '{"name":"smoke-agent","model":"llama3.1","tools":["search_web"]}' | jq -r .id)

echo "==> create task"
TASK_JSON=$(curl -sf -X POST "$GATE_URL/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Space-Id: $SPACE_ID" \
  -H 'Content-Type: application/json' \
  -d "{\"prompt\":\"say hello\",\"agent_id\":\"$AGENT_ID\"}")
CORR_ID=$(echo "$TASK_JSON" | jq -r .correlation_id)

echo "==> poll task status"
for i in $(seq 1 30); do
  STATUS=$(curl -sf "$GATE_URL/tasks/$CORR_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "X-Space-Id: $SPACE_ID" | jq -r .status)
  echo "   attempt $i: $STATUS"
  if [[ "$STATUS" == "completed" || "$STATUS" == "degraded" || "$STATUS" == "failed" ]]; then
    curl -sf "$GATE_URL/tasks/$CORR_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "X-Space-Id: $SPACE_ID" | jq .
    exit 0
  fi
  sleep 2
done

echo "timeout waiting for task completion"
exit 1