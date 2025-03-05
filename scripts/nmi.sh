#! /bin/bash

query_url='https://secure.nmi.com/api/query.php'
transact_url='https://secure.nmi.com/api/transact.php'

query() {
    curl --data "security_key=$NMI_TOKEN" $query_url
}

transact() {
    curl --data "security_key=$NMI_TOKEN&$1" $transact_url
}

transact $*
