'use client';

import { useState } from 'react';

export default function Page() {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState('');
  const [width, setWidth] = useState(null);
  const [height, setHeight] = useState(null);
  const [filename, setFilename] = useState('');

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage('ファイルが選択されていません。');
      return;
    }
    const formData = new FormData();
    formData.append('image', file);

    const response = await fetch('http://localhost:1000/api/v1/upload', {
      method: 'POST',
      body: formData,
    });

    const data = await response.json();
    setMessage(data.message);
    setWidth(data.width);
    setHeight(data.height);
    setFilename(data.filename);
  };

  return (
    <div>
      <h1>画像アップロード</h1>
      <input type="file" onChange={handleFileChange} />
      <button onClick={handleUpload}>アップロード</button>
      {message && (
        <div>
          <p>{message}</p>
          <p>ファイル名: {filename}</p>
          <p>幅: {width}</p>
          <p>高さ: {height}</p>
        </div>
      )}
    </div>
  );
}
