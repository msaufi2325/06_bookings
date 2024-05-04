#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=postgres -dbpass=password -cache=false -production=false

# chmod +x run.sh
# ./run.sh
