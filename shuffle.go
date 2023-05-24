package main

import "strings"

func decryptInner(data string, key uint32) string {
	decrypted := strings.Split(data, ",")
	arr1 := make([]int, len(decrypted))
	arr2 := make([]int, len(decrypted))

	for jv := 0; jv < len(decrypted); jv++ {
		dV := ((key >> 8) & 65535) % uint32(len(decrypted))
		key = ((key * 65793) + 4282663) & 8388607
		rV := ((key >> 8) & 65535) % uint32(len(decrypted))
		key = ((key * 65793) + 4282663) & 8388607

		arr1[len(decrypted)-1-jv] = int(dV)
		arr2[len(decrypted)-1-jv] = int(rV)
	}

	for i := 0; i < len(decrypted); i++ {
		decrypted[arr1[i]], decrypted[arr2[i]] = decrypted[arr2[i]], decrypted[arr1[i]]
	}

	return strings.Join(decrypted, ",")
}
