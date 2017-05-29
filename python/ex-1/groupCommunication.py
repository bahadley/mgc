#!/usr/bin/python

from time import sleep

from mininet.net import Mininet
from mininet.node import Controller, OVSKernelAP
from mininet.link import TCLink
from mininet.cli import CLI
from mininet.log import setLogLevel

def topology():
    "Create a network."
    net = Mininet(controller=Controller, link=TCLink, accessPoint=OVSKernelAP)

    print "*** Creating nodes"
    sta1 = net.addStation('sta1', mac='00:00:00:00:00:02', ip='10.0.0.2/8')
    sta2 = net.addStation('sta2', mac='00:00:00:00:00:03', ip='10.0.0.3/8')
    sta3 = net.addStation('sta3', mac='00:00:00:00:00:04', ip='10.0.0.4/8')
    ap1 = net.addAccessPoint('ap1', ssid='ssid-ap1', mode='g', channel='11', position='115,62,0')
    ap2 = net.addAccessPoint('ap2', ssid='ssid-ap2', mode='g', channel='1', position='57,142,0')
    c1 = net.addController('c1', controller=Controller)

    print "*** Configuring wifi nodes"
    net.configureWifiNodes()

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

    print "*** Waiting until stations are in range of ap2..."
    sleep(55)

    print "*** Starting test..."
    s1 = net.get('sta1')
    s2 = net.get('sta2')
    outfile = '/tmp/%s.out' % s2.name
    s2.cmd('./rcv.py 2>&1 > %s &' % outfile)
    sleep(1)
    s1.cmd('./send.py')
    sleep(1)
    s1.cmd('kill %send.py')
    s2.cmd('kill %rcv.py')
    print "*** Ending test..."

    print "*** Running CLI"
    CLI(net)

    print "*** Stopping network"
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    topology()
