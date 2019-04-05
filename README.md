# email-events
A command line tool written in Golang.
It takes the command line arguments and does the following;
  a. creates the link from the provided arguments.
  b. using the created link it hits the elastic email api and return the status of the url in JASON format.
  3. extracts the url from the returned JASON and hits the internet to download the logs in .csv format.
  
  
# How to run the Program?
1. To build the program:
go build email.go

it will create the email executables, to run this executable just type;
./email -apiky=xxxx xxxx xxxx xxxx -statuses=1,2 from=2019-02-25T00:00:00 to=2019-03-03T23:59:59


2. To run without build:
go run email.go -apiky=xxxx xxxx xxxx xxxx -statuses=1,2 from=2019-02-25T00:00:00 to=2019-03-03T23:59:59

