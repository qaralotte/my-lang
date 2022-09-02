isPrime(n) = {
    i = 2
    for i < n {
        if n % i == 0 {
            return false
        }
        i = i + 1
    }
    return true
}

i = 3
for i < 10 {
    print isPrime(i)
    print i
    i = i + 1
}

