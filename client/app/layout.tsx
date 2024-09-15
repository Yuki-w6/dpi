import Link from 'next/link';

import './ui/global.css';
import { Noto_Sans_JP } from 'next/font/google';

// Googleフォントを有効化
const fnt = Noto_Sans_JP({ subsets: ['latin'] });

export const metadata = {
  title: 'reduce image size',
  description: '画像サイズを下げるためのアプリ',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body className={fnt.className}>
        <header className="h-16 shadow">
          <nav></nav>
        </header>
        <div className="w-full p-6">
          <div className="tool">
            <h1 className="text-4xl font-semibold text-center">
              画像のサイズを下げる
            </h1>
            <h2 className="text-xl text-center text-gray-800 mt-2">
              アップロードした画像の解像度を下げ、画像のサイズを下げます。
            </h2>
            <div className="ml-2">{children}</div>
          </div>
        </div>
      </body>
    </html>
  );
}
