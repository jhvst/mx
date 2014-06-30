MX checks domain DNS zone for MX records. This is useful if you have bunch of emails that you need to validate.

##How to use

1. `git clone https://github.com/9uuso/mx` -> `cd mx` -> `go build`
2. Insert text file called input.txt in the same directory as the compiled binary, with each email seperated by line ending.
3. Run the binary and wait for the program to finish.