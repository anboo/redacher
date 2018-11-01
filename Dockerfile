FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/PuerkitoBio/goquery && go get -u github.com/vorkytaka/easyvk-go/easyvk && go get github.com/joho/godotenv && go build -o main .
CMD ["/app/main"]