import itertools

vals = [2, 9, 5, 7, 3]

def f(a):
    return a[0] + a[1] * (a[2] ** 2) + a[3] ** 3 - a[4]

for order in itertools.permutations(vals):
    if f(order) == 399:
        print(order)
