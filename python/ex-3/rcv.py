#!/usr/bin/python

import socket
import sys
from datetime import datetime

UDP_IP = "127.0.0.1"
UDP_PORT = 30609 
FILE_NAME = '/tmp/sta.out'

if len(sys.argv) > 1:
  UDP_IP = sys.argv[1]
  FILE_NAME = '/tmp/%s' % sys.argv[2] 

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP
sock.bind((UDP_IP, UDP_PORT))

with open(FILE_NAME, 'w') as f:
  while True:
    data, addr = sock.recvfrom(1024) # buffer size is 1024 bytes
    f.write("rcvd msg: '%s', from: %s, at: %s\n" % 
      (data, addr, str(datetime.now())))

f.close()
