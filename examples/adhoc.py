#!/usr/bin/python

from mininet.net import Mininet
from mininet.node import Controller, OVSKernelAP
from mininet.cli import CLI
from mininet.link import TCLink
from mininet.log import setLogLevel

def topology():
    "Create a network."
    net = Mininet(controller=Controller, link=TCLink, accessPoint=OVSKernelAP)

    print "*** Creating nodes"
    m1 = net.addStation('m1', wlans=2)
    m2 = net.addStation('m2', wlans=2)
    b1 = net.addStation('b1', position='35,142,0')
    ap1 = net.addAccessPoint('ap1', ssid='ssid-ap1', mode='g', channel='11', position='115,62,0')
    ap2 = net.addAccessPoint('ap2', ssid='ssid-ap2', mode='g', channel='1', position='57,142,0')
    c1 = net.addController('c1', controller=Controller)

    print "*** Configuring wifi nodes"
    net.configureWifiNodes()

    m1.setIP('10.0.0.1/8', intf="m1-wlan0")
    m2.setIP('10.0.0.2/8', intf="m2-wlan0")
    m1.setIP('192.168.10.1/24', intf="m1-wlan1")
    m2.setIP('192.168.10.2/24', intf="m2-wlan1")
    b1.setIP('192.168.10.3/24', intf="b1-wlan0")

    print "*** Creating links"
    net.addHoc(m1, ssid='adhocNet', mode='g')
    net.addHoc(m2, ssid='adhocNet', mode='g')
    net.addLink(ap2, b1)

    print "*** Starting network"
    net.build()
    c1.start()
    ap1.start([c1])
    ap2.start([c1])

    net.plotGraph(max_x=200, max_y=200)

    net.associationControl('ssf')

    net.startMobility(time=0)
    net.mobility(m1, 'start', time=1, position='86,188,0')
    net.mobility(m2, 'start', time=1, position='78,195,0')
    net.mobility(m1, 'stop', time=250, position='86,0,0')
    net.mobility(m2, 'stop', time=250, position='78,7,0')
    net.stopMobility(time=250)

    print "*** Running CLI"
    CLI(net)

    print "*** Stopping network"
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    topology()
