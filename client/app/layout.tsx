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
    <html lang="ja" className="h-full">
      <body className={`h-full overflow-hidden ${fnt.className}`}>
        <header className="block fixed top-0 left-0 right-0 h-16 z-10">
          <nav></nav>
        </header>
        <div className="h-full mt-16">
          <div className="h-full w-full p-6 bg-gray-100">{children}</div>
        </div>
        <div className="block fixed left-0 right-0 bottom-0 h-12 z-10 bg-white">
          <div className="text-center p-2.5">test</div>
        </div>
      </body>
    </html>
  );
}
