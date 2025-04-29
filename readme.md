# MyTravel.com – Blockchain-Based Ticket Management System

## Overview
# Team Members:
- Deepanshu(210310)
- Narendra(210649)

This was our course project for CS731(Prof. Angshuman karmakar) at IIT Kanpur, where we explored the integration of blockchain technology with a web application.

MyTravel.com is a blockchain-based ticket management platform developed using Hyperledger Fabric.  
It enables customers to book, update, and cancel transport tickets securely, while service providers can register and manage transport services.

This project includes:
- Frontend: React.js
- Backend: Golang (API server)
- Blockchain Network: Hyperledger Fabric
- Database: MongoDB

[Full Project Report (CS731)](https://github.com/Deepanshu18099/Blockchain_cs731/)


Setup Instructions

## 🛠 Prerequisites

- Docker and Docker Compose
- Node.js and npm (for React frontend)
- Golang installed (>= 1.20)
- Hyperledger Fabric binaries installed
- MongoDB installed locally



# 1. Clone the Repository

```bash
cd ~{HyperledgerPath}/fabric-samples/asset-transfer-basic
git clone https://github.com/Deepanshu18099/Blockchain_cs731/ chaincode-go
```


# To Start the Hyperledger Fabric Network
### Set environment variables

```bash
export HyperledgerPath=~/fabric-samples-path
export CORE_PEER_TLS_ENABLED=true
export FABRIC_CFG_PATH=${HyperledgerPath}/fabric-samples/test-network/../config
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=${HyperledgerPath}/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=${HyperledgerPath}/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

### Start Fabric Network
${HyperledgerPath}/fabric-samples/test-network/network.sh up

### Create a Channel
${HyperledgerPath}/fabric-samples/test-network/network.sh createChannel -c ticketsystem

### Deploy the Chaincode
```
${HyperledgerPath}/fabric-samples/test-network/network.sh deployCC -ccn keyvalchaincode -ccp ${HyperledgerPath}/fabric-samples/asset-transfer-basic/chaincode-go -ccl go -c ticketsystem
```


## Run the Backend API Server
cd ./backend

### Install dependencies
go mod tidy

### Run backend server (main.go is inside /cmd)
cd cmd
go run main.go



## Run the Frontend (React App)
### Navigate to frontend directory
cd chaincode-go/frontend

### Install frontend dependencies
npm install

### Start the React app
npm start



# Useful Peer CLI Commands (for Hyperledger Fabric)

## Setup Chaincode
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
--tls --cafile ${HyperledgerPath}/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
-C ticketsystem -n keyvalchaincode \
--peerAddresses localhost:7051 \
--tlsRootCertFiles ${HyperledgerPath}/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
-c '{"function":"CreateUser","Args":["email","name","phone","userid"]}'
```
