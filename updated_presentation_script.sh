#!/bin/bash





# Ensure current directory ends with 'test-network'
if [[ "$(pwd)" != */test-network ]]; then
  echo "ERROR: This script must be run from the 'test-network' directory."
  echo "Please 'cd' into the correct directory of hyperledger, e.g.:"
  echo "  Example cd ~/Hyperledger/fabric/fabric-samples/test-network/"
  exit 1
fi


providers=()
users=()
sources=()
destinations=()

# Fill users and providers array
for i in $(seq -w 1 20); do
    providers+=("provider${i}@gmail.com")
    users+=("user${i}@gmail.com")
    sources+=("src${i}")
    destinations+=("dest${i}")
done

restart() {
    echo "===== Restarting Chaincode and Network ====="
    ./network.sh down
    ./network.sh up createChannel -c ticketsystem

    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    ./network.sh deployCC -ccn keyvalchaincode -ccp ../asset-transfer-basic/chaincode-go -ccl go -c ticketsystem
}

invoke(){
    ./network.sh deployCC -ccn keyvalchaincode -ccp ../asset-transfer-basic/chaincode-go -ccl go -c ticketsystem
}

create_20_users() {
    echo "===== Creating 20 Users ====="
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    for gmail in "${users[@]}"; do
        name="User ${gmail//[^0-9]/}"  # Extracting number from gmail
        phone="1234567890"
        id="${gmail}"

        echo "Creating user ${gmail}................"
        result=$(peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateEntity\",\"Args\":[\"${gmail}\", \"${name}\", \"${phone}\", \"${id}\", \"user\"]}")
    done
}


create_20_providers() {
    echo "===== Creating 20 Providers ====="
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    for gmail in "${providers[@]}"; do
        name="Provider ${gmail//[^0-9]/}"  # Extracting number from gmail
        phone="1234567890"
        id="${gmail}"

        echo "Creating provider ${gmail}..............."
        result=$(peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateEntity\",\"Args\":[\"${gmail}\", \"${name}\", \"${phone}\", \"${id}\", \"provider\"]}")
    done
}


create_transportations() {
    echo "===== Creating Transportation Services ====="
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    modes=("Bus" "Train" "Flight")
    base_prices=(500 300 1000)

    for src in "${sources[@]}"; do
        for dest in "${destinations[@]}"; do
            [ "$src" == "$dest" ] && continue

            echo "Creating transports from $src to $dest"
            flag=true
            for ((k=1; k<=8; k++)); do
                index1=$((k*2))
                index2=$((k*2-1))
                index=$index1
                if [ "$flag" == true ]; then
                    index=$index2
                    flag=false
                else
                    flag=true
                fi

                provider_gmail="${providers[$index-1]}"
                vehicle_id="VHC-${src}-${dest}-$k"
                dept_time="06:00"
                arr_time="11:30"
                duration="5h30m"

                mode_index=$((RANDOM % 3))
                mode="${modes[$mode_index]}"
                base_price="${base_prices[$mode_index]}"
                total_seats="100"

                random_extra=$(( RANDOM % 100 ))
                price=$(( base_price + random_extra ))

                start_date="2025-05-01"
                end_date="2030-05-30"

                peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{
                    \"function\":\"AddTransportService\",
                    \"Args\":[
                        \"${provider_gmail}\",
                        \"${vehicle_id}\",
                        \"${src}\",
                        \"${dest}\",
                        \"${dept_time}\",
                        \"${arr_time}\",
                        \"${duration}\",
                        \"${mode}\",
                        \"${total_seats}\",
                        \"${price}\",
                        \"${start_date}\",
                        \"${end_date}\"
                    ]
                }"
            done
        done
    done
}

#till here the script is working 

check_status_user() {
    echo "Enter User gmail (e.g., user01@gmail.com): "
    read user_gmail
    
    echo "Checking status for user: $user_gmail"
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true


    
    peer chaincode query -C ticketsystem -n keyvalchaincode -c "{
        \"function\":\"GetDetailUser\",
        \"Args\":[\"$user_gmail\"]
    }" | jq

  
}

