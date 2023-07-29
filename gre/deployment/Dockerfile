FROM greatestworks

WORKDIR business/server/world

RUN go build main.go.go

EXPOSE 1101

RUN chmod +x main.go

ENTRYPOINT ["./main"]


