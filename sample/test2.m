fb(x) = {
    if x == 1 {
        return 1
    }
    if x == 2 {
        return 1
    }
    return fb(x - 1) + fb(x - 2)
}

i = 1
for i != 10 {
    print fb(i)
    i = i + 1
}