# multipart-go-demo
Multipart File Upload Server

### running it

```
docker run -p 8199:8199 -it skhatri/multipart-file-upload
```

### testing it

```
curl -X POST -H "Content-Type: multipart/form-data; boundary=123-UPLOAD-SEPARATOR" \
-d '--123-UPLOAD-SEPARATOR
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

test
--123-UPLOAD-SEPARATOR--
' "http://localhost:8199"
```
