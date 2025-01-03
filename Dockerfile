# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

ADD . /app/ 

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ARG CLUSTERENV="dev"

ENV clusterEnv=$CLUSTERENV

COPY *.go ./


RUN go build -o /docker-gs-ping

EXPOSE 3000
 
CMD [ "/docker-gs-ping" ]



#docker build -t service-news-app-backend:dev .
#docker build -t service-news-app-backend:prod .

# AWS Commands
#docker tag service-news-app-backend:dev 528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:dv222
#docker tag service-news-app-backend:prod 528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:pv26

#docker push 528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:dv222
#docker push 528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:pv26

#kubectl set image deployments/service-news-app-backend service-news-app-backend=528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:dv222
#kubectl set image deployments/service-news-app-backend service-news-app-backend=528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:pv26


## Test Deployment for dev
#kubectl set image deployments/service-news-app-backend-test service-news-app-backend=528458746222.dkr.ecr.ap-south-1.amazonaws.com/service-news-app-backend:dv222
