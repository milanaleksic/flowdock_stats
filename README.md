# Flowdock statistics

[![Build Status](https://semaphoreci.com/api/v1/milanaleksic/flowdock_stats/branches/master/badge.svg)](https://semaphoreci.com/milanaleksic/flowdock_stats)

This is a small application written in Go language that digests Flowdock REST API and gives some numbers.
 
_Statistics_ that it calculates during a particular time span are:
- which user had how many comments on a chosen flow,
- percentage of them fixing (_editing_) their comments after sending the comment,
- number of words per comment.

To avoid fetching all the users all the time, it will serialize locally. Since it uses proto3 which is currently in BETA

## How to run

Following arguments are needed for application to know which flow it should access and parse through:

    flowdock_stats -flowdockApiToken=API_TOKEN -companyToAnalyze=COMPANY_NAME -flowToAnalyze=FLOW_NAME -days=DAYS
    
Both `COMPANY_NAME` and `FLOW_NAME` can be read from the URL of a flow, currently it's sth like: 
https://www.flowdock.com/app/COMPANY_NAME/FLOW_NAME 

And the `API_TOKEN` can be taken from the page https://www.flowdock.com/account/tokens.

If you don't want to look only for messages posted in last 24 hours but also longer, set `days` parameter to some value

## How to develop

If you have Go SDK (1.5.1 is the one I used), all you need to do is

    go get github.com/milanaleksic/flowdock_stats

_In case you wish to change the serialization format_ you will need to do following steps until GA of proto3 is announced:
- download suitable protobuf compiler package from https://github.com/google/protobuf/releases
- install Go plugin for protobuf compiler https://github.com/grpc/grpc-go/tree/master/examples 
- this will allow you to do `protoc --go_out=plugins=grpc:. *.proto` from within `serialization` directory
