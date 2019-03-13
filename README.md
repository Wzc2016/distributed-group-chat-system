# distributed-group-chat-system

This is a course project of CS425 in UIUC.

The document is available [here](https://courses.engr.illinois.edu/ece428/sp2019/mps/mp1.html)

## Install and run 

	go build main.go

	./main <name> <port> <n>

Below is a sample run of the program:

<pre><samp>$ ./main Alice 4444 3
READY
<kbd>Hi everyone!</kbd>
Bob: Hi Alice!
Charlie: Hey!
Bob: whoops, gotta go
Bob has left
</samp></pre>

You may **NOT** run this program on your own machine. You can only run this program on our VMs. 

**OR** you can modify the `hosts` to be a list of your own machines's hostnames.

## Report

You can view the design of the system in [report.pdf](https://blog.tsugumi.pro/pdf/report.pdf).
