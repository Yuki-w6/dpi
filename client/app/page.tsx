'use client';

import { useState } from 'react';
import './ui/global.css';
import { Noto_Sans_JP } from 'next/font/google';

// Googleフォントを有効化
const fnt = Noto_Sans_JP({ subsets: ['latin'] });

export default function Page() {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState('');
  const [width, setWidth] = useState(null);
  const [height, setHeight] = useState(null);
  const [x_dpi, setXDPI] = useState(null);
  const [y_dpi, setYDPI] = useState(null);
  const [filename, setFilename] = useState('');

  const handleFileChange = (e: any) => {
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
    setXDPI(data.xDPI);
    setYDPI(data.yDPI);
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
          <p>解像度（水平方向）: {x_dpi}</p>
          <p>解像度（垂直方向）: {y_dpi}</p>
        </div>
      )}
    </div>
  );
}
