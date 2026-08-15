package helper

import (
	"regexp"
	"strconv"
	"strings"
)

// IsValidCNPJ valida um número de CNPJ (com ou sem máscara)
func IsValidCNPJ(cnpj string) bool {
	// Remove caracteres não numéricos
	re := regexp.MustCompile(`\D`)
	cleanCNPJ := re.ReplaceAllString(cnpj, "")

	// CNPJ deve ter exatamente 14 dígitos
	if len(cleanCNPJ) != 14 {
		return false
	}

	// Rejeita sequências iguais (ex: 00000000000000)
	if strings.Repeat(string(cleanCNPJ[0]), 14) == cleanCNPJ {
		return false
	}

	// Cálculo do 1º dígito verificador
	sizes1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum1 := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(cleanCNPJ[i]))
		sum1 += digit * sizes1[i]
	}
	result1 := sum1 % 11
	digit1 := 0
	if result1 >= 2 {
		digit1 = 11 - result1
	}

	// Cálculo do 2º dígito verificador
	sizes2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum2 := 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(cleanCNPJ[i]))
		sum2 += digit * sizes2[i]
	}
	result2 := sum2 % 11
	digit2 := 0
	if result2 >= 2 {
		digit2 = 11 - result2
	}

	// Verifica se os dígitos calculados batem com os informados
	calcDigit1, _ := strconv.Atoi(string(cleanCNPJ[12]))
	calcDigit2, _ := strconv.Atoi(string(cleanCNPJ[13]))

	return calcDigit1 == digit1 && calcDigit2 == digit2
}
