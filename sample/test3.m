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

j = 2
for j < 10 {
    print isPrime(j)
    j = j + 1
}

