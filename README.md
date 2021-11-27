Commands to execute in terminal start after "-->"
For example, for the command "ls", in this file will be "-->ls"


Install Go on Ubuntu
-->sudo apt install golang


Check the version (This project was written on version 1.17.3)
-->go version


Install libraries for Gioui (Go's graphic's library)
-->apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev


Build go executable file (terminal should be in the same folder as Main.go)
-->go build -o="TOMS" .


Run TOMS (Totally Ordered Multicast System) manager server
-->./TOMS -server


Server will log the address on which it is listening.
To run a worker, write that addres in place of ADDRESS in the following script
-->./TOMS -mAddr="ADDRESS"


Start all workers before sending any messages