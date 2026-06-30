# Manhwa Tools API Documentation

## Base URL
All API requests should be made to `http://localhost:8080`

---

## 1. Text Eraser Tool
Removes text from speech bubbles in Manhwa pages using AI instance segmentation.

### Endpoint
`POST /api/clean`

### Request Format
- **Content-Type**: `multipart/form-data`
- **Body**:
  - `image`: The image file to be processed (supported formats: `.jpg`, `.png`, `.jpeg`).

### Response
- **Status 200 OK**: Returns the cleaned raw image binary (`image/jpeg`).
- **Status 400 Bad Request**: Missing or invalid image file.
- **Status 500 Internal Server Error**: ONNX model or processing failure.

### Example (cURL)
```bash
# Upload a page to the eraser and save the output
curl -X POST -F "image=@/path/to/page.jpg" http://localhost:8080/api/clean -o cleaned_page.jpg
```

---

## Modularity and Extending the API
This project is built using a modular Gin router. To add new tools (e.g., `POST /api/translate`), create a new handler in `backend/internal/api/tools/` and register it in `backend/internal/api/router.go`.
