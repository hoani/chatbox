import RGB1602
import argparse

lcd=RGB1602.RGB1602(16,2)

parser = argparse.ArgumentParser(description='RGB1602 control')
parser.add_argument("--line1")
parser.add_argument("--line2")
parser.add_argument("--rgb")

args = parser.parse_args()

if args.rgb:
    [r,g,b]=args.rgb.split(',', 3)
    r = int(r)
    g = int(g)
    b = int(b)

    lcd.setRGB(r, g, b)

if args.line1:
    lcd.setCursor(0, 0)
    lcd.printout(args.line1)

if args.line2:
    lcd.setCursor(0, 1)
    lcd.printout(args.line2)



