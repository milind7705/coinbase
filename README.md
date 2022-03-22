# coinbase
The goal of this project is to create a real-time VWAP (volume-weighted average price) calculation engine. The calculations this project is following is as per the following blog:

https://analyzingalpha.com/blog/vwap

The gorilla socket has been used here as per following comparison:
![image](https://user-images.githubusercontent.com/5960727/159239386-d3525a2a-9ef3-42a6-bc6d-369e771dd90b.png)


For building the project, use 

**make build**


The build command creates the binary as bin/vwap which can be run as
** ./bin/vwap**

<img width="893" alt="image" src="https://user-images.githubusercontent.com/5960727/159240161-2e220f40-b14b-467b-9f11-14aaf77a30e2.png">

If the binary is run without any parameters, the **./bin/vwap** takes the default parameters defined as per problem. 
The binary can also be run using yaml file with command **./bin/vwap <yaml_file>** . The example yaml file is present in config/coinbase.yaml

Approach:

The code uses a gorilla websocket client which initiates the socket connection with the exchange(coinbase in this case) and performs the initial handshake via subscription. Another goroutine is listening on the channel for receiving the client response and initializes internal data structures(Queue).The Queue stores the individual trades with a sliding window of 200, summation of quantity and vwap at a given point in time.

For unit and integration testing: 

**make test**

