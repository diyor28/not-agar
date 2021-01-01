import math

def intersection_area(R, r, dist):
    dist1 = (R**2 - r**2 + dist**2)/(2*dist)
    dist2 = dist - dist1
    return R**2 * math.acos(dist1/R) - dist1 * math.sqrt(R**2 - dist1**2) + r**2 * math.acos(dist2/r) - dist2*math.sqrt(r**2 - dist2**2)

R = 3
r = 3
dist = 0.2
print(math.pi * R * R)
print(math.pi * r * r)
print(intersection_area(3, 3, 0.2))
