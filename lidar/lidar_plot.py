import sys
import socket
from pyqtgraph.Qt import QtGui, QtCore
from PyQt4.QtCore import QThread
import pyqtgraph as pg
import numpy as np

SIGNAL_THRESHOLD = 50
SERVER_URL = "192.168.1.138"
SERVER_PORT = 9000

p6 = None
curve = None
ptr = 0

xs = []
ys = []

class LidarSocketThread(QThread):
    def __init__(self):
        QThread.__init__(self)

    def __del__(self):
        self.wait()

    def run(self):
        s = socket.socket()
        s.connect((SERVER_URL, SERVER_PORT))

        i = 0

        while True:
            b = s.recv(512).decode("utf-8")
            lines = b.split("\n")

            for l in lines:
                nums = l.split(",")
                if len(nums) == 3:
                    if nums[0].startswith("*") and nums[2].endswith("!"):
                        # valid message
                        index = int(nums[0][1:])
                        distance = int(nums[1])
                        signal = int(nums[2][:-1])

                        if signal > SIGNAL_THRESHOLD:
                            x, y = distance_to_xy(index, distance, signal)

                            xs.append(x)
                            ys.append(y)

                            if i < 359:
                                i = i + 1
                            else:
                                xs.pop(0)
                                ys.pop(0)

def distance_to_xy(angle, distance, signal):
    x = 0
    y = 0
    # if signal > SIGNAL_THRESHOLD:
    x = distance * np.cos(np.deg2rad(angle))
    y = distance * np.sin(np.deg2rad(angle))
    return x, y

def update():
    global curve, ptr, p6, xs, ys
    curve.setData(xs, ys)
    if ptr == 0:
        p6.enableAutoRange('xy', False)  ## stop auto-scaling after the first data set is plotted
    ptr += 1

def main():
    app = QtGui.QApplication([])
    win = pg.GraphicsWindow(title="LIDAR Real-time Display")
    win.resize(1280,720)
    win.setWindowTitle('LIDAR Display')

    pg.setConfigOptions(antialias=True)

    global p6, curve, ptr
    p6 = win.addPlot(title="LIDAR")
    curve = p6.plot(pen=None, symbolSize=10)
    ptr = 0

    timer = QtCore.QTimer()
    timer.timeout.connect(update)
    timer.start(50)

    # Start background thread
    t = LidarSocketThread()
    t.start()

    QtGui.QApplication.instance().exec_()

if __name__ == "__main__":
    import sys
    if (sys.flags.interactive != 1) or not hasattr(QtCore, 'PYQT_VERSION'):
        main()