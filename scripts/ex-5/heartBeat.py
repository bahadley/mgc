#!/usr/bin/python

from time import sleep, time
from signal import SIGINT

from mininet.cli import CLI
from mininet.link import TCLink
from mininet.log import setLogLevel
from mininet.net import Mininet
from mininet.node import Controller, OVSKernelAP
from mininet.util import pmonitor

OUTPUT_FILE = '/tmp/mgc.out'
SYNC_START = 10
EXPERIMENT_DURATION = 60 
EXECUTABLE_PATH = '../../mgc'

def topology():
    "Create a network."
    net = Mininet(controller=Controller, link=TCLink, accessPoint=OVSKernelAP)

    print "*** Creating nodes"
    sta1 = net.addStation('sta1')
    sta2 = net.addStation('sta2')
    sta3 = net.addStation('sta3')
    ap1 = net.addAccessPoint('ap1', ssid='ssid-ap1', mode='g', channel='11', position='100,100,0', range=100)
    c1 = net.addController('c1', controller=Controller)

    print "*** Configuring wifi nodes"
    net.configureWifiNodes()

    sta1.setIP('10.0.0.1/8', intf="sta1-wlan0")
    sta2.setIP('10.0.0.2/8', intf="sta2-wlan0")
    sta3.setIP('10.0.0.3/8', intf="sta3-wlan0")

    print "*** Starting network"
    net.build()
    c1.start()
    ap1.start([c1])

    net.plotGraph(max_x=200, max_y=200)

    net.associationControl('ssf')

    net.startMobility(time=0)
    net.mobility(sta1, 'start', time=1, position='8,100,0')
    net.mobility(sta2, 'start', time=1, position='1,107,0')
    net.mobility(sta3, 'start', time=1, position='1,93,0')
    net.mobility(sta1, 'stop', time=250, position='199,100,0')
    net.mobility(sta2, 'stop', time=250, position='192,107,0')
    net.mobility(sta3, 'stop', time=250, position='192,93,0')
    net.stopMobility(time=250)

    sleep(10)

    s1 = net.get('sta1')
    s2 = net.get('sta2')
    s3 = net.get('sta3')
    popens = {} # Python subprocess.Popen objects keyed by Mininet hosts
    startTime = int(time()) + SYNC_START
    endTime = startTime + EXPERIMENT_DURATION 

    print "*** Starting %d second experiment in %d second(s) - at Unix epoch: %d..." % (
      EXPERIMENT_DURATION, (startTime - int(time())), startTime)

    popens[s1] = s1.popen(EXECUTABLE_PATH, '-role=l', '-addr=%s' % s1.IP(), 
      '-dsts=%s,%s' % (s2.IP(), s3.IP()), '-start=%d' % startTime)
    popens[s2] = s2.popen(EXECUTABLE_PATH, '-role=f', 
      '-addr=%s' % s2.IP(), '-start=%d' % startTime)
    popens[s3] = s3.popen(EXECUTABLE_PATH, '-role=f', 
      '-addr=%s' % s3.IP(), '-start=%d' % startTime)

    with open(OUTPUT_FILE, 'w') as f:
      for h, line in pmonitor(popens, timeoutms=500):
        if h:
          f.write('<%s>: %s' % (h.name, line))
        if time() >= endTime:
          break

    popens[s1].send_signal(SIGINT)
    popens[s2].send_signal(SIGINT)
    popens[s3].send_signal(SIGINT)

    f.close()

    print "*** Ending experiment..."

    print "*** Running CLI"
    CLI(net)

    print "*** Stopping network"
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    topology()
