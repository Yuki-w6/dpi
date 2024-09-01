import React, { useState } from "react";
import axios from "axios";

const TestImageUpload: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);

  const onFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      const selectedFile = files[0];
      setFile(selectedFile);
    } else {
      setFile(null);
    }
  };

  const uploadImage = async () => {
    if (!file) {
      console.error("No file selected");
      return;
    }
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await axios.post(
        "{GO_API_ENDPOINT}",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );
      console.log(response.data);
    } catch (error) {
      console.error("Upload failed:", error);
    }
  };

  return (
    <form
      onSubmit={(e) => {
        uploadImage();
        e.preventDefault();
      }}
    >
      <input type="file" accept="image/*" onChange={(e) => onFileChange(e)} />
      <button type="submit">upload image</button>
    </form>
  );
};

export default TestImageUpload;