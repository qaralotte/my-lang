fb(x) = {
    if x == 0 {
        return 0
    }
    if x == 1 {
        return 1
    }
    return fb(x - 1) + fb(x - 2)
}

i = 0
for i != 10 {
    print fb(i)
    i = i + 1
}