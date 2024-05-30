#!/usr/bin/env bash

which ollama
ollama --version
ollama serve &
sleep 5
while [[ "$(curl -s -o /dev/nul -w ''%{http_code}'' http://localhost:11434)" != "200" ]]; do
  sleep 5
done
echo ${EMBALLM_EXECUTABLE} "$@"
${EMBALLM_EXECUTABLE} "$@"