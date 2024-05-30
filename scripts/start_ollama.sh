#!/usr/bin/env bash

which ollama
ollama --version
ollama serve > /addon/ollama.log &
curl --connect-timeout 5 \
  --max-time 10 \
  --retry 5 \
  --retry-delay 0 \
  --retry-max-time 40 \
  http://localhost:11434/api/generate -d '{"model": "gemma:7b"}'

echo ${EMBALLM_EXECUTABLE} "$@"
${EMBALLM_EXECUTABLE} "$@"