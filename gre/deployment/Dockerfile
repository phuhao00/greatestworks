FROM greatestworks

WORKDIR business/server/world

RUN go build main.go

EXPOSE 1101

RUN chmod +x main

ENTRYPOINT ["./main"]


