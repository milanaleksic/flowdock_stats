# Flowdock statistics

![status of the build on Travis](https://travis-ci.org/milanaleksic/flowdock_stats.svg "")

## Idea

**work under progress**

This is a prototype of a digesting small application written in Go language that digests Flowdock REST API.
 
The only _statistics_ that it gives is "who is the biggest talker" during a particular time span.

To avoid fetching all the users all the time, it will serialize locally. Since it uses proto3 which is currently in BETA
you will need to do following steps until GA of proto3 is announced:
- to download suitable compiled protobuf compiler package from https://github.com/google/protobuf/releases
- install Go plugin for protobuf compiler https://github.com/grpc/grpc-go/tree/master/examples 
 
## How to install

If you have Go SDK (1.5.1 is the one I used), all you need to do is

    go get github.com/milanaleksic/flowdock_stats

## How to run

Following arguments are needed for application to know which flow it should access and parse through:

    flowdock_stats -flowdockApiToken=API_TOKEN -companyToAnalyze=COMPANY_NAME -flowToAnalyze=FLOW_NAME
    
Both `COMPANY_NAME` and `FLOW_NAME` can be read from the URL of a flow, currently it's sth like: 
https://www.flowdock.com/app/COMPANY_NAME/FLOW_NAME 

And the `API_TOKEN` can be taken from the page https://www.flowdock.com/account/tokens 