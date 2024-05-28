#!/bin/bash

# check and make sure ollama is installed and running
which ollama
ollama --version
ollama serve &
# TODO should poll for ollama to be available
sleep 30
ollama pull gemma:2b

# run emballm
. /bin/startEmballm.sh "$@"