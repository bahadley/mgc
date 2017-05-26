#!/usr/bin/python

import socket

UDP_IP = "10.0.0.3"
#UDP_IP = "127.0.0.1"
UDP_PORT = 30609 

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP
sock.bind((UDP_IP, UDP_PORT))

with open('/tmp/sta2.out', 'w') as f:
  while True:
    data, addr = sock.recvfrom(1024) # buffer size is 1024 bytes
    f.write("received message: %s from: %s\n" % (data, addr))

f.close()
