# Simple QR-generator on Go

It can generate black&wite and colored codes from URL and provides two ways to work with it:
1) CLI
2) HTTP

## Download

```bash
go install github.com/Godreck/qr-gen@latest

# or local:
go build -o qr .
```

## Usage

**CLI**
```
./qr -url https://example.com -out qr.png -fg 3498db -bg ffffff
```

**HTTP**
```
./qr -http
# -> http://localhost:8080/qr?url=https://example.com&fg=1e88e5&bg=ffffff
```

**Docker**
```
docker build -t qr-gen .
docker run -p 8080:8080 qr-gen

---

