#!/bin/bash
which ollama
ollama --version
ollama serve &
# TODO should poll for ollama to be available
sleep 30
ollama pull gemma:2b

/bin/emballm --file test