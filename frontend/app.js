document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('file-input');
    const loadingState = document.getElementById('loading-state');
    const resultsSection = document.getElementById('results-section');
    const originalImage = document.getElementById('original-image');
    const cleanedImage = document.getElementById('cleaned-image');
    const downloadBtn = document.getElementById('download-btn');
    const resetBtn = document.getElementById('reset-btn');

    const toolErase = document.getElementById('tool-erase');
    const toolTranslate = document.getElementById('tool-translate');
    
    let currentEndpoint = '/api/clean';

    // Tool Selection Logic
    toolErase.addEventListener('click', (e) => {
        e.preventDefault();
        currentEndpoint = '/api/clean';
        toolErase.classList.add('active');
        toolTranslate.classList.remove('active');
        document.querySelector('header h1').innerText = 'Text Eraser';
        document.querySelector('header p').innerText = 'Pixel-Perfect AI Text Removal';
        document.querySelector('#loading-state p').innerText = 'YOLO ONNX Segmenting Bubbles...';
    });

    toolTranslate.addEventListener('click', (e) => {
        e.preventDefault();
        currentEndpoint = '/api/translate';
        toolTranslate.classList.add('active');
        toolErase.classList.remove('active');
        document.querySelector('header h1').innerText = 'Auto Translator';
        document.querySelector('header p').innerText = 'AI Bangla Translation Rendering Test';
        document.querySelector('#loading-state p').innerText = 'Translating & Rendering Bangla Text...';
    });

    // Drag and Drop Events
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.querySelector('.upload-box').classList.add('dragover');
        }, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.querySelector('.upload-box').classList.remove('dragover');
        }, false);
    });

    dropZone.addEventListener('drop', handleDrop, false);
    fileInput.addEventListener('change', handleFileSelect, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const file = dt.files[0];
        handleFile(file);
    }

    function handleFileSelect(e) {
        const file = e.target.files[0];
        handleFile(file);
    }

    function handleFile(file) {
        if (!file || !file.type.startsWith('image/')) {
            alert('Please select a valid image file (JPG, PNG).');
            return;
        }

        // Display original image
        const reader = new FileReader();
        reader.onload = (e) => {
            originalImage.src = e.target.result;
        };
        reader.readAsDataURL(file);

        // Upload and process
        processImage(file);
    }

    async function processImage(file) {
        loadingState.classList.remove('hidden');
        dropZone.querySelector('.upload-box').style.opacity = '0.5';

        const formData = new FormData();
        formData.append('image', file);

        try {
            const response = await fetch(currentEndpoint, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                const errJson = await response.json().catch(() => ({}));
                throw new Error(errJson.error || 'Failed to process image');
            }

            const blob = await response.blob();
            const objectUrl = URL.createObjectURL(blob);
            
            cleanedImage.src = objectUrl;
            downloadBtn.href = objectUrl;
            downloadBtn.download = `cleaned_${file.name}`;

            dropZone.classList.add('hidden');
            resultsSection.classList.remove('hidden');
        } catch (error) {
            alert(`Error: ${error.message}`);
        } finally {
            loadingState.classList.add('hidden');
            dropZone.querySelector('.upload-box').style.opacity = '1';
        }
    }

    resetBtn.addEventListener('click', () => {
        resultsSection.classList.add('hidden');
        dropZone.classList.remove('hidden');
        fileInput.value = '';
    });
});
