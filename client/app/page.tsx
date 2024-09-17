'use client';

import { useState } from 'react';
import './ui/global.css';

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
    <div className="text-center">
      <h1 className="text-4xl font-semibold text-center">
        画像のサイズを下げる
      </h1>
      <h2 className="text-xl text-center text-gray-800 mt-2 pb-8">
        アップロードした画像の解像度を下げ、画像のサイズを下げます。
      </h2>

      <button className="bg-blue-400 hover:bg-blue-600 text-white text-xl font-bold py-6 px-24 rounded-xl">
        画像を選択
      </button>
      {/* ドロップメッセージ */}
      <p className="mt-4 text-gray-500">
        または、ここに画像をドロップしてください
      </p>
    </div>
  );
}
