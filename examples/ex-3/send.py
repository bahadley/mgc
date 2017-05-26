#!/usr/bin/python

import socket
import time

UDP_IP = "10.0.0.3"
#UDP_IP = "127.0.0.1"
UDP_PORT = 30609 
MESSAGE = "Alive"

print "UDP target IP:", UDP_IP
print "UDP target port:", UDP_PORT
print "message:", MESSAGE

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP

while True:
  sock.sendto(MESSAGE, (UDP_IP, UDP_PORT))
  time.sleep(1)
