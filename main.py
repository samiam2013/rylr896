#!/usr/bin/python3
import time 
import argparse
from luma.core.interface.serial import i2c
from luma.core.render import canvas
from luma.oled.device import sh1106

# I2C interface and address
serial = i2c(port=1, address=0x3c)
# OLED Device
device = sh1106(serial, rotate=2)

parser = argparse.ArgumentParser()
parser.add_argument("-r","--rssi", help="Recieved Signal Strength Indicator RSSI (digit)", type=int)
parser.add_argument("-p","--persist", help="Seconds to hold the displayed digit", type=int)
parser.add_argument("-l","--last", help="Time since last recieve", type=str)
args = parser.parse_args()
# STARTUP SCREEN
with canvas(device) as draw:
    draw.rectangle(device.bounding_box, outline="white", fill="black")
    draw.text((30, 15), "LAST: "+args.last, fill="white")
    draw.text((35, 30), "RSSI: "+str(args.rssi), fill="white")
    
time.sleep(args.persist)