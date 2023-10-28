#! /bin/bash

query_url='https://secure.nmi.com/api/query.php'

curl --data "security_key=$NMI_TOKEN" $query_url
