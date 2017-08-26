import matplotlib.pyplot as plot
import numpy as np
import pandas

FILE_NAME = "dump.csv"

data = pandas.read_csv(FILE_NAME)

angles = data["index"]
distances = data["distance"]

xs = {}
ys = {}

# Convert the distance and angle measurements to x-y coords
for i in range(0, len(angles)):
    xs[i] = distances[i] * np.cos(np.deg2rad(angles[i]))
    ys[i] = distances[i] * np.sin(np.deg2rad(angles[i]))

plot.title("LIDAR Surroundings")
plot.grid(True)
plot.gca().set_aspect('equal', adjustable='box')

plot.scatter(xs.values(), ys.values(), c="blue", s=3)
plot.scatter(0, 0, c="red", s=100)

plot.show()