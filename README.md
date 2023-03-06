# selfhosted-conduit
A kind of reverse proxy server intended to connect clients to self hosted backends that cannot listen for incoming connections

### Send a message to the backend
curl -v --cert client.pem --key client.key -k -X POST https://0:8090/client/request/69YL9I -H "Content-Type: text/plain" --data 'aGVsbG8='

### Register backend with the conduit

curl --cert certs/client.pem --key certs/client.key -k -X PUT https://0:8090/backend/register

### Establish persistent HTTP connection to conduit

curl --cert certs/client.pem --key certs/client.key -k -X POST https://0:8090/backend/connect --data '69YL9I'