check_status_provider() {
    echo "Enter Provider gmail (e.g., provider01@gmail.com): "
    read provider_gmail
    
    echo "Checking status for provider: $provider_gmail"
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
    
    peer chaincode query -C ticketsystem -n keyvalchaincode -c "{
        \"function\":\"GetDetailProvider\",
        \"Args\":[\"$provider_gmail\"]
    }" | jq 
}


add_balance() {
  echo "Enter User gmail (e.g., user01@gmail.com): "
  read user_gmail

    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  
  echo "Enter Amount to Add (e.g., 5000): "
  read amount
  
  echo "Adding balance $amount to user: $user_gmail"
  
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{
    \"function\":\"AddBalance\",
    \"Args\":[\"$user_gmail\", \"$amount\"]
  }"
}

check_ticket_status() {
  echo "Enter Ticket ID (e.g., TICKET12345): "
  read ticket_id
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  
  echo "Checking status for ticket: $ticket_id"
  
   peer chaincode query -C ticketsystem -n keyvalchaincode -c "{
    \"function\":\"GetDetailTicket\",
    \"Args\":[\"$ticket_id\"]
  }" | jq
}

check_transport_status(){
    echo "Enter the transportID (e.g. transport-VHC-src20-dest20-1-src20-dest20)"
    read transport_id
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true

    echo "Checking transport status: ${transport_id}"
    
    peer chaincode query -C ticketsystem -n keyvalchaincode -c "{
    \"function\":\"GetTransportationDetail\",
    \"Args\":[\"$transport_id\"]
  }" | jq
}


enquire_available_transports() {
  echo "Enter Source (e.g., Kanpur): "
  read source
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  
  echo "Enter Destination (e.g., Delhi): "
  read destination
  
  echo "Enter Date (YYYY-MM-DD, e.g., 2025-05-15): "
  read date
  
  echo "Enter Mode (e.g., Bus, Train, Flight): "
  read mode
  
  echo "Querying available transports from $source to $destination on $date via $mode..."
  
   peer chaincode query -C ticketsystem -n keyvalchaincode -c "{
    \"Args\": [\"GetAvailableTransports\", \"$source\", \"$destination\", \"$date\", \"$mode\"]
  }" | jq
}


book_ticket() {
  echo "Enter User gmail (e.g., dave@gmail.com): "
  read user_gmail
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  echo "Enter Transport ID (e.g., transport-UP78DJ2025-Kanpur-Delhi): "
  read transport_id
  
  echo "Enter Travel Date (YYYY-MM-DD, e.g., 2025-07-10): "
  read travel_date
  
  echo "Enter Seat Number to Book (e.g., 1): "
  read seat_number
  
  echo "Booking ticket for $user_gmail on $transport_id for $travel_date with $seat_count seats..."
  
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{
    \"function\":\"BookTicket\",
    \"Args\":[\"$user_gmail\", \"$transport_id\", \"$travel_date\", \"$seat_number\"]
  }"
}

update_ticket() {
  echo "Enter Ticket ID (e.g., transport-UP78DJ2025-Kanpur-Delhi-2025-05-15-1): "
  read ticket_id
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  echo "Enter New Travel Date (YYYY-MM-DD, e.g., 2025-07-10): "
  read new_date
  
  echo "Enter New Seat Number (e.g., 2): "
  read new_seat_number
  
  echo "Updating ticket $ticket_id to new date $new_date with $new_seat_count seats..."
  
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{
    \"function\":\"UpdateTicket\",
    \"Args\":[\"$ticket_id\", \"$new_date\", \"$new_seat_number\"]
  }"
}

cancel_ticket() {
  echo "Enter User gmail (e.g., dave@gmail.com): "
  read user_gmail
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
  echo "Enter Ticket ID (e.g., transport-UP78DJ2025-Kanpur-Delhi-2025-07-10-2): "
  read ticket_id
  
  echo "Cancelling ticket $ticket_id for user $user_gmail..."
  
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C ticketsystem -n keyvalchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{
    \"function\":\"CancelTicket\",
    \"Args\":[\"$user_gmail\", \"$ticket_id\"]
  }"
}

if [[ $# -gt 0 ]]; then
  "$@"
fi
