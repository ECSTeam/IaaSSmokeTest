The purpose of this utility is to validate communication between two networks.  

To run:

Install go on VMs in both networks.

create an environment variable IAASTESTCONFIGDIR which points to the directory that hosts the "config.json" file.  

Modify the config.json file on each of the VMs to include:

1.  the hostname of the remote host
2.  the ports and protocols on which the remote host is expected to recieve messages
3.  the ports and protocols on which the local VM will need to recieve messages.

The localConnectionDetails section on one machine should match the remoteConnectionDetails on the other, and vice versa.

from a command line on each VM, type "go run check_readiness.go"

The process begins by setting up the listsners defined in the "localConnectionDetails" section of config.json.  If those are correctly initialized, press the enter key to continue.

The process next sends messages to the remote machine as configured the the "remoteHost" and "remoteConnectionDetails" sections of config.json.  You should be able to validate that messages were recieved remotely by watching the console.  Once the messages have been sent, again press enter to continue.

Lastly, a summary of events will appear.  You should expect to see that for every TCP listener instantiated, a message was recieved.  For every UDP listener instantiated, a message will be recieved.  And, for every TCP message sent, a reply was recieved.  If these numbers do not match, there should be error information in the console or there are simply errors in the configuration.

