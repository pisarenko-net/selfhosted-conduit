# selfhosted-conduit
A kind of reverse proxy server intended to connect clients to self hosted backends that cannot listen for incoming connections

Register backend with the conduit
curl --cert certs/client.pem --key certs/client.key -k -X PUT https://0:8090/backend/register

Establish persistent HTTP connection to conduit
curl --cert certs/client.pem --key certs/client.key -k -X POST https://0:8090/backend/connect --data '69YL9I'
