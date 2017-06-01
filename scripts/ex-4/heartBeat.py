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
EXPERIMENT_DURATION = 15
EXECUTABLE_PATH = '../../mgc'

def topology():
    "Create a network."
    net = Mininet(controller=Controller, link=TCLink, accessPoint=OVSKernelAP)

    print "*** Creating nodes"
    sta1 = net.addStation('sta1', wlans=2)
    sta2 = net.addStation('sta2', wlans=2)
    sta3 = net.addStation('sta3', wlans=2)
    ap1 = net.addAccessPoint('ap1', ssid='ssid-ap1', mode='g', channel='11', position='115,62,0')
    ap2 = net.addAccessPoint('ap2', ssid='ssid-ap2', mode='g', channel='1', position='57,142,0')
    c1 = net.addController('c1', controller=Controller)

    print "*** Configuring wifi nodes"
    net.configureWifiNodes()

    sta1.setIP('10.0.0.2/8', intf="sta1-wlan0")
    sta2.setIP('10.0.0.3/8', intf="sta2-wlan0")
    sta3.setIP('10.0.0.4/8', intf="sta3-wlan0")
    sta1.setIP('192.168.10.1/24', intf="sta1-wlan1")
    sta2.setIP('192.168.10.2/24', intf="sta2-wlan1")
    sta3.setIP('192.168.10.3/24', intf="sta3-wlan1")

    net.addHoc(sta1, ssid='adhocNet', mode='g')
    net.addHoc(sta2, ssid='adhocNet', mode='g')
    net.addHoc(sta3, ssid='adhocNet', mode='g')

    print "*** Starting network"
    net.build()
    c1.start()
    ap1.start([c1])
    ap2.start([c1])

    net.plotGraph(max_x=200, max_y=200)

    net.associationControl('ssf')

    net.startMobility(time=0)
    net.mobility(sta1, 'start', time=1, position='86,188,0')
    net.mobility(sta2, 'start', time=1, position='78,195,0')
    net.mobility(sta3, 'start', time=1, position='93,195,0')
    net.mobility(sta1, 'stop', time=250, position='86,0,0')
    net.mobility(sta2, 'stop', time=250, position='78,7,0')
    net.mobility(sta3, 'stop', time=250, position='93,7,0')
    net.stopMobility(time=250)

    sleep(10)

    print "*** Starting experiment..."
    s1 = net.get('sta1')
    s2 = net.get('sta2')
    s3 = net.get('sta3')

    popens = {}
    popens[s2] = s2.popen(EXECUTABLE_PATH, '-role=F', 
      '-addr=%s' % s2.IP())
    popens[s3] = s3.popen(EXECUTABLE_PATH, '-role=F', 
      '-addr=%s' % s3.IP())

    sleep(1)
    leader = s1.popen(EXECUTABLE_PATH, '-addr=%s' % s1.IP(), 
      '-dsts=%s,%s' % (s2.IP(), s3.IP()))

    endTime = time() + EXPERIMENT_DURATION 
    with open(OUTPUT_FILE, 'w') as f:
      for h, line in pmonitor(popens, timeoutms=500):
        if h:
          f.write('<%s>: %s' % (h.name, line))
        if time() >= endTime:
          leader.send_signal( SIGINT )
          popens[s2].send_signal( SIGINT )
          popens[s3].send_signal( SIGINT )

    f.close()

    print "*** Ending experiment..."

    print "*** Running CLI"
    CLI(net)

    print "*** Stopping network"
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    topology()
