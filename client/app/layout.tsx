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
        <h1 className="text-4xl text-indigo-800 font-bold my-2">
          画像サイズを下げる
        </h1>
        <ul className="flex bg-blue-600 mb-4 pl-2">
          <li className="block px-4 py-2 my-1 hover:bg-gray-100 rounded">
            <Link className="no-underline text-blue-300" href="/">
              Home
            </Link>
          </li>
        </ul>
        <div className="ml-2">{children}</div>
      </body>
    </html>
  );
}
