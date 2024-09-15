import Link from 'next/link';

import './ui/global.css';
import { Inconsolata } from 'next/font/google';

// Googleフォントを有効化
const fnt = Inconsolata({ subsets: ['latin'] });

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
        <header className="header">
          <nav></nav>
        </header>
        <div className="main">
          <h1 className="text-4xl text-indigo-800 font-bold my-2 center">
            画像サイズを下げる
          </h1>
          <div className="ml-2">{children}</div>
        </div>
      </body>
    </html>
  );
}
