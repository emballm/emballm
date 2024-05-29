FROM golang:1.22.3 as base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build


FROM golang:1.22.3 as final
LABEL authors="andershokinson,andrewmollohan,aubreyklaft,gabrielaboy,williamwisseman"

#copy executable and run scripts
COPY --from=base /app/emballm /bin/emballm
COPY --from=base --chmod=755 /app/scripts/*.sh /bin/

# Install curl and ollama
RUN apt update && apt upgrade && apt install -y curl
RUN curl -fsSL https://ollama.com/install.sh | sh

ENV EMBALLM_EXECUTABLE=/bin/emballm

ENTRYPOINT ["/bin/startOllama.sh"]
CMD [ "-directory", "/scratch/harness" ]