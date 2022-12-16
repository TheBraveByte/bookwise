#pull the base image from docker hub
FROM golang:1.19.2

#Set the working directory to use
WORKDIR /usr/src/app

#Copy the golang dependencies file into the work directory
COPY go.mod .
COPY go.sum .

#download all the dependencies packages needed and verify the packages
RUN go mod download && go mod verify

#copy all the files in the root project directory
COPY . .

# build your code from the main folder/file in the parent folder
# to a executable file catalogueAPI and disabled using C code and
# use just the built-in
RUN CGO_ENABLED=0 go build -o bookwiseAPI ./web/cmd

#Run the executable
RUN chmod +x /usr/src/app/bookwiseAPI

#Compile the executable file
CMD ["/usr/src/app/bookwiseAPI"]

