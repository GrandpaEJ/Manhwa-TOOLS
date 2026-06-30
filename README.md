# Manhwa Tools 🚀

A highly scalable, Go-based web application and API for advanced AI processing of Manhwa/Manga images. Currently features a pixel-perfect **Text Eraser** utilizing YOLOv8 and ONNX runtime.

## Features ✨
- **Modular Architecture**: Built with a custom scalable Gin router to effortlessly plug in new tools (Translators, Upscalers, etc.).
- **Text Eraser**: Uses a custom YOLO ONNX segmentation model (`segmentor_best.onnx`) and Moore Neighborhood contour tracing to perfectly mask and white-out speech bubbles.
- **Glassmorphism UI**: A beautiful drag-and-drop frontend interface for quick manual processing.
- **High Performance**: Native Go ONNX inference completely eliminating heavy Python dependencies.

## Installation 🛠️

1. Clone the repository:
```bash
git clone https://github.com/GrandpaEJ/Manhwa-TOOLS.git
cd Manhwa-TOOLS
```

2. Start the backend server:
```bash
cd backend
go run main.go
```

3. Access the web app at [http://localhost:8080](http://localhost:8080).

## Documentation 📚
- [API Documentation](docs/api.md) - Details on how to interface with the REST API endpoints.

## Future Roadmap 🗺️
- Auto Translator Integration (LLM-based)
- Typography Rendering
- Image Upscaling
