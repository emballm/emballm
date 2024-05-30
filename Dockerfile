FROM --platform=linux/amd64 ollama/ollama:latest as final

WORKDIR /bin

ENV OLLAMA_MODELS=/app/models
COPY scratch/models /app/models

COPY scratch/emballm /bin/emballm
COPY --chmod=755 scripts/*.sh /bin/
COPY config.yaml /bin/config.yaml

ENV EMBALLM_EXECUTABLE=/bin/emballm

RUN apt update && apt install -y curl

ENTRYPOINT ["/bin/start_ollama.sh"]
CMD [ "-directory", "/harness", "-config",  "/bin/config.yaml", "-model", "gemma:7b", "-output", "/addon/output.json"]
