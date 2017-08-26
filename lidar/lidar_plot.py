import matplotlib.pyplot as plot
import numpy as np
import pandas
import socket

FILE_NAME = "dump.csv"
SIGNAL_THRESHOLD = 50

# Convert the distance and angle measurements to x-y coords
def distances_to_xy(angles, distances, signal):
    xs = {}
    ys = {}
    for i in range(0, len(angles)):
        if signal[i] > SIGNAL_THRESHOLD:
            xs[i] = distances[i] * np.cos(np.deg2rad(angles[i]))
            ys[i] = distances[i] * np.sin(np.deg2rad(angles[i]))
    return xs.values(), ys.values()

def distance_to_xy(angle, distance, signal):
    x = 0
    y = 0
    # if signal > SIGNAL_THRESHOLD:
    x = distance * np.cos(np.deg2rad(angle))
    y = distance * np.sin(np.deg2rad(angle))
    return x, y

def plot_setup():
    plot.title("LIDAR Surroundings")
    plot.grid(True)
    plot.axis((-400, 400, -400, 400))
    plot.gca().set_aspect('equal', adjustable='box')

def main():
    plot_setup()
    plot.ion()
    plot.show()

    s = socket.socket()

    s.connect(("192.168.1.138", 9000))

    xs = {}
    ys = {}

    i = 0
    j = 0

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

                    xs[i], ys[i] = distance_to_xy(index, distance, signal)

                    i = i + 1
                    j = j + 1

                    if i == 71:
                        print(i)
                        plot.scatter(xs.values(), ys.values(), c="blue", s=3)
                        plot.pause(0.01)
                        i = 0

                    if j == 359:
                        plot.clf()
                        plot_setup()
                        j = 0

if __name__ == "__main__":
    main()
