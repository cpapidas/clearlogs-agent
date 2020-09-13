# UNDER DEVELOPMENT
![Under construction](https://thumbs.gfycat.com/JoyfulAfraidAltiplanochinchillamouse-small.gif)

# Clearlog Agent

Clearlog Agent is responsible to listen the STOUT and STDERR of the 
processes and send the logs to a remote server.

# Development Notes

## Test it locally
 
Run the following command in order to build and run the docker image with 
an example http server

`docker build --rm -t clagent:v1 .`


# How to start the agent

## Using UDP
```
    $ cd cmd/clagent
    $ go build && ./clagent \
        -token=AR-c12d47f5-22a3-459d-bb43-bc1273fdeb97 \
        -udp=true \
        -ignoreKillSimilarProcess=false \ 
        -baseUrl=http://localhost:3001 
```

single line
`go build && ./clagent -token=AR-c12d47f5-22a3-459d-bb43-bc1273fdeb97 -udp=true -ignoreKillSimilarProcess=false -baseUrl=http://localhost:3001`

