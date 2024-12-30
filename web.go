package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	OriginalText    string
	CompressedText  string
	CompressedLen   int
	OriginalLen     int
	OriginalBitLen  int
	CompressionRate float64
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Dinamik Huffman Sıkıştırma</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f5f5f5;
            padding: 20px;
            border-radius: 5px;
            margin-top: 20px;
        }
        .result {
            background-color: #fff;
            padding: 15px;
            border-radius: 5px;
            margin-top: 10px;
        }
        textarea {
            width: 100%;
            padding: 10px;
            margin: 10px 0;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #45a049;
        }
    </style>
</head>
<body>
    <h1>Dinamik Huffman Sıkıştırma</h1>
    <div class="container">
        <form method="POST">
            <h3>Sıkıştırılacak Metni Girin:</h3>
            <textarea name="text" rows="4" required>{{.OriginalText}}</textarea>
            <br>
            <button type="submit">Sıkıştır</button>
        </form>
    </div>

    {{if .CompressedText}}
    <div class="container">
        <h3>Sonuçlar:</h3>
        <div class="result">
            <p><strong>Orijinal Metin:</strong> {{.OriginalText}}</p>
            <p><strong>Sıkıştırılmış Metin (Binary):</strong> {{.CompressedText}}</p>
            <p><strong>Orijinal Boyut:</strong> {{.OriginalLen}} byte ({{.OriginalBitLen}} bit)</p>
            <p><strong>Sıkıştırılmış Boyut:</strong> {{.CompressedLen}} bit</p>
            <p><strong>Sıkıştırma Oranı:</strong> %{{printf "%.2f" .CompressionRate}}</p>
        </div>
    </div>
    {{end}}
</body>
</html>
`

func handleCompression(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("page").Parse(htmlTemplate))
	data := &PageData{}

	if r.Method == "POST" {
		text := r.FormValue("text")
		if text != "" {
			// Dinamik Huffman ağacını oluştur
			dh := NewDynamicHuffman()

			// Metni sıkıştır
			compressed := ""
			for i := 0; i < len(text); i++ {
				char := text[i]
				code := dh.GetCode(char)
				compressed += code
				dh.UpdateTree(char)
			}

			data.OriginalText = text
			data.CompressedText = compressed
			data.CompressedLen = len(compressed)
			data.OriginalLen = len(text)
			data.OriginalBitLen = len(text) * 8
			data.CompressionRate = 100 - (float64(len(compressed)) / float64(len(text)*8) * 100)
		}
	}

	tmpl.Execute(w, data)
}

func startServer() {
	http.HandleFunc("/", handleCompression)
	println("Web sunucusu başlatıldı. http://localhost:8080 adresini ziyaret edin")
	http.ListenAndServe(":8080", nil)
}
