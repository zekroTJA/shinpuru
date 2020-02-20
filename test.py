import random
import math


'Schiebespiel'

field = [

    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 0],

]


def Anzeige():

    global field

    print('\n')
    print('-+---+---+---+')
    i = ("")
    for line in field:
        line_str = "%s | %s | %s | %s |" % (i, line[0], line[1], line[2])
        print(line_str)
        print('-+---+---+---+')


def number_input():
    print("Enter number to change: ")
    global x
    x = input()


def getcoords():
    a = field

    if a[0][0] == x:
        return 0, 0

    if a[1][0] == x:
        return 1, 0

    if a[2][0] == x:
        return 2, 0

    if a[0][1] == x:
        return 0, 1

    if a[0][2] == x:
        return 0, 2

    if a[1][1] == x:
        return 1, 1

    if a[1][2] == x:
        return 1, 2

    if a[2][1] == x:
        return 2, 1

    if a[2][2] == x:
        return 2, 2

    return 0, 0


def nextToZero(x1, y1, x2, y2):
    f = math.sqrt(math.pow(float(x1-x2), 2.0) + math.Pow(float(y1-y2), 2.0))
    return f == 1


def move():
    ix, iy = getcoords(field, x)

    jx, jy = getcoords(field, 0)

    field[jx][jy] = field[ix][iy]

    field[ix][iy] = 0


def checkwon():
    a = field
    if (a[0][0] == 1) and (a[1][0] == 4) and (a[2][0] == 7) and (a[0][1] == 2) and (a[0][2] == 3) and (a[1][1] == 5) and (a[1][2] == 6) and (a[2][1] == 8) and (a[2][2] == 0):
        print("")
        print("! You finished it !")
        print("")
        return True
    else:
        return False


def fillfield():

    arr = [0, 1, 2, 3, 4, 5, 6, 7, 8]

    i = 0

    while i != 9:

        r = random.randint(0, 8)

        zr = arr[r]

        if zr == -1:
            continue

        arr[r] = -1

        x, y = i/3, i % 3

        i + 1

        # intx = int(x)

        # inty = int(y)

        field[x][y] = zr

        if i == 9:
            break
    Anzeige()


def game():

    fillfield()

    number_input()

    getcoords()

    nextToZero()

    if nextToZero() is True:
        move()

    Anzeige()

    checkwon()

    if checkwon() is False:
        return


game()
