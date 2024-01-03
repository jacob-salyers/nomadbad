#!/bin/bash

send_email() {
    curl -v --user "api:$(cat cred/mailgun_api.txt)"                \
        https://api.mailgun.net/v3/mg.nomad-jiujitsu.com/messages \
        -F from='Nomad Form <mailgun@mg.nomad-jiujitsu.com>' \
        -F to=caravancollective@outlook.com \
        -F subject='Test 3' \
        -F text='Mailgun test'
}

send_email
