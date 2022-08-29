fb(x) = {
    if x == 1 {
        return 1
    }

    if x == 2 {
        return 1
    }

    return fb(x - 1) + fb(x - 2)

}

print fb(1)
print fb(2)
print fb(3)
print fb(4)
print fb(5)