package main

import (
	"encoding/json"
	"fmt"
	"github.com/for6to9si/vpnparser/pkgs/outbound"
	"net/url"
	"os"
	"strings"
	"unicode"
)

// replaceInvalidChars заменяет недопустимые символы в имени файла на подчёркивания
func replaceInvalidChars(name string) (string, error) {
	// Рекурсивно декодируем URL-кодированные символы
	decodedName := name
	for {
		newDecoded, err := url.QueryUnescape(decodedName)
		if err != nil || newDecoded == decodedName {
			break
		}
		decodedName = newDecoded
	}

	// Недопустимые символы для имён файлов в большинстве ОС
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}

	// Заменяем недопустимые символы
	for _, char := range invalidChars {
		decodedName = strings.ReplaceAll(decodedName, char, "_")
	}

	// Удаляем или заменяем непечатаемые символы и эмодзи
	var result strings.Builder
	for _, r := range decodedName {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}

	// Удаляем пробелы в начале и конце
	cleaned := strings.TrimSpace(result.String())

	// Удаляем повторяющиеся подчёркивания
	cleaned = strings.ReplaceAll(cleaned, "__", "_")
	cleaned = strings.ReplaceAll(cleaned, "__", "_") // Дважды на случай тройных

	// Ограничиваем длину имени файла
	if len(cleaned) > 100 {
		cleaned = cleaned[:100]
	}

	// Если после всех преобразований имя пустое, создаем дефолтное
	if cleaned == "" {
		cleaned = "unnamed_config"
	}

	return cleaned, nil
}

// isEmoji проверяет, является ли символ эмодзи
func isEmoji(r rune) bool {
	// Более точная проверка эмодзи (диапазоны Unicode для эмодзи)
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1F1E6 && r <= 0x1F1FF) // Flags
}

func extractComment(input string) string {
	// Разделяем строку по символу '#'
	parts := strings.SplitN(input, "#", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func createFile(filename string) error {
	// Очищаем имя файла от недопустимых символов
	cleanedFilename := filename
	if cleanedFilename == "" {
		return fmt.Errorf("имя файла пустое после очистки")
	}

	// Проверяем длину имени файла
	if len(cleanedFilename) > 255 {
		cleanedFilename = cleanedFilename[:255]
	}

	// Создаём файл
	file, err := os.Create(cleanedFilename)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла '%s': %v", cleanedFilename, err)
	}
	defer file.Close()
	return nil
}

// decodeURLComment декодирует URL-кодированную строку рекурсивно
func decodeURLComment(comment string) (string, error) {
	decoded := comment
	for {
		newDecoded, err := url.QueryUnescape(decoded)
		if err != nil {
			return "", fmt.Errorf("ошибка декодирования URL: %v", err)
		}
		if newDecoded == decoded {
			break
		}
		decoded = newDecoded
	}
	return decoded, nil
}

// main function parses a VLESS URI and outputs formatted JSON
func main() {
	// Define the VLESS URI
	vlessURI := []string{"vless://e7cc66ec-4b1dafd8-20283e67d656@power.dlmddr:2020?security=none&encryption=none&host=Soddr&headerType=http&type=tcp#qnsh"}

	for i, input := range vlessURI {

		// Обработка каждой строки

		comment := extractComment(input)
		if comment == "" {
			fmt.Printf("Строка %d: Пропущена (нет комментария после #)\n", i+1)
			continue
		}

		// Декодируем URL-кодированную строку
		decodedComment, err := replaceInvalidChars(comment)
		if err != nil {
			fmt.Printf("Строка %d: Ошибка декодирования: %v\n", i+1, err)
			continue
		}

		// Инициализируем и парсим конфигурацию (с обработкой ошибок)
		ob := outbound.GetOutbound(outbound.XrayCore, input)
		if ob == nil {
			fmt.Printf("Строка %d: Неподдерживаемый протокол: %s\n", i+1, input)
			continue
		}
		// Парсим с обработкой паники
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Строка %d: Ошибка парсинга (panic recovered): %v\n", i+1, r)
				}
			}()
			ob.Parse(input)
		}()

		// Get the outbound configuration
		config := ob.GetOutboundStr()

		if config == "" {
			fmt.Printf("Строка %d: Не удалось распарсить конфигурацию\n", i+1)
			continue
		}

		// Check if config is already a JSON string
		var jsonData []byte
		//var err error

		// Try to treat config as a JSON string first
		var temp map[string]interface{}
		if err := json.Unmarshal([]byte(config), &temp); err == nil {
			// If config is a valid JSON string, re-serialize it with proper formattin
			temp["tag"] = decodedComment
			jsonData, err = json.MarshalIndent(temp, "", "  ")
		} else {
			// Если config не JSON, создаем новый объект
			temp = map[string]interface{}{
				"config": config,
				"tag":    decodedComment,
			}
			// If config is not a JSON string, assume it's a struct and serialize it
			jsonData, err = json.MarshalIndent(temp, "", "  ")
		}

		if err != nil {
			fmt.Printf("Строка %d: Ошибка сериализации JSON: %v\n", i+1, err)
			continue
		}

		// Print the formatted JSON to console
		fmt.Println(string(jsonData))

		fmt.Printf("Строка %d: %s\n", i+1, decodedComment)
		//err = createFile(decodedComment)
		//if err != nil {
		//	fmt.Printf("Ошибка при создании файла для строки %d: %v\n", i+1, err)
		//} else {
		//	fmt.Printf("Файл '%s' успешно создан\n", decodedComment)
		//}

		// Создаем файл с очищенным именем
		cleanedName, err := replaceInvalidChars(decodedComment)
		if cleanedName == "" {
			cleanedName = fmt.Sprintf("config_%d", i+1)
		}

		// Save the formatted JSON to a file for verification
		err = os.WriteFile(cleanedName+".json", jsonData, 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		fmt.Println("JSON configuration saved to %v\n", decodedComment)
	}
}
